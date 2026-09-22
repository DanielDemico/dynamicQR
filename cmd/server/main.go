package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"dynamicqr/internal/config"
	"dynamicqr/internal/database"
	"dynamicqr/internal/handler"
	"dynamicqr/internal/repository"
	"dynamicqr/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	cfg := config.Load()
	log.Printf("[DynamicQR] Iniciando servidor na porta %s...", cfg.Port)
	log.Printf("[DynamicQR] Base URL configurada: %s", cfg.BaseURL)
	log.Printf("[DynamicQR] Conectando ao MongoDB em: %s (DB: %s)", cfg.MongoURI, cfg.MongoDB)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	mongoDB, err := database.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("[DynamicQR] Falha crítica ao conectar no MongoDB: %v", err)
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer closeCancel()
		_ = mongoDB.Close(closeCtx)
	}()
	log.Println("[DynamicQR] Conexão com MongoDB e índices configurados com sucesso.")

	// Inicialização de camadas
	qrRepo := repository.NewQRCodeRepository(mongoDB.Database)
	qrService := service.NewQRService()
	qrHandler := handler.NewQRHandler(cfg, qrRepo, qrService)
	redirectHandler := handler.NewRedirectHandler(qrRepo)

	// Configuração do Router
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Configuração de CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link", "Content-Disposition"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Healthcheck
	r.Get("/healthz", qrHandler.Healthz)

	// Rotas da API REST
	r.Route("/api/v1/qr-codes", func(api chi.Router) {
		api.Post("/", qrHandler.Create)
		api.Get("/", qrHandler.List)
		api.Get("/{id}", qrHandler.GetByID)
		api.Put("/{id}", qrHandler.Update)
		api.Delete("/{id}", qrHandler.Delete)
		api.Get("/{id}/image", qrHandler.GetImage)
	})

	// Servir arquivos estáticos do frontend
	staticDir := http.Dir("./web/static")
	r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(staticDir)))

	// Rota raiz para o painel de gestão
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/static/index.html")
	})

	// Redirecionamento público por slug
	r.Get("/{slug}", redirectHandler.Redirect)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}

	// Inicialização assíncrona do servidor HTTP
	serverErrors := make(chan error, 1)
	go func() {
		log.Printf("[DynamicQR] Servidor pronto para atender requisições em http://0.0.0.0:%s", cfg.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Graceful Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("[DynamicQR] Erro fatal no servidor HTTP: %v", err)
		}
	case sig := <-shutdown:
		log.Printf("[DynamicQR] Sinal recebido: %v. Encerrando graciosamente...", sig)
		ctxShutdown, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancelShutdown()

		if err := server.Shutdown(ctxShutdown); err != nil {
			log.Printf("[DynamicQR] Erro ao encerrar servidor HTTP: %v", err)
			_ = server.Close()
		}
	}
	log.Println("[DynamicQR] Servidor finalizado com sucesso.")
}

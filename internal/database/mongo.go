package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect conecta ao MongoDB, valida com ping e cria os índices obrigatórios
func Connect(ctx context.Context, uri, dbName string) (*MongoDB, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOpts)
	if err != nil {
		return nil, fmt.Errorf("erro ao inicializar cliente MongoDB: %w", err)
	}

	// Ping com timeout
	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, nil); err != nil {
		return nil, fmt.Errorf("erro no ping do MongoDB: %w", err)
	}

	db := client.Database(dbName)
	m := &MongoDB{
		Client:   client,
		Database: db,
	}

	if err := m.ensureIndexes(ctx); err != nil {
		return nil, fmt.Errorf("erro ao configurar índices do MongoDB: %w", err)
	}

	return m, nil
}

func (m *MongoDB) ensureIndexes(ctx context.Context) error {
	col := m.Database.Collection("qr_codes")

	// 1. Índice único em "slug"
	slugIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "slug", Value: 1}},
		Options: options.Index().SetUnique(true).SetName("idx_slug_unique"),
	}

	// 2. Índice em "created_at" decrescente
	createdAtIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "created_at", Value: -1}},
		Options: options.Index().SetName("idx_created_at_desc"),
	}

	// 3. Índice textual em "title"
	titleIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "title", Value: "text"}},
		Options: options.Index().SetName("idx_title_text"),
	}

	_, err := col.Indexes().CreateMany(ctx, []mongo.IndexModel{slugIndex, createdAtIndex, titleIndex})
	return err
}

func (m *MongoDB) Close(ctx context.Context) error {
	return m.Client.Disconnect(ctx)
}

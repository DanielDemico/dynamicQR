package config

import (
	"bufio"
	"os"
	"strings"
	"time"
)

type Config struct {
	Port         string
	BaseURL      string
	MongoURI     string
	MongoDB      string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// Load reads configuration from environment variables and an optional .env file
func Load() *Config {
	loadDotEnv(".env")

	port := getEnv("PORT", "8080")
	baseURL := getEnv("BASE_URL", "http://localhost:8080")
	mongoURI := getEnv("MONGO_URI", "mongodb://localhost:27017")
	mongoDB := getEnv("MONGO_DB", "dynamicqr_db")
	readTimeoutStr := getEnv("READ_TIMEOUT", "5s")
	writeTimeoutStr := getEnv("WRITE_TIMEOUT", "10s")

	readTimeout, err := time.ParseDuration(readTimeoutStr)
	if err != nil {
		readTimeout = 5 * time.Second
	}

	writeTimeout, err := time.ParseDuration(writeTimeoutStr)
	if err != nil {
		writeTimeout = 10 * time.Second
	}

	return &Config{
		Port:         port,
		BaseURL:      baseURL,
		MongoURI:     mongoURI,
		MongoDB:      mongoDB,
		ReadTimeout:  readTimeout,
		WriteTimeout: writeTimeout,
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && strings.TrimSpace(val) != "" {
		return strings.TrimSpace(val)
	}
	return defaultVal
}

func loadDotEnv(filepath string) {
	file, err := os.Open(filepath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			k := strings.TrimSpace(parts[0])
			v := strings.TrimSpace(parts[1])
			// remove quotes if present
			v = strings.Trim(v, `"'`)
			if _, exists := os.LookupEnv(k); !exists {
				_ = os.Setenv(k, v)
			}
		}
	}
}

package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoConfig agrupa as configurações de ligação ao MongoDB.
type MongoConfig struct {
	URI      string
	Database string
}

// MongoConfigFromEnv lê a configuração de variáveis de ambiente.
func MongoConfigFromEnv() MongoConfig {
	return MongoConfig{
		URI:      getEnv("MONGODB_URL", "mongodb://localhost:27017"),
		Database: getEnv("MONGODB_DB", "gorm_crm_logs"),
	}
}

// NewMongo liga ao MongoDB e verifica a ligação com Ping.
// MongoDB usa uma ligação lazy — só liga realmente no primeiro acesso.
// Ping força a verificação imediata, útil no startup.
func NewMongo(cfg MongoConfig) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("connect mongodb: %w", err)
	}

	// Verifica se o servidor está acessível
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("ping mongodb: %w", err)
	}

	return client.Database(cfg.Database), nil
}


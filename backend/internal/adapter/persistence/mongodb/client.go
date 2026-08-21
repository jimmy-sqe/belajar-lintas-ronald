// Package mongodb is the MongoDB persistence adapter. It implements the
// item.Repository port using the official Mongo Go driver.
package mongodb

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds MongoDB connection settings.
type Config struct {
	URI string `mapstructure:"MONGODB_URI"`
	DB  string `mapstructure:"MONGODB_DB"`
}

// Connect returns a connected mongo.Client.
func Connect(ctx context.Context, cfg Config) (*mongo.Client, error) {
	if cfg.URI == "" {
		return nil, fmt.Errorf("mongodb: URI is required")
	}
	connCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connCtx, options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("mongodb: connect: %w", err)
	}
	if err := client.Ping(connCtx, nil); err != nil {
		_ = client.Disconnect(context.Background())
		return nil, fmt.Errorf("mongodb: ping: %w", err)
	}
	return client, nil
}

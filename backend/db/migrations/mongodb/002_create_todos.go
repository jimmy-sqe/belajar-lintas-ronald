package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func Up_002_create_todos(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("todos").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{Keys: bson.D{{Key: "created_by", Value: 1}}},
		{Keys: bson.D{{Key: "created_by", Value: 1}, {Key: "created_at", Value: -1}}},
	})
	return err
}

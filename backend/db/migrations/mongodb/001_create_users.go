package migrations

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Up_001_create_users creates indexes on the users collection.
// MongoDB doesn't require CREATE TABLE; this registers the unique
// email index needed for FindByEmail lookups.
func Up_001_create_users(ctx context.Context, db *mongo.Database) error {
	_, err := db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})
	return err
}

package mongodb

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

type UserRepository struct {
	col *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{col: db.Collection("users")}
}

type userDoc struct {
	ID           string `bson:"_id"`
	Email        string `bson:"email"`
	PasswordHash string `bson:"password_hash"`
	Name         string `bson:"name,omitempty"`
	CreatedBy    string `bson:"created_by,omitempty"`
	ModifiedBy   string `bson:"modified_by,omitempty"`
}

func (d userDoc) toDomain() (*authdomain.User, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, err
	}
	u := &authdomain.User{
		ID: id, Email: d.Email, PasswordHash: d.PasswordHash, Name: d.Name,
	}
	if d.CreatedBy != "" {
		if cb, err := uuid.Parse(d.CreatedBy); err == nil {
			u.CreatedBy = &cb
		}
	}
	if d.ModifiedBy != "" {
		if mb, err := uuid.Parse(d.ModifiedBy); err == nil {
			u.ModifiedBy = &mb
		}
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	var doc userDoc
	err := r.col.FindOne(ctx, bson.M{"email": email}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, authdomain.ErrUserNotFound
		}
		return nil, err
	}
	return doc.toDomain()
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	var doc userDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, authdomain.ErrUserNotFound
		}
		return nil, err
	}
	return doc.toDomain()
}

var _ authdomain.UserRepository = (*UserRepository)(nil)

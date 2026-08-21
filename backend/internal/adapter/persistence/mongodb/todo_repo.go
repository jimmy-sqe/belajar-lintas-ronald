package mongodb

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// TodoRepository is a MongoDB-backed implementation of todo.Repository.
type TodoRepository struct {
	col *mongo.Collection
}

func NewTodoRepository(db *mongo.Database) *TodoRepository {
	return &TodoRepository{col: db.Collection("todos")}
}

type todoDoc struct {
	ID          string     `bson:"_id"`
	Title       string     `bson:"title"`
	Description string     `bson:"description,omitempty"`
	DueDate     *time.Time `bson:"due_date,omitempty"`
	CreatedAt   time.Time  `bson:"created_at"`
	CreatedBy   string     `bson:"created_by"`
	ModifiedAt  time.Time  `bson:"modified_at"`
	ModifiedBy  string     `bson:"modified_by"`
}

func docFromDomain(t *todo.Todo) todoDoc {
	return todoDoc{
		ID: t.ID.String(), Title: t.Title, Description: t.Description, DueDate: t.DueDate,
		CreatedAt: t.CreatedAt, CreatedBy: t.CreatedBy.String(),
		ModifiedAt: t.ModifiedAt, ModifiedBy: t.ModifiedBy.String(),
	}
}

func (d todoDoc) toDomain() (*todo.Todo, error) {
	id, err := uuid.Parse(d.ID)
	if err != nil {
		return nil, err
	}
	cb, err := uuid.Parse(d.CreatedBy)
	if err != nil {
		return nil, err
	}
	mb, err := uuid.Parse(d.ModifiedBy)
	if err != nil {
		return nil, err
	}
	return &todo.Todo{
		ID: id, Title: d.Title, Description: d.Description, DueDate: d.DueDate,
		CreatedAt: d.CreatedAt, CreatedBy: cb,
		ModifiedAt: d.ModifiedAt, ModifiedBy: mb,
	}, nil
}

func (r *TodoRepository) Create(ctx context.Context, t *todo.Todo) error {
	_, err := r.col.InsertOne(ctx, docFromDomain(t))
	return err
}

func (r *TodoRepository) FindByID(ctx context.Context, id uuid.UUID) (*todo.Todo, error) {
	var doc todoDoc
	err := r.col.FindOne(ctx, bson.M{"_id": id.String()}).Decode(&doc)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, todo.ErrNotFound
		}
		return nil, err
	}
	return doc.toDomain()
}

func (r *TodoRepository) List(ctx context.Context, filter todo.ListFilter) ([]todo.Todo, int64, error) {
	count, err := r.col.CountDocuments(ctx, bson.M{"created_by": filter.OwnerID.String()})
	if err != nil {
		return nil, 0, err
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetSkip(int64(filter.Offset)).
		SetLimit(int64(filter.Limit))
	cur, err := r.col.Find(ctx, bson.M{"created_by": filter.OwnerID.String()}, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cur.Close(ctx) //nolint:errcheck // best-effort cursor close
	var out []todo.Todo
	for cur.Next(ctx) {
		var doc todoDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, 0, err
		}
		t, err := doc.toDomain()
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *t)
	}
	return out, count, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *todo.Todo) error {
	res, err := r.col.ReplaceOne(ctx, bson.M{"_id": t.ID.String()}, docFromDomain(t))
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return todo.ErrNotFound
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.col.DeleteOne(ctx, bson.M{"_id": id.String()})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return todo.ErrNotFound
	}
	return nil
}

var _ todo.Repository = (*TodoRepository)(nil)

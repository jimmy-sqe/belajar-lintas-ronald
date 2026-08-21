package mysql

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
)

// TodoRepository is a MySQL-backed implementation of todo.Repository.
type TodoRepository struct {
	db *sqlx.DB
}

func NewTodoRepository(db *sqlx.DB) *TodoRepository {
	return &TodoRepository{db: db}
}

type todoRow struct {
	ID          string         `db:"id"`
	Title       string         `db:"title"`
	Description sql.NullString `db:"description"`
	DueDate     sql.NullTime   `db:"due_date"`
	CreatedAt   sql.NullTime   `db:"created_at"`
	CreatedBy   string         `db:"created_by"`
	ModifiedAt  sql.NullTime   `db:"modified_at"`
	ModifiedBy  string         `db:"modified_by"`
}

func (r todoRow) toDomain() (*todo.Todo, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, err
	}
	createdBy, err := uuid.Parse(r.CreatedBy)
	if err != nil {
		return nil, err
	}
	modifiedBy, err := uuid.Parse(r.ModifiedBy)
	if err != nil {
		return nil, err
	}
	t := &todo.Todo{
		ID:         id,
		Title:      r.Title,
		CreatedAt:  r.CreatedAt.Time,
		CreatedBy:  createdBy,
		ModifiedAt: r.ModifiedAt.Time,
		ModifiedBy: modifiedBy,
	}
	if r.Description.Valid {
		t.Description = r.Description.String
	}
	if r.DueDate.Valid {
		dd := r.DueDate.Time
		t.DueDate = &dd
	}
	return t, nil
}

func (r *TodoRepository) Create(ctx context.Context, t *todo.Todo) error {
	var due sql.NullTime
	if t.DueDate != nil {
		due = sql.NullTime{Time: *t.DueDate, Valid: true}
	}
	var desc sql.NullString
	if t.Description != "" {
		desc = sql.NullString{String: t.Description, Valid: true}
	}
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO todos (id, title, description, due_date, created_at, created_by, modified_at, modified_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		t.ID.String(), t.Title, desc, due,
		t.CreatedAt, t.CreatedBy.String(),
		t.ModifiedAt, t.ModifiedBy.String())
	return err
}

func (r *TodoRepository) FindByID(ctx context.Context, id uuid.UUID) (*todo.Todo, error) {
	var row todoRow
	err := r.db.GetContext(ctx, &row,
		`SELECT id, title, description, due_date, created_at, created_by, modified_at, modified_by
		 FROM todos WHERE id = ? LIMIT 1`,
		id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, todo.ErrNotFound
		}
		return nil, err
	}
	return row.toDomain()
}

func (r *TodoRepository) List(ctx context.Context, filter todo.ListFilter) ([]todo.Todo, int64, error) {
	var rows []todoRow
	err := r.db.SelectContext(ctx, &rows,
		`SELECT id, title, description, due_date, created_at, created_by, modified_at, modified_by
		 FROM todos WHERE created_by = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		filter.OwnerID.String(), filter.Limit, filter.Offset)
	if err != nil {
		return nil, 0, err
	}

	var total int64
	err = r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM todos WHERE created_by = ?`,
		filter.OwnerID.String())
	if err != nil {
		return nil, 0, err
	}

	out := make([]todo.Todo, 0, len(rows))
	for _, row := range rows {
		t, err := row.toDomain()
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *t)
	}
	return out, total, nil
}

func (r *TodoRepository) Update(ctx context.Context, t *todo.Todo) error {
	var due sql.NullTime
	if t.DueDate != nil {
		due = sql.NullTime{Time: *t.DueDate, Valid: true}
	}
	var desc sql.NullString
	if t.Description != "" {
		desc = sql.NullString{String: t.Description, Valid: true}
	}
	res, err := r.db.ExecContext(ctx,
		`UPDATE todos SET title = ?, description = ?, due_date = ?, modified_at = ?, modified_by = ?
		 WHERE id = ?`,
		t.Title, desc, due, t.ModifiedAt, t.ModifiedBy.String(), t.ID.String())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return todo.ErrNotFound
	}
	return nil
}

func (r *TodoRepository) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := r.db.ExecContext(ctx,
		`DELETE FROM todos WHERE id = ?`, id.String())
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return todo.ErrNotFound
	}
	return nil
}

var _ todo.Repository = (*TodoRepository)(nil)

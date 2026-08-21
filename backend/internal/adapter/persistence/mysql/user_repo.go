package mysql

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

type userRow struct {
	ID           string         `db:"id"`
	Email        string         `db:"email"`
	PasswordHash string         `db:"password_hash"`
	Name         sql.NullString `db:"name"`
	CreatedAt    sql.NullTime   `db:"created_at"`
	CreatedBy    sql.NullString `db:"created_by"`
	ModifiedAt   sql.NullTime   `db:"modified_at"`
	ModifiedBy   sql.NullString `db:"modified_by"`
}

func (r userRow) toDomain() (*authdomain.User, error) {
	id, err := uuid.Parse(r.ID)
	if err != nil {
		return nil, err
	}
	u := &authdomain.User{
		ID: id, Email: r.Email, PasswordHash: r.PasswordHash,
		Name: r.Name.String, CreatedAt: r.CreatedAt.Time, ModifiedAt: r.ModifiedAt.Time,
	}
	if r.CreatedBy.Valid {
		if cb, err := uuid.Parse(r.CreatedBy.String); err == nil {
			u.CreatedBy = &cb
		}
	}
	if r.ModifiedBy.Valid {
		if mb, err := uuid.Parse(r.ModifiedBy.String); err == nil {
			u.ModifiedBy = &mb
		}
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*authdomain.User, error) {
	var row userRow
	err := r.db.GetContext(ctx, &row,
		`SELECT id, email, password_hash, name, created_at, created_by, modified_at, modified_by
		 FROM users WHERE LOWER(email) = LOWER(?) LIMIT 1`,
		strings.ToLower(email))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authdomain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain()
}

func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*authdomain.User, error) {
	var row userRow
	err := r.db.GetContext(ctx, &row,
		`SELECT id, email, password_hash, name, created_at, created_by, modified_at, modified_by
		 FROM users WHERE id = ? LIMIT 1`,
		id.String())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, authdomain.ErrUserNotFound
		}
		return nil, err
	}
	return row.toDomain()
}

var _ authdomain.UserRepository = (*UserRepository)(nil)

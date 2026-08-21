package todo

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, input CreateInput) (*Todo, error)
	FindByID(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*Todo, error)
	List(ctx context.Context, filter ListFilter) ([]Todo, int64, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*Todo, error)
	Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error
}

type service struct {
	repo  Repository
	cache Cache
}

func NewService(repo Repository, cache Cache) Service {
	return &service{repo: repo, cache: cache}
}

func (s *service) Create(ctx context.Context, input CreateInput) (*Todo, error) {
	now := time.Now().UTC()
	t := &Todo{
		ID:          uuid.New(),
		Title:       input.Title,
		Description: input.Description,
		DueDate:     input.DueDate,
		CreatedAt:   now,
		CreatedBy:   input.CreatedBy,
		ModifiedAt:  now,
		ModifiedBy:  input.CreatedBy,
	}
	if err := s.repo.Create(ctx, t); err != nil {
		return nil, err
	}
	_ = s.cache.Set(ctx, t)
	return t, nil
}

func (s *service) FindByID(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) (*Todo, error) {
	if cached, ok := s.cache.Get(ctx, id); ok {
		if cached.CreatedBy != ownerID {
			return nil, ErrNotFound
		}
		return cached, nil
	}
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.CreatedBy != ownerID {
		return nil, ErrNotFound
	}
	_ = s.cache.Set(ctx, t)
	return t, nil
}

func (s *service) List(ctx context.Context, filter ListFilter) ([]Todo, int64, error) {
	return s.repo.List(ctx, filter)
}

func (s *service) Update(ctx context.Context, id uuid.UUID, input UpdateInput) (*Todo, error) {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if t.CreatedBy != input.ModifiedBy {
		return nil, ErrNotFound
	}
	if input.Title != nil {
		t.Title = *input.Title
	}
	if input.Description != nil {
		t.Description = *input.Description
	}
	if input.DueDate != nil {
		t.DueDate = *input.DueDate
	}
	t.ModifiedAt = time.Now().UTC()
	t.ModifiedBy = input.ModifiedBy
	if err := s.repo.Update(ctx, t); err != nil {
		return nil, err
	}
	_ = s.cache.Delete(ctx, id)
	return t, nil
}

func (s *service) Delete(ctx context.Context, id uuid.UUID, ownerID uuid.UUID) error {
	t, err := s.repo.FindByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}
		return err
	}
	if t.CreatedBy != ownerID {
		return ErrNotFound
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}
	_ = s.cache.Delete(ctx, id)
	return nil
}

package todo_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo"
	todomock "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/todo/mock"
)

func tctx() context.Context { return context.Background() }

func TestService_Create_Success(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()

	repo.On("Create", mock.Anything, mock.MatchedBy(func(td *todo.Todo) bool {
		return td.ID != uuid.Nil &&
			td.Title == "Buy milk" &&
			td.CreatedBy == ownerID &&
			td.ModifiedBy == ownerID &&
			!td.CreatedAt.IsZero() &&
			!td.ModifiedAt.IsZero()
	})).Return(nil)
	cache.On("Set", mock.Anything, mock.Anything).Return(nil)

	svc := todo.NewService(repo, cache)
	got, err := svc.Create(tctx(), todo.CreateInput{
		Title:     "Buy milk",
		CreatedBy: ownerID,
	})

	assert.NoError(t, err)
	assert.NotEqual(t, uuid.Nil, got.ID)
	assert.Equal(t, ownerID, got.CreatedBy)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestService_FindByID_Success(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	id := uuid.New()
	stored := &todo.Todo{
		ID:        id,
		Title:     "Buy milk",
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	cache.On("Get", mock.Anything, id).Return((*todo.Todo)(nil), false)
	repo.On("FindByID", mock.Anything, id).Return(stored, nil)
	cache.On("Set", mock.Anything, stored).Return(nil)

	svc := todo.NewService(repo, cache)
	got, err := svc.FindByID(tctx(), id, ownerID)

	assert.NoError(t, err)
	assert.Equal(t, stored, got)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestService_FindByID_NotFoundForeignOwner(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	foreignID := uuid.New()
	id := uuid.New()
	stored := &todo.Todo{
		ID:        id,
		Title:     "Buy milk",
		CreatedBy: foreignID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	cache.On("Get", mock.Anything, id).Return((*todo.Todo)(nil), false)
	repo.On("FindByID", mock.Anything, id).Return(stored, nil)

	svc := todo.NewService(repo, cache)
	got, err := svc.FindByID(tctx(), id, ownerID)

	assert.ErrorIs(t, err, todo.ErrNotFound)
	assert.Nil(t, got)
}

func TestService_List_FiltersByOwner(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	filter := todo.ListFilter{OwnerID: ownerID, Limit: 10, Offset: 0}
	repo.On("List", mock.Anything, filter).Return([]todo.Todo{{ID: uuid.New(), CreatedBy: ownerID}}, int64(1), nil)

	svc := todo.NewService(repo, cache)
	todos, total, err := svc.List(tctx(), filter)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), total)
	assert.Len(t, todos, 1)
	repo.AssertExpectations(t)
}

func TestService_Update_Success(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	id := uuid.New()
	existing := &todo.Todo{
		ID:        id,
		Title:     "old",
		CreatedBy: ownerID,
		CreatedAt: time.Now().Add(-time.Hour),
		ModifiedAt: time.Now().Add(-time.Hour),
	}
	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("Update", mock.Anything, mock.MatchedBy(func(td *todo.Todo) bool {
		return td.Title == "new" && td.ModifiedBy == ownerID
	})).Return(nil)
	cache.On("Delete", mock.Anything, id).Return(nil)

	svc := todo.NewService(repo, cache)
	newTitle := "new"
	got, err := svc.Update(tctx(), id, todo.UpdateInput{
		Title:      &newTitle,
		ModifiedBy: ownerID,
	})

	assert.NoError(t, err)
	assert.Equal(t, "new", got.Title)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestService_Update_ForeignOwnerReturnsNotFound(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	foreignID := uuid.New()
	id := uuid.New()
	existing := &todo.Todo{
		ID:        id,
		Title:     "old",
		CreatedBy: foreignID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	repo.On("FindByID", mock.Anything, id).Return(existing, nil)

	svc := todo.NewService(repo, cache)
	newTitle := "new"
	got, err := svc.Update(tctx(), id, todo.UpdateInput{
		Title:      &newTitle,
		ModifiedBy: ownerID,
	})

	assert.ErrorIs(t, err, todo.ErrNotFound)
	assert.Nil(t, got)
}

func TestService_Delete_Success(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	id := uuid.New()
	existing := &todo.Todo{
		ID:        id,
		Title:     "Buy milk",
		CreatedBy: ownerID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	repo.On("FindByID", mock.Anything, id).Return(existing, nil)
	repo.On("Delete", mock.Anything, id).Return(nil)
	cache.On("Delete", mock.Anything, id).Return(nil)

	svc := todo.NewService(repo, cache)
	err := svc.Delete(tctx(), id, ownerID)

	assert.NoError(t, err)
	repo.AssertExpectations(t)
	cache.AssertExpectations(t)
}

func TestService_Delete_ForeignOwnerReturnsNotFound(t *testing.T) {
	repo := new(todomock.Repository)
	cache := new(todomock.Cache)
	ownerID := uuid.New()
	foreignID := uuid.New()
	id := uuid.New()
	existing := &todo.Todo{
		ID:        id,
		Title:     "Buy milk",
		CreatedBy: foreignID,
		CreatedAt: time.Now(),
		ModifiedAt: time.Now(),
	}
	repo.On("FindByID", mock.Anything, id).Return(existing, nil)

	svc := todo.NewService(repo, cache)
	err := svc.Delete(tctx(), id, ownerID)

	assert.ErrorIs(t, err, todo.ErrNotFound)
}

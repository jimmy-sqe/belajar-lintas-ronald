// Package redis implements SessionStore using a Redis backend.
package redis

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"

	authdomain "github.com/S-Quantum-Engine/belajar-lintas-ronald/backend/internal/domain/auth"
)

const keyPrefix = "refresh:"

// SessionStore stores refresh-session records in Redis with TTL matching
// the session's expiry.
type SessionStore struct {
	client *redis.Client
}

// NewSessionStore connects to Redis at the given URL and returns a store.
func NewSessionStore(url string) (*SessionStore, error) {
	opt, err := redis.ParseURL(url)
	if err != nil {
		return nil, err
	}
	return &SessionStore{client: redis.NewClient(opt)}, nil
}

// NewSessionStoreWithClient lets callers reuse an existing client
// (e.g. wired in the composition root).
func NewSessionStoreWithClient(client *redis.Client) *SessionStore {
	return &SessionStore{client: client}
}

type record struct {
	UserID    string `json:"user_id"`
	ExpiresAt int64  `json:"expires_at"`
}

// Put stores the session with TTL equal to ExpiresAt - now.
func (s *SessionStore) Put(ctx context.Context, sess authdomain.RefreshSession) error {
	ttl := time.Until(sess.ExpiresAt)
	if ttl <= 0 {
		return errors.New("session already expired")
	}
	data, err := json.Marshal(record{
		UserID:    sess.UserID.String(),
		ExpiresAt: sess.ExpiresAt.Unix(),
	})
	if err != nil {
		return err
	}
	return s.client.Set(ctx, keyPrefix+sess.TokenHash, data, ttl).Err()
}

// Get returns the session record, or nil if not found / expired.
func (s *SessionStore) Get(ctx context.Context, tokenHash string) (*authdomain.RefreshSession, error) {
	data, err := s.client.Get(ctx, keyPrefix+tokenHash).Bytes()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var rec record
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	id, err := uuid.Parse(rec.UserID)
	if err != nil {
		return nil, err
	}
	return &authdomain.RefreshSession{
		TokenHash: tokenHash,
		UserID:    id,
		ExpiresAt: time.Unix(rec.ExpiresAt, 0).UTC(),
	}, nil
}

// Delete removes the session record. Returns nil on missing key.
func (s *SessionStore) Delete(ctx context.Context, tokenHash string) error {
	return s.client.Del(ctx, keyPrefix+tokenHash).Err()
}

// Compile-time check.
var _ authdomain.SessionStore = (*SessionStore)(nil)

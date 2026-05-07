package token

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	TokenPrefix = "token:"
	TokenTTL    = 24 * time.Hour
	GraceTTL    = 5 * time.Minute
)

type Manager struct {
	rdb *redis.Client
}

func NewManager(rdb *redis.Client) *Manager {
	return &Manager{rdb: rdb}
}

func Generate(userID string) string {
	src := fmt.Sprintf("%s_%d_%d", userID, time.Now().UnixNano(), rand.Int63())
	hash := sha256.Sum256([]byte(src))
	return hex.EncodeToString(hash[:])[:32]
}

func (m *Manager) Store(ctx context.Context, token, userID string) error {
	key := TokenPrefix + token
	return m.rdb.Set(ctx, key, userID, TokenTTL).Err()
}

func (m *Manager) Validate(ctx context.Context, token string) (string, error) {
	key := TokenPrefix + token
	userID, err := m.rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", err
	}
	return userID, nil
}

func (m *Manager) Refresh(ctx context.Context, oldToken, userID string) (string, error) {
	oldKey := TokenPrefix + oldToken
	m.rdb.Expire(ctx, oldKey, GraceTTL)

	newToken := Generate(userID)
	return newToken, m.Store(ctx, newToken, userID)
}

func (m *Manager) Revoke(ctx context.Context, token string) error {
	return m.rdb.Del(ctx, TokenPrefix+token).Err()
}

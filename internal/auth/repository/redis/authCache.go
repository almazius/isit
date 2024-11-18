package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

const keyPrefixForSession = "auth:session:"

type AuthCache struct {
	rdb *redis.Client
}

func NewAuthCache(rdb *redis.Client) *AuthCache {
	return &AuthCache{
		rdb: rdb,
	}
}

func (c *AuthCache) CreateNewSession(ctx context.Context, fingerPrint, userAgent, sessionKey string) error {
	err := c.rdb.Set(ctx,
		genKeyPrefixForSession(fingerPrint, userAgent),
		sessionKey,
		1*time.Hour).Err()
	if err != nil {
		slog.Error("failed set session to cache", "error", err)
		return err
	}

	return nil
}

func (c *AuthCache) ValidateSession(ctx context.Context, fingerPrint, userAgent string) (string, error) {
	res := c.rdb.Get(ctx, genKeyPrefixForSession(fingerPrint, userAgent))
	if res.Err() != nil {
		slog.Error("failed get session from cache", "error", res.Err())
		return "", res.Err()
	}

	return res.Val(), nil
}

func (c *AuthCache) DeleteNewSession(ctx context.Context, fingerPrint, userAgent string) error {
	err := c.rdb.Del(ctx, genKeyPrefixForSession(fingerPrint, userAgent)).Err()
	if err != nil {
		slog.Error("failed delete session", "error", err)
		return err
	}

	return nil
}
func genKeyPrefixForSession(fingerPrint, userAgent string) string {
	return fmt.Sprintf("%s%s:%s", keyPrefixForSession, fingerPrint, userAgent)
}

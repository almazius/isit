package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"isit/config"
	"time"
)

func Init(config *config.Config) (*redis.Client, error) {
	ctx, cansel := context.WithTimeout(context.Background(), time.Second*5)
	defer cansel()
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Post),
		Username: "default",
		Password: "default",
		DB:       config.Redis.DB,
	})
	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}

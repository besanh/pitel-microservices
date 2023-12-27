package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
)

type (
	IRedisCache interface {
		Set(key string, value any, ttl time.Duration)
		Get(key string) any
	}
	RedisCache struct {
		client *redis.Client
	}
)

func NewRedisCache(client *redis.Client) IRedisCache {
	return &RedisCache{
		client: client,
	}
}

const (
	REDIS_KEEP_TTL = redis.KeepTTL
)

func (r *RedisCache) Set(key string, value any, ttl time.Duration) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	val, _ := util.ParseAnyToString(value)
	if _, err := r.client.Set(ctx, key, val, ttl).Result(); err != nil {
		log.Error(err)
	}
}

func (r *RedisCache) Get(key string) any {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	val, err := r.client.Get(ctx, key).Result()
	if err != nil {
		log.Error(err)
		return nil
	}
	return val
}

package cache

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
)

type (
	IRedisCache interface {
		Set(key string, value any, ttl time.Duration)
		SetTTL(key string, value any, t time.Duration) (string, error)
		Get(key string) any
		IsExisted(key string) (bool, error)
		IsHExisted(list, key string) (bool, error)
		HGet(list, key string) (string, error)
		HGetAll(list string) (map[string]string, error)
		HSet(key string, values []any) error
		HMGet(key string, fields ...string) ([]any, error)
		HMSet(key string, values ...any) error
		HMDel(key string, fields ...string) error
		FLUSHALL() any
		Del(key []string) error
		HDel(key string, fields ...string) error
		GetKeysPattern(pattern string) ([]string, error)
		Close()
	}
	RedisCache struct {
		client *redis.Client
	}
)

var RCache IRedisCache

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

func (r *RedisCache) SetTTL(key string, value any, t time.Duration) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.Set(ctx, key, value, t).Result()
	return ret, err
}

func (r *RedisCache) IsExisted(key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := r.client.Exists(ctx, key).Result()
	if res == 0 || err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisCache) IsHExisted(list, key string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	res, err := r.client.HExists(ctx, list, key).Result()
	if !res || err != nil {
		return false, err
	}
	return true, nil
}

func (r *RedisCache) HGet(list, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.HGet(ctx, list, key).Result()
	return ret, err
}

func (r *RedisCache) HGetAll(list string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.HGetAll(ctx, list).Result()
	return ret, err
}

func (r *RedisCache) HSet(key string, values []any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	_, err := r.client.HSet(ctx, key, values...).Result()
	return err
}

func (r *RedisCache) Del(key []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := r.client.Del(ctx, key...).Err()
	return err
}

func (r *RedisCache) HMSet(key string, values ...any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.HMSet(ctx, key, values...).Result()
	if err != nil {
		return err
	}
	if !ret {
		err = errors.New("HashMap Set failed")
	}
	return err
}

func (r *RedisCache) HMDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := r.client.HDel(ctx, key, fields...).Err()
	return err
}

func (r *RedisCache) FLUSHALL() any {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret := r.client.FlushAll(ctx)
	return ret
}

func (r *RedisCache) HMGet(key string, fields ...string) ([]any, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.HMGet(ctx, key, fields...).Result()
	return ret, err
}

func (r *RedisCache) HDel(key string, fields ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := r.client.HDel(ctx, key, fields...).Err()
	return err
}

func (r *RedisCache) GetKeysPattern(pattern string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	ret, err := r.client.Keys(ctx, pattern).Result()
	return ret, err
}

func (r *RedisCache) Close() {
	r.client.Close()
}

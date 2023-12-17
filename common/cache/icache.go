package cache

import (
	"time"
)

type IMemCache interface {
	Set(key string, value any) error
	SetTTL(key string, value any, t time.Duration) error
	Get(key string) (any, error)
	Del(key string) error
	Close()
}

type IRedisCache interface {
	Set(key string, value any) error
	SetTTL(key string, value any, t time.Duration) error
	Get(key string) (string, error)
	Del(key string) error
	Close()
}

var RCache IRedisCache
var MCache IMemCache

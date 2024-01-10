package cache

import (
	"context"
	"time"

	"github.com/jellydator/ttlcache/v3"
	"github.com/tel4vn/fins-microservices/common/log"
	"github.com/tel4vn/fins-microservices/common/util"
)

type (
	IMemCache interface {
		Set(key string, value any, ttl time.Duration)
		Get(key string) any
		Del(key string)
		Close()
	}
	MemCache struct {
		*ttlcache.Cache[string, any]
	}
)

const (
	DEFAULT_TTL = ttlcache.DefaultTTL
)

var MCache IMemCache

func NewMemCache() IMemCache {
	service := &MemCache{}
	cache := ttlcache.New(
		ttlcache.WithTTL[string, any](30 * time.Minute),
	)
	cache.OnInsertion(func(ctx context.Context, item *ttlcache.Item[string, any]) {
		log.Infof("memcache: inserted %s, expires at %s", item.Key(), item.ExpiresAt())
	})
	cache.OnEviction(func(ctx context.Context, reason ttlcache.EvictionReason, item *ttlcache.Item[string, any]) {
		if reason == ttlcache.EvictionReasonCapacityReached {
			val, _ := util.ParseAnyToString(item.Value())
			log.Infof("memcache: removed %s, value: %v", item.Key(), val)
		}
	})
	service.Cache = cache
	return service
}

func (s *MemCache) Set(key string, value any, ttl time.Duration) {
	s.Cache.Set(key, value, ttl)
}

func (s *MemCache) Get(key string) any {
	val := s.Cache.Get(key)
	if val == nil {
		return nil
	}
	return val.Value()
}

func (s *MemCache) Del(key string) {
	s.Cache.Delete(key)
}

func (s *MemCache) Close() {
	s.Cache.Stop()
}

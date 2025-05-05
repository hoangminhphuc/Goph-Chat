package cache

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/go-redis/cache/v9"
	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"
)

type Cache interface {
	Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
}

var (
	defaultCache *TwoLevelCache
	defaultTTL   = 5 * time.Minute
	defaultLRUSize = 1000
)
// Implementing a LRU in-memory cache and Redis cache.
type TwoLevelCache struct {
	local    *lru.Cache[string, interface{}]
	remote   *cache.Cache
	loader   singleflight.Group
	ttl      time.Duration
}

// NewTwoLevelCache creates a TwoLevelCache.
// lruSize is max entries in local cache, ttl is default expiration.
func NewTwoLevelCache(rdb *redis.Client) (*TwoLevelCache, error) {
	// create LRU cache (Level 1)
	localCache, err := lru.New[string, interface{}](defaultLRUSize)
	if err != nil {
		return nil, err
	}
	// create Redis-backed cache without internal micro-cache
	r := cache.New(&cache.Options{
		Redis:      rdb,
		// Only use when Redis latency spikes under heavy concurrent bursts to 
		// absorb extreme bursts of traffic.
		LocalCache: nil, 
	})

	return &TwoLevelCache{
		local:  localCache,
		remote: r,
		ttl:    defaultTTL,
	}, nil
}

func InitDefaultCache(rdb *redis.Client) error {
	cache, err := NewTwoLevelCache(rdb)
	if err != nil {
		return err
	}
	defaultCache = cache
	return nil
}

func GetCache() *TwoLevelCache {
	if defaultCache == nil {
		panic("cache not initialized")
	}
	return defaultCache
}



// Set writes to both local and Redis caches.
func (c *TwoLevelCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// write to Redis (Level 2)
	if err := c.remote.Set(&cache.Item{Ctx: ctx, Key: key, Value: value, TTL: ttl}); err != nil {
		return err
	}
	// write to local LRU (Level 1)
	c.local.Add(key, value)
	return nil
}

// Get attempts local cache, then remote with singleflight to prevent stampede.
func (c *TwoLevelCache) Get(ctx context.Context, key string, dest interface{}) error {
	// Local LRU lookup
	if v, ok := c.local.Get(key); ok {
			switch d := dest.(type) {
			case *string:
					if str, ok := v.(string); ok {
							*d = str
							return nil
					}
			default:
					// dest is pointer to struct or value type
					rv := reflect.ValueOf(dest)
					vp := reflect.ValueOf(v)
					if rv.Kind() == reflect.Ptr && vp.Kind() == reflect.Ptr {
							rv.Elem().Set(vp.Elem())
							return nil
					}
			}
	}

	// Load from Redis (Level 2) or fallback to DB
	val, err, _ := c.loader.Do(key, func() (interface{}, error) {
			var err error
			if err := c.remote.Get(ctx, key, dest); err == nil {
					c.local.Add(key, dest)
					return dest, nil
			}
			return nil, err
	})

	if err != nil {
			return err
	}

	// Assign result into dest
	rv := reflect.ValueOf(dest)
	vp := reflect.ValueOf(val)
	if rv.Kind() == reflect.Ptr && vp.Kind() == reflect.Ptr {
			rv.Elem().Set(vp.Elem())
			return nil
	}
	return fmt.Errorf("cache error")
}

// Delete removes from both caches.
func (c *TwoLevelCache) Delete(ctx context.Context, key string) error {
	c.local.Remove(key)
	return c.remote.Delete(ctx, key)
}
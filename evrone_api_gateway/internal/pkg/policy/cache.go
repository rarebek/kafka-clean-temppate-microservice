package policy

import (
	"context"
	"time"

	"github.com/casbin/casbin/v2/persist/cache"
	"github.com/go-redis/redis/v8"
)

type casbinCache struct {
	db *redis.Client
}

func NewCache(db *redis.Client) *casbinCache {
	return &casbinCache{
		db: db,
	}
}

func (c *casbinCache) Set(key string, value bool, extra ...interface{}) error {
	var expiration time.Duration = 0
	c.db.Set(context.Background(), key, value, expiration)
	return nil
}

func (c *casbinCache) Get(key string) (bool, error) {
	res, err := c.db.Get(context.Background(), key).Bool()
	if err != nil {
		return false, cache.ErrNoSuchKey
	}
	return res, nil
}

func (c *casbinCache) Delete(key string) error {
	if res := c.db.Del(context.Background(), key); res.Err() != nil {
		return cache.ErrNoSuchKey
	}
	return nil
}

func (c *casbinCache) Clear() error {
	return c.db.FlushAll(context.Background()).Err()
}

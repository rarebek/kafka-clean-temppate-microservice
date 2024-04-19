package redis

//
//import (
//	"context"
//	"encoding/json"
//	"time"
//
//	"go.opentelemetry.io/otel/attribute"
//
//	"evrone_api_gateway/internal/pkg/redis"
//)
//
//type Cache interface {
//	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
//	Get(ctx context.Context, key string) ([]byte, error)
//	Del(ctx context.Context, key string) error
//}
//
//func NewCache(rdb *redis.RedisDB) *cache {
//	return &cache{
//		rdb: rdb,
//	}
//}
//
//type cache struct {
//	rdb *redis.RedisDB
//}
//
//func (c *cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
//	// tracing
//	ctx, span := otlp_pkg.Start(ctx, "cecheService", "CasheRepoSet")
//	defer span.End()
//	byteData, err := json.Marshal(value)
//	if err != nil {
//		return err
//	}
//	err = c.rdb.Client.Set(ctx, key, string(byteData), expiration).Err()
//	if err != nil {
//		return err
//	}
//	return nil
//}
//
//func (c *cache) Get(ctx context.Context, key string) ([]byte, error) {
//	// tracing
//	ctx, span := otlp_pkg.Start(ctx, "cecheService", "CasheRepoGet")
//	defer span.End()
//
//	data, err := c.rdb.Client.Get(ctx, key).Result()
//	if err != nil {
//		return nil, err
//	}
//
//	span.SetAttributes(
//		attribute.Key("data").String(data),
//	)
//
//	return []byte(data), nil
//}
//
//func (c *cache) Del(ctx context.Context, key string) error {
//	// tracing
//	ctx, span := otlp_pkg.Start(ctx, "cecheService", "CasheRepoDel")
//	defer span.End()
//
//	err := c.rdb.Client.Del(ctx, key).Err()
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

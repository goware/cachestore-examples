package cachestore_e2e_test

import (
	"context"
	"time"

	cachestore "github.com/goware/cachestore2"
	"github.com/redis/go-redis/v9"
)

func redisFlushAll() {
	redisClient := redis.NewClient(&redis.Options{Addr: "localhost:6379", DB: 9})
	redisClient.FlushAll(context.Background())
}

type metricsStore[T any] struct {
	store            cachestore.Store[T]
	telemetryCounter uint64
}

var _ cachestore.Store[any] = &metricsStore[any]{}

func NewMetricsStore[V any](store cachestore.Store[V]) cachestore.Store[V] {
	return &metricsStore[V]{
		store: store,
	}
}

func (s *metricsStore[V]) Name() string {
	return "cacheWithTelemetry"
}

func (s *metricsStore[V]) Options() cachestore.StoreOptions {
	s.telemetryCounter++
	return cachestore.StoreOptions{}
}

func (s *metricsStore[V]) Exists(ctx context.Context, key string) (bool, error) {
	s.telemetryCounter++
	return s.store.Exists(ctx, key)
}

func (s *metricsStore[V]) Set(ctx context.Context, key string, value V) error {
	s.telemetryCounter++
	return s.store.Set(ctx, key, value)
}

func (s *metricsStore[V]) SetEx(ctx context.Context, key string, value V, ttl time.Duration) error {
	s.telemetryCounter++
	return s.store.SetEx(ctx, key, value, ttl)
}

func (s *metricsStore[V]) BatchSet(ctx context.Context, keys []string, values []V) error {
	s.telemetryCounter++
	return s.store.BatchSet(ctx, keys, values)
}

func (s *metricsStore[V]) BatchSetEx(ctx context.Context, keys []string, values []V, ttl time.Duration) error {
	s.telemetryCounter++
	return s.store.BatchSetEx(ctx, keys, values, ttl)
}

func (s *metricsStore[V]) Get(ctx context.Context, key string) (V, bool, error) {
	s.telemetryCounter++
	return s.store.Get(ctx, key)
}

func (s *metricsStore[V]) BatchGet(ctx context.Context, keys []string) ([]V, []bool, error) {
	s.telemetryCounter++
	return s.store.BatchGet(ctx, keys)
}

func (s *metricsStore[V]) Delete(ctx context.Context, key string) error {
	s.telemetryCounter++
	return s.store.Delete(ctx, key)
}

func (s *metricsStore[V]) DeletePrefix(ctx context.Context, keyPrefix string) error {
	s.telemetryCounter++
	return s.store.DeletePrefix(ctx, keyPrefix)
}

func (s *metricsStore[V]) ClearAll(ctx context.Context) error {
	s.telemetryCounter++
	return s.store.ClearAll(ctx)
}

func (s *metricsStore[V]) GetOrSetWithLock(ctx context.Context, key string, getter func(context.Context, string) (V, error)) (V, error) {
	s.telemetryCounter++
	return s.store.GetOrSetWithLock(ctx, key, getter)
}

func (s *metricsStore[V]) GetOrSetWithLockEx(ctx context.Context, key string, getter func(context.Context, string) (V, error), ttl time.Duration) (V, error) {
	s.telemetryCounter++
	return s.store.GetOrSetWithLockEx(ctx, key, getter, ttl)
}

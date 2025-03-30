package cachestore_e2e_test

import (
	"context"
	"testing"

	"github.com/goware/cachestore"
	memcache "github.com/goware/cachestore-mem"
	rediscache "github.com/goware/cachestore-redis"
	"github.com/stretchr/testify/require"
)

func TestCachestoreE2E(t *testing.T) {

	t.Run("memcache direct", func(t *testing.T) {
		mem, err := memcache.NewCacheWithSize[string](10)
		require.NoError(t, err)

		ctx := context.Background()

		ok, err := mem.Exists(ctx, "foo")
		require.NoError(t, err)
		require.False(t, ok)

		err = mem.Set(ctx, "foo", "bar")
		require.NoError(t, err)

		ok, err = mem.Exists(ctx, "foo")
		require.NoError(t, err)
		require.True(t, ok)

		val, ok, err := mem.Get(ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, "bar", val)
		require.True(t, ok)

		// TODO: a few more tests later..
	})

	t.Run("redis direct", func(t *testing.T) {
		redisFlushAll()

		red, err := rediscache.NewCache[string](&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		ctx := context.Background()

		ok, err := red.Exists(ctx, "foo")
		require.NoError(t, err)
		require.False(t, ok)

		err = red.Set(ctx, "foo", "bar")
		require.NoError(t, err)

		ok, err = red.Exists(ctx, "foo")
		require.NoError(t, err)
		require.True(t, ok)

		val, ok, err := red.Get(ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, "bar", val)
		require.True(t, ok)
	})

	t.Run("memcache backend", func(t *testing.T) {
		backend, err := memcache.NewBackend(10)
		require.NoError(t, err)

		mem := cachestore.OpenStore[string](backend)

		ctx := context.Background()

		ok, err := mem.Exists(ctx, "foo")
		require.NoError(t, err)
		require.False(t, ok)

		err = mem.Set(ctx, "foo", "bar")
		require.NoError(t, err)

		ok, err = mem.Exists(ctx, "foo")
		require.NoError(t, err)
		require.True(t, ok)

		val, ok, err := mem.Get(ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, "bar", val)
		require.True(t, ok)

		//--

		mem2 := cachestore.OpenStore[int](backend)
		mem2.Set(ctx, "age", 123)
		val2, ok, err := mem2.Get(ctx, "age")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 123, val2)

		//--

		mem3 := cachestore.OpenStore[apiResponse](backend)
		mem3.Set(ctx, "resp1", apiResponse{
			Status:  200,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    []byte(`{"message": "Hello, World!"}`),
		})
		val3, ok, err := mem3.Get(ctx, "resp1")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 200, val3.Status)
		require.Equal(t, "application/json", val3.Headers["Content-Type"])
		require.Equal(t, []byte(`{"message": "Hello, World!"}`), val3.Body)
	})

	t.Run("redis backend", func(t *testing.T) {
		redisFlushAll()
		backend, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		red := cachestore.OpenStore[string](backend)

		ctx := context.Background()

		ok, err := red.Exists(ctx, "foo")
		require.NoError(t, err)
		require.False(t, ok)

		err = red.Set(ctx, "foo", "bar")
		require.NoError(t, err)

		ok, err = red.Exists(ctx, "foo")
		require.NoError(t, err)
		require.True(t, ok)

		val, ok, err := red.Get(ctx, "foo")
		require.NoError(t, err)
		require.Equal(t, "bar", val)
		require.True(t, ok)

		//--

		red2 := cachestore.OpenStore[int](backend)
		red2.Set(ctx, "age", 123)
		val2, ok, err := red2.Get(ctx, "age")
		require.NoError(t, err)
		require.Equal(t, 123, val2)
		require.True(t, ok)

		//--

		red3 := cachestore.OpenStore[apiResponse](backend)
		red3.Set(ctx, "resp1", apiResponse{
			Status:  200,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    []byte(`{"message": "Hello, World!"}`),
		})
		val3, ok, err := red3.Get(ctx, "resp1")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, 200, val3.Status)
		require.Equal(t, "application/json", val3.Headers["Content-Type"])
		require.Equal(t, []byte(`{"message": "Hello, World!"}`), val3.Body)

	})

}

type apiResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

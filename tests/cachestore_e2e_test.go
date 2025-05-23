package cachestore_e2e_test

import (
	"context"
	"math/big"
	"testing"

	memcache "github.com/goware/cachestore-mem"
	rediscache "github.com/goware/cachestore-redis"
	cachestore "github.com/goware/cachestore2"
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

		// TODO: a few more tests later...
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

		ctx := context.Background()

		t.Run("type: string", func(t *testing.T) {
			red := cachestore.OpenStore[string](backend)
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

		t.Run("type: []byte", func(t *testing.T) {
			red := cachestore.OpenStore[[]byte](backend)

			ok, err := red.Exists(ctx, "name")
			require.NoError(t, err)
			require.False(t, ok)

			err = red.Set(ctx, "name", []byte{1, 2, 3})
			require.NoError(t, err)

			ok, err = red.Exists(ctx, "name")
			require.NoError(t, err)
			require.True(t, ok)

			val, ok, err := red.Get(ctx, "name")
			require.NoError(t, err)
			require.Equal(t, []byte{1, 2, 3}, val)
			require.True(t, ok)
		})

		t.Run("type: int", func(t *testing.T) {
			red := cachestore.OpenStore[int](backend)
			red.Set(ctx, "age", 123)
			val, ok, err := red.Get(ctx, "age")
			require.NoError(t, err)
			require.Equal(t, 123, val)
			require.True(t, ok)
		})

		t.Run("type: apiResponse", func(t *testing.T) {
			red := cachestore.OpenStore[apiResponse](backend)
			red.Set(ctx, "resp1", apiResponse{
				Status:  200,
				Headers: map[string]string{"Content-Type": "application/json"},
				Body:    []byte(`{"message": "Hello, World!"}`),
			})
			val, ok, err := red.Get(ctx, "resp1")
			require.NoError(t, err)
			require.True(t, ok)
			require.Equal(t, 200, val.Status)
			require.Equal(t, "application/json", val.Headers["Content-Type"])
			require.Equal(t, []byte(`{"message": "Hello, World!"}`), val.Body)
		})
	})

	t.Run("compose memcache and rediscache direct", func(t *testing.T) {
		redisFlushAll()

		mem, err := memcache.NewCacheWithSize[string](10)
		require.NoError(t, err)

		red, err := rediscache.NewCache[string](&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		composed, err := cachestore.ComposeStores(mem, red)
		require.NoError(t, err)

		ctx := context.Background()

		err = composed.Set(ctx, "a", "1")
		require.NoError(t, err)
		err = composed.Set(ctx, "b", "2")
		require.NoError(t, err)
		err = composed.Set(ctx, "c", "3")
		require.NoError(t, err)

		val, ok, err := composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "1", val)

		val, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "2", val)

		val, ok, err = composed.Get(ctx, "c")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "3", val)

		err = mem.Delete(ctx, "a")
		require.NoError(t, err)
		err = mem.Delete(ctx, "c")
		require.NoError(t, err)

		val, ok, err = composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "1", val)

		val, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "2", val)

		val, ok, err = composed.Get(ctx, "c")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "3", val)
	})

	t.Run("compose memcache and rediscache backend", func(t *testing.T) {
		backend, err := memcache.NewBackend(10)
		require.NoError(t, err)

		backend2, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		mem := cachestore.OpenStore[string](backend)
		red := cachestore.OpenStore[string](backend2)

		composed, err := cachestore.ComposeBackends[string](backend, backend2)
		require.NoError(t, err)

		ctx := context.Background()

		err = composed.Set(ctx, "a", "1")
		require.NoError(t, err)
		err = composed.Set(ctx, "b", "2")
		require.NoError(t, err)
		err = composed.Set(ctx, "c", "3")
		require.NoError(t, err)

		val, ok, err := composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "1", val)

		val, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "2", val)

		val, ok, err = composed.Get(ctx, "c")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "3", val)

		err = mem.Delete(ctx, "a")
		require.NoError(t, err)
		err = red.Delete(ctx, "c")
		require.NoError(t, err)

		val, ok, err = composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "1", val)

		val, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "2", val)

		val, ok, err = composed.Get(ctx, "c")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, "3", val)
	})

	t.Run("compose memcache and rediscache backend with struct values", func(t *testing.T) {
		backend, err := memcache.NewBackend(10)
		require.NoError(t, err)

		backend2, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		mem := cachestore.OpenStore[apiResponse](backend)
		red := cachestore.OpenStore[apiResponse](backend2)

		composed, err := cachestore.ComposeBackends[apiResponse](backend, backend2)
		require.NoError(t, err)

		ctx := context.Background()

		resp1 := apiResponse{
			Status:  200,
			Headers: map[string]string{"Content-Type": "application/json"},
			Body:    []byte(`{"message": "Hello, Alice!"}`),
		}

		resp2 := apiResponse{
			Status:  201,
			Headers: map[string]string{"Content-Type": "application/text"},
			Body:    []byte(`{"message": "Hello, Bob!"}`),
		}

		err = composed.Set(ctx, "a", resp1)
		require.NoError(t, err)
		err = composed.Set(ctx, "b", resp2)
		require.NoError(t, err)

		val1, ok, err := composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, resp1, val1)

		val2, ok, err := composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, resp2, val2)

		err = mem.Delete(ctx, "a")
		require.NoError(t, err)
		err = red.Delete(ctx, "c")
		require.NoError(t, err)

		val1, ok, err = composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, resp1, val1)

		resp2, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, resp2, val2)
	})

	t.Run("compose memcache and rediscache backend with *big.Int values", func(t *testing.T) {
		backend, err := memcache.NewBackend(10)
		require.NoError(t, err)

		backend2, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		mem := cachestore.OpenStore[*big.Int](backend)
		red := cachestore.OpenStore[*big.Int](backend2)

		composed, err := cachestore.ComposeBackends[*big.Int](backend, backend2)
		require.NoError(t, err)

		ctx := context.Background()

		bigNum1 := big.NewInt(1234567890)
		bigNum2 := big.NewInt(9876543210)

		err = composed.Set(ctx, "a", bigNum1)
		require.NoError(t, err)
		err = composed.Set(ctx, "b", bigNum2)
		require.NoError(t, err)

		val1, ok, err := composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, bigNum1, val1)

		val2, ok, err := composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, bigNum2, val2)

		err = mem.Delete(ctx, "a")
		require.NoError(t, err)
		err = red.Delete(ctx, "c")
		require.NoError(t, err)

		val1, ok, err = composed.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, bigNum1, val1)

		bigNum2, ok, err = composed.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, bigNum2, val2)
	})
}

type apiResponse struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    []byte            `json:"body"`
}

func TestTelemetryWrappedMetricsStore(t *testing.T) {
	ctx := context.Background()

	t.Run("memcache backend", func(t *testing.T) {
		backend, err := memcache.NewBackend(10)
		require.NoError(t, err)

		store := NewMetricsStore(cachestore.OpenStore[apiResponse](backend))

		err = store.Set(ctx, "a", apiResponse{Status: 1})
		require.NoError(t, err)

		err = store.Set(ctx, "b", apiResponse{Status: 2})
		require.NoError(t, err)

		ms, ok := store.(*metricsStore[apiResponse])
		require.True(t, ok)

		require.Equal(t, uint64(2), ms.telemetryCounter)

		v1, ok, err := store.Get(ctx, "a")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, apiResponse{Status: 1}, v1)

		v2, ok, err := store.Get(ctx, "b")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, apiResponse{Status: 2}, v2)

		require.Equal(t, uint64(4), ms.telemetryCounter)
	})

	t.Run("redis backend", func(t *testing.T) {
		backend, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379, DBIndex: 9})
		require.NoError(t, err)

		store := NewMetricsStore(cachestore.OpenStore[apiResponse](backend))

		err = store.Set(ctx, "c", apiResponse{Status: 3})
		require.NoError(t, err)

		err = store.Set(ctx, "d", apiResponse{Status: 4})
		require.NoError(t, err)

		ms, ok := store.(*metricsStore[apiResponse])
		require.True(t, ok)

		require.Equal(t, uint64(2), ms.telemetryCounter)

		v1, ok, err := store.Get(ctx, "c")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, apiResponse{Status: 3}, v1)

		v2, ok, err := store.Get(ctx, "d")
		require.NoError(t, err)
		require.True(t, ok)
		require.Equal(t, apiResponse{Status: 4}, v2)

		require.Equal(t, uint64(4), ms.telemetryCounter)
	})
}

package main

import (
	"context"
	"fmt"
	"log"

	memcache "github.com/goware/cachestore-mem"
	rediscache "github.com/goware/cachestore-redis"
	cachestore "github.com/goware/cachestore2"
)

func main() {
	localBackend, err := memcache.NewBackend(10)
	if err != nil {
		log.Fatal(err)
	}

	remoteBackend, err := rediscache.NewBackend(&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("!!! MEMCACHE BACKEND:")
	l1 := cachestore.OpenStore[string](localBackend)
	l1.Set(context.Background(), "key", "value!")
	value, _, err := l1.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)

	l2 := cachestore.OpenStore[int](localBackend)
	l2.Set(context.Background(), "key", 123)
	value2, _, err := l2.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value2)

	//--
	fmt.Println("")
	fmt.Println("")
	fmt.Println("!!! REDIS BACKEND:")

	r1 := cachestore.OpenStore[string](remoteBackend)
	err = r1.Set(context.Background(), "key", "zvalue!")
	if err != nil {
		panic(err)
	}
	value3, _, err := r1.Get(context.Background(), "key")
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	fmt.Println(value3)

	r2 := cachestore.OpenStore[int](remoteBackend)
	r2.Set(context.Background(), "key2", 123)
	value4, _, err := r2.Get(context.Background(), "key2")
	if err != nil {
		panic(err)
		log.Fatal(err)
	}
	fmt.Println(value4)

	//--

	// lets do compose next..
}

func main1() {
	fmt.Println("Hello, World!")

	var cache cachestore.Store[string]
	var err error
	cache, err = memcache.NewCacheWithSize[string](10)
	if err != nil {
		log.Fatal(err)
	}
	_ = cache

	err = cache.Set(context.Background(), "key", "value!")
	if err != nil {
		log.Fatal(err)
	}

	value, _, err := cache.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)

	var cache2 cachestore.Store[string]
	cache2, err = rediscache.NewCache[string](&rediscache.Config{Enabled: true, Host: "localhost", Port: 6379})
	if err != nil {
		log.Fatal(err)
	}
	_ = cache2

	err = cache2.Set(context.Background(), "key", "value2!")
	if err != nil {
		log.Fatal(err)
	}

	value, _, err = cache2.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(value)

}

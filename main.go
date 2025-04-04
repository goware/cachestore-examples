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

	fmt.Println("=> MEMCACHE BACKEND:")
	l1 := cachestore.OpenStore[string](localBackend)
	l1.Set(context.Background(), "key", "value!")
	value, _, err := l1.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get key:", value)

	l2 := cachestore.OpenStore[int](localBackend)
	l2.Set(context.Background(), "key", 123)
	value2, _, err := l2.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get key:", value2)

	//--
	fmt.Println("")
	fmt.Println("")
	fmt.Println("=> REDIS BACKEND:")

	r1 := cachestore.OpenStore[string](remoteBackend)
	err = r1.Set(context.Background(), "key", "zvalue!")
	if err != nil {
		panic(err)
	}
	value3, _, err := r1.Get(context.Background(), "key")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get key:", value3)

	r2 := cachestore.OpenStore[int](remoteBackend)
	r2.Set(context.Background(), "key2", 123)
	value4, _, err := r2.Get(context.Background(), "key2")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("get key:", value4)

	// Please see cachestore_e2e_test.go for more examples.
}

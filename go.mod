module github.com/goware/cachestore-examples

// replace github.com/goware/cachestore2 => ../cachestore2

// replace github.com/goware/cachestore-mem => ../cachestore-mem

// replace github.com/goware/cachestore-redis => ../cachestore-redis

go 1.23.0

require (
	github.com/goware/cachestore-mem v0.2.1
	github.com/goware/cachestore-redis v0.2.0
	github.com/goware/cachestore2 v0.12.0
	github.com/redis/go-redis/v9 v9.7.3
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/elastic/go-freelru v0.16.0 // indirect
	github.com/goware/singleflight v0.3.0 // indirect
	github.com/klauspost/cpuid/v2 v2.2.10 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/zeebo/xxh3 v1.0.2 // indirect
	golang.org/x/sys v0.30.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

module github.com/goware/cachestore-examples

replace github.com/goware/cachestore => ../cachestore

replace github.com/goware/cachestore-mem => ../cachestore-mem

replace github.com/goware/cachestore-redis => ../cachestore-redis

go 1.24.1

require (
	github.com/goware/cachestore v0.11.0
	github.com/goware/cachestore-mem v0.0.0-00010101000000-000000000000
	github.com/goware/cachestore-redis v0.0.0-00010101000000-000000000000
	github.com/redis/go-redis/v9 v9.7.3
	github.com/stretchr/testify v1.10.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/goware/singleflight v0.3.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

package client

// Cache implements a cache for the securitron service.
type Cache interface {
	Get(path, key string) KV
	Put(path string, val KV) bool
	Del(path, key string) (KV, bool)
}

// NoCache disables caching
type NoCache struct{}

func (NoCache) Del(path, key string) (KV, bool) { return nil, false }
func (NoCache) Get(path, key string) KV         { return nil }
func (NoCache) Put(path string, kv KV) bool     { return false }

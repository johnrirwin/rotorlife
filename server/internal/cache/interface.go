package cache

import "time"

// Cache defines the interface for cache backends
type Cache interface {
	Get(key string) (interface{}, bool)
	Set(key string, value interface{})
	SetWithTTL(key string, value interface{}, ttl time.Duration)
	Delete(key string)
	Clear()
}

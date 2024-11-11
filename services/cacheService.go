package services

import (
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

// CacheService provides in-memory caching with expiration
type CacheService struct {
	items map[string]CacheItem
	mu    sync.RWMutex
}

// NewCacheService initializes CacheService
func NewCacheService() *CacheService {
	return &CacheService{
		items: make(map[string]CacheItem),
	}
}

// Set stores data in the cache with a 10-minute expiration time
func (c *CacheService) Set(key string, data interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = CacheItem{
		Data:      data,
		ExpiresAt: time.Now().Add(10 * time.Minute),
	}
}

// Get retrieves data from the cache if not expired
func (c *CacheService) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Data, true
}

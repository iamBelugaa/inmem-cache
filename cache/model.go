package cache

import (
	"sync"
	"time"
)

type Cache[K comparable, V any] struct {
	cleanupChan   chan bool
	cleanupTicker *time.Ticker
	mutex         sync.RWMutex
	cacheDuration time.Duration
	actions       *Actions[K, V]
	cache         map[K]*CacheItem[V]
}

type CacheItem[V any] struct {
	value        V
	lastAccessed time.Time
}

type Actions[K comparable, V any] struct {
	OnCleanup func(key K, val V)
}

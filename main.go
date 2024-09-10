package main

import (
	"fmt"
	"sync"
	"time"

	"github.com/apache/arrow/go/v14/arrow"
	"github.com/apache/arrow/go/v14/arrow/array"
	"github.com/apache/arrow/go/v14/arrow/memory"
)

type Mode int

const (
	CPU Mode = iota
	GPU
)

type CacheItem struct {
	Value     []byte
	ExpiresAt time.Time
}

type Cache struct {
	items map[string]CacheItem
	mu    sync.RWMutex
	mode  Mode
	pool  *memory.GoAllocator
}

func NewCache(mode Mode) *Cache {
	return &Cache{
		items: make(map[string]CacheItem),
		mode:  mode,
		pool:  memory.NewGoAllocator(),
	}
}

func (c *Cache) SetWithArrow(key string, value []byte, ttl time.Duration) {
	builder := array.NewBinaryBuilder(c.pool, arrow.BinaryTypes.Binary)
	defer builder.Release()

	builder.Append(value)
	arr := builder.NewArray()
	defer arr.Release()

	c.mu.Lock()
	c.items[key] = CacheItem{
		Value:     arr.(*array.Binary).Value(0),
		ExpiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}

func (c *Cache) SetWithGPU(key string, value []byte, ttl time.Duration) {
	//fmt.Printf("Processing data on the GPU for key: %s\n", key)
	c.mu.Lock()
	c.items[key] = CacheItem{
		Value:     value,
		ExpiresAt: time.Now().Add(ttl),
	}
	c.mu.Unlock()
}

func (c *Cache) Set(key string, value []byte, ttl time.Duration) {
	if c.mode == CPU {
		c.SetWithArrow(key, value, ttl)
	} else if c.mode == GPU {
		c.SetWithGPU(key, value, ttl)
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, found := c.items[key]
	if !found || time.Now().After(item.ExpiresAt) {
		return nil, false
	}
	return item.Value, true
}

func (c *Cache) CleanExpiredItems() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for k, v := range c.items {
		if time.Now().After(v.ExpiresAt) {
			delete(c.items, k)
		}
	}
}

func BenchmarkSet(cache *Cache, wg *sync.WaitGroup, numKeys int) {
	startTime := time.Now()

	for i := 0; i < numKeys; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Set(fmt.Sprintf("key%d", i), []byte("some data"), 10*time.Minute)
		}(i)
	}

	wg.Wait()
	endTime := time.Now()
	fmt.Printf("Time taken to Set %d keys in %s mode: %v\n", numKeys, cache.mode, endTime.Sub(startTime))
}

func BenchmarkGet(cache *Cache, wg *sync.WaitGroup, numKeys int) {
	startTime := time.Now()

	for i := 0; i < numKeys; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cache.Get(fmt.Sprintf("key%d", i))
		}(i)
	}

	wg.Wait()
	endTime := time.Now()
	fmt.Printf("Time taken to Get %d keys in %s mode: %v\n", numKeys, cache.mode, endTime.Sub(startTime))
}

func main() {
	numKeys := 1000000
	var wg sync.WaitGroup

	// Initialize the cache in CPU mode
	cacheCPU := NewCache(CPU)
	fmt.Println("Benchmarking CPU Mode")
	BenchmarkSet(cacheCPU, &wg, numKeys)
	BenchmarkGet(cacheCPU, &wg, numKeys)

	// Initialize the cache in GPU mode
	cacheGPU := NewCache(GPU)
	fmt.Println("Benchmarking GPU Mode")
	BenchmarkSet(cacheGPU, &wg, numKeys)
	BenchmarkGet(cacheGPU, &wg, numKeys)
}

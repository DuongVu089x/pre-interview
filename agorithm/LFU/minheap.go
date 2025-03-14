package lfu

import (
	"container/heap"
	"fmt"
)

// CacheItem represents an item in the cache
type CacheItem struct {
	key       string
	value     interface{}
	frequency int
	timestamp int64 // Để xử lý trường hợp frequency bằng nhau
}

// MinHeap implementation
type MinHeap []*CacheItem

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	// So sánh frequency trước, nếu bằng nhau thì so sánh timestamp
	if h[i].frequency == h[j].frequency {
		return h[i].timestamp < h[j].timestamp
	}
	return h[i].frequency < h[j].frequency
}
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(*CacheItem)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[0 : n-1]
	return item
}

// LFUCache structure
type LFUCache struct {
	capacity int
	items    map[string]*CacheItem
	heap     *MinHeap
	time     int64 // Để track timestamp
}

// NewLFUCache creates a new LFU Cache
func NewLFUCache(capacity int) *LFUCache {
	h := &MinHeap{}
	heap.Init(h)
	return &LFUCache{
		capacity: capacity,
		items:    make(map[string]*CacheItem),
		heap:     h,
		time:     0,
	}
}

// Get retrieves an item from the cache
func (c *LFUCache) Get(key string) (interface{}, bool) {
	if item, exists := c.items[key]; exists {
		// Tăng frequency và cập nhật timestamp
		item.frequency++
		c.time++
		item.timestamp = c.time
		// Rebuild heap để maintain order
		heap.Init(c.heap)
		return item.value, true
	}
	return nil, false
}

// Put adds an item to the cache
func (c *LFUCache) Put(key string, value interface{}) {
	c.time++

	// Nếu key đã tồn tại, cập nhật value và frequency
	if item, exists := c.items[key]; exists {
		item.value = value
		item.frequency++
		item.timestamp = c.time
		heap.Init(c.heap)
		return
	}

	// Tạo item mới
	newItem := &CacheItem{
		key:       key,
		value:     value,
		frequency: 1,
		timestamp: c.time,
	}

	// Nếu cache đầy, remove item có frequency thấp nhất
	if len(c.items) >= c.capacity {
		leastFreq := heap.Pop(c.heap).(*CacheItem)
		delete(c.items, leastFreq.key)
	}

	// Thêm item mới
	c.items[key] = newItem
	heap.Push(c.heap, newItem)
}

// PrintCache prints the current state of the cache
func (c *LFUCache) PrintCache() {
	fmt.Println("\nCache state:")
	for key, item := range c.items {
		fmt.Printf("Key: %s, Value: %v, Frequency: %d\n",
			key, item.value, item.frequency)
	}
}

func TestLFU() {
	// Tạo cache với capacity = 3
	cache := NewLFUCache(3)

	// Test cases
	cache.Put("A", 1)
	cache.Put("B", 2)
	cache.Put("C", 3)
	cache.PrintCache()

	// Access "A" twice
	cache.Get("A")
	cache.Get("A")
	cache.PrintCache()

	// Add new item when cache is full
	cache.Put("D", 4) // Should remove least frequently used item
	cache.PrintCache()

	// Access "B"
	cache.Get("B")
	cache.PrintCache()

	// Update existing item
	cache.Put("B", 20)
	cache.PrintCache()
}

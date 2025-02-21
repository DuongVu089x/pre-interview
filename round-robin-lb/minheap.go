package roundrobinlb

import (
	"container/heap"
	"fmt"
)

// Server represents a server with its index and availability time
type Server struct {
	index     int
	available int // When the server becomes available
}

// MinHeap implements a priority queue for servers
type MinHeap []Server

func (h MinHeap) Len() int { return len(h) }
func (h MinHeap) Less(i, j int) bool {
	// Sort by available time, then by index
	if h[i].available == h[j].available {
		return h[i].index < h[j].index
	}
	return h[i].available < h[j].available
}
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Server)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getServerIndex(n int, arrival []int, burstTime []int) []int {
	result := make([]int, len(arrival))
	serverHeap := &MinHeap{}
	heap.Init(serverHeap)

	// Initialize servers
	for i := 1; i <= n; i++ {
		heap.Push(serverHeap, Server{index: i, available: 0})
	}

	for i := 0; i < len(arrival); i++ {
		requestTime := arrival[i]
		burst := burstTime[i]

		// Find the first available server
		for serverHeap.Len() > 0 && (*serverHeap)[0].available <= requestTime {
			server := heap.Pop(serverHeap).(Server)
			result[i] = server.index
			server.available = requestTime + burst
			heap.Push(serverHeap, server)
			break
		}

		// If no server is available at this request time
		if result[i] == 0 {
			result[i] = -1
		}
	}

	return result
}

func TestMinHeap() {
	n := 3
	arrival := []int{2, 4, 1, 8, 9}
	burstTime := []int{7, 9, 2, 4, 5}
	result := getServerIndex(n, arrival, burstTime)
	fmt.Println(result) // Output: [2, 1, 1, 3, 2]
}

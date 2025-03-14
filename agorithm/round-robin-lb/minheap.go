package roundrobinlb

import (
	"fmt"
)

type Server struct {
	index int32
	avaAt int32
}

func getServerIndex(n int32, arrival []int32, burstTime []int32) []int32 {
	m := int32(len(arrival))
	result := make([]int32, m)
	servers := make([]Server, n)

	// Initialize servers
	for i := int32(0); i < n; i++ {
		servers[i] = Server{index: i + 1, avaAt: 0}
	}

	// Process each request
	for i := int32(0); i < m; i++ {
		requestTime := arrival[i]
		burst := burstTime[i]
		assigned := false

		// Find the first available server
		for j := int32(0); j < n; j++ {
			if servers[j].avaAt <= requestTime {
				result[i] = servers[j].index
				servers[j].avaAt = requestTime + burst
				assigned = true
				break
			}
		}

		// If no server is available, mark as -1
		if !assigned {
			result[i] = -1
		}
	}

	return result
}

func TestMinHeap() {
	n := int32(2)
	arrival := []int32{}
	burstTime := []int32{1, 2, 1, 2}
	result := getServerIndex(n, arrival, burstTime)
	fmt.Println(result) // Should print [1 2 -1 1]
}

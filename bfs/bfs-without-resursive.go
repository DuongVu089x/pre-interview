package bfs

import (
	"fmt"
)

type Graph struct {
	adj map[int][]int
}

func (g *Graph) addEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	g.adj[v] = append(g.adj[v], u) // Giả sử đồ thị vô hướng
}

func (g *Graph) BFSWithoutRecursive(start int) {
	visited := make(map[int]bool)
	queue := []int{start}

	for len(queue) > 0 {
		// Lấy phần tử đầu tiên trong queue
		node := queue[0]
		queue = queue[1:] // pop first element

		// Nếu chưa thăm thì xử lý và đánh dấu đã thăm
		if !visited[node] {
			fmt.Print(node, " ")
			visited[node] = true
		}

		// Đưa các đỉnh kề chưa thăm vào queue
		for _, neighbor := range g.adj[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}
}

func TestWithoutRecursive() {
	g := Graph{adj: make(map[int][]int)}
	g.addEdge(0, 1)
	g.addEdge(0, 2)
	g.addEdge(1, 3)
	g.addEdge(1, 4)
	g.addEdge(2, 5)
	g.addEdge(2, 6)

	fmt.Println("BFS Traversal:")
	g.BFSWithoutRecursive(0)
}

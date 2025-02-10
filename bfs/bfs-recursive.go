package bfs

import "fmt"

func (g *Graph) BFS(visited map[int]bool, queue []int) {

	if len(queue) == 0 {
		return
	}

	// Lấy phần tử đầu tiên trong queue
	node := queue[0]
	queue = queue[1:]

	// Nếu chưa thăm thì xử lý và đánh dấu đã thăm
	if !visited[node] {
		fmt.Print(node, " ")
		visited[node] = true

		// Thêm các đỉnh kề vào hàng đợi
		for _, neighbor := range g.adj[node] {
			if !visited[neighbor] {
				queue = append(queue, neighbor)
			}
		}
	}

	// Đệ quy với danh sách queue mới
	g.BFS(visited, queue)
}

func Test() {
	g := Graph{adj: make(map[int][]int)}
	g.addEdge(0, 1)
	g.addEdge(0, 2)
	g.addEdge(1, 3)
	g.addEdge(1, 4)
	g.addEdge(2, 5)
	g.addEdge(2, 6)

	visited := make(map[int]bool)
	queue := []int{0}

	fmt.Println("BFS Traversal:")
	g.BFS(visited, queue)
}

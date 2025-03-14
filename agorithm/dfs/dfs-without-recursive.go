package dfs

import "fmt"

// // Graph representation using an adjacency list
// type Graph struct {
// 	adj map[int][]int
// }

// // Add an edge to the graph
// func (g *Graph) addEdge(u, v int) {
// 	g.adj[u] = append(g.adj[u], v)
// 	g.adj[v] = append(g.adj[v], u) // Assuming an undirected graph
// }

// DFS without recursion (Using a stack)
func (g *Graph) DFSWithoutRecursive(start int) {
	visited := make(map[int]bool)
	stack := []int{start}

	for len(stack) > 0 {
		// Pop the last element
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		// If not visited, process the node
		if !visited[node] {
			fmt.Print(node, " ")
			visited[node] = true
		}

		// Push all unvisited neighbors onto the stack
		for _, neighbor := range g.adj[node] {
			if !visited[neighbor] {
				stack = append(stack, neighbor)
			}
		}
	}
}

// func main() {
// 	g := Graph{adj: make(map[int][]int)}
// 	g.addEdge(0, 1)
// 	g.addEdge(0, 2)
// 	g.addEdge(1, 3)
// 	g.addEdge(1, 4)
// 	g.addEdge(2, 5)
// 	g.addEdge(2, 6)

// 	fmt.Println("DFS Traversal (Non-Recursive):")
// 	g.DFS(0)
// }

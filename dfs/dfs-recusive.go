package main

import "fmt"

// Graph representation using an adjacency list
type Graph struct {
	adj map[int][]int
}

// Add an edge to the graph
func (g *Graph) addEdge(u, v int) {
	g.adj[u] = append(g.adj[u], v)
	g.adj[v] = append(g.adj[v], u) // Assuming an undirected graph
}

// DFS recursive helper function
func (g *Graph) dfsHelper(node int, visited map[int]bool) {
	// Process current node
	fmt.Print(node, " ")
	visited[node] = true

	// Recursively visit all unvisited neighbors
	for _, neighbor := range g.adj[node] {
		if !visited[neighbor] {
			g.dfsHelper(neighbor, visited)
		}
	}
}

// DFS using recursion
func (g *Graph) DFS(start int) {
	visited := make(map[int]bool)
	g.dfsHelper(start, visited)
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

package bfs

import "fmt"

type Point struct {
	x, y int
}

var directions = []Point{
	{0, 1},  // Right
	{1, 0},  // Down
	{0, -1}, // Left
	{-1, 0}, // Up
}

func solveMaze(maze [][]int, start, end Point) bool {

	row, col := len(maze), len(maze[0])
	queue := []Point{start}

	visited := make(map[Point]bool)
	parent := make(map[Point]Point)

	for len(queue) > 0 {
		curr := queue[0]  // get first element
		queue = queue[1:] // pop first element

		if curr == end {
			fmt.Println("Path found!")
			printPath(parent, start, end)
			return true
		}

		visited[curr] = true

		for _, d := range directions {
			next := Point{curr.x + d.x, curr.y + d.y}

			if next.x >= 0 && next.x < row && next.y >= 0 && next.y < col && // check bounds
				maze[next.x][next.y] == 0 && // check if the cell is open
				!visited[next] { // check if the cell is not visited

				if next == end {
					fmt.Println("Path found!")
				}

				queue = append(queue, next)
				parent[next] = curr // save path
			}
		}
	}

	fmt.Println("No path found.")
	return false
}

// Hàm in đường đi từ start -> end
func printPath(parent map[Point]Point, start, end Point) {
	path := []Point{}
	for at := end; at != start; at = parent[at] {
		path = append([]Point{at}, path...)
	}
	path = append([]Point{start}, path...)

	fmt.Println("Shortest Path:", path)
}

func TestMaze() {
	maze := [][]int{
		{0, 1, 0, 0, 0},
		{0, 1, 0, 1, 0},
		{0, 0, 0, 1, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 0, 1, 0},
	}

	start := Point{0, 0}
	end := Point{4, 4}

	solveMaze(maze, start, end)
}

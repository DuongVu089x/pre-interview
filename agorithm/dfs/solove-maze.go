package dfs

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

	stack := []Point{start}
	visited := make(map[Point]bool)

	for len(stack) > 0 {

		// current point is the last point in stack
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1] // pop last point

		// if reached the end, return true
		if curr == end {
			fmt.Println("Path found!")
			return true
		}

		// mark visited
		visited[curr] = true

		// try all with 4 directions
		for _, d := range directions {
			next := Point{curr.x + d.x, curr.y + d.y}

			// check bounds and if the cell is open and not visited
			if next.x >= 0 && next.x < row && next.y >= 0 && next.y < col && // check bounds
				maze[next.x][next.y] == 0 && // check if the cell is open
				!visited[next] { // check if the cell is not visited
				stack = append(stack, next)
			}
		}
	}

	fmt.Println("No path found.")
	return false
}

// func main() {
// 	maze := [][]int{
// 		{0, 0, 1, 0, 0},
// 		{1, 0, 1, 0, 1},
// 		{1, 0, 0, 0, 1},
// 		{1, 1, 1, 0, 0},
// 	}

// 	start := Point{0, 0}
// 	end := Point{3, 4}

// 	solveMaze(maze, start, end)
// }

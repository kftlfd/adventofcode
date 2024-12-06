package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	inpPath := os.Args[1]
	inpFile, err := os.Open(inpPath)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(inpFile)
	fileScanner.Split(bufio.ScanLines)

	grid := [][]string{}

	for fileScanner.Scan() {
		row := fileScanner.Text()
		grid = append(grid, strings.Split(row, ""))
	}

	m := len(grid)
	n := len(grid[0])
	guard_start := [2]int{0, 0}
	for r, row := range grid {
		stop := false
		for c, cell := range row {
			if cell == "^" {
				guard_start = [2]int{r, c}
				stop = true
				break
			}
		}
		if stop {
			break
		}
	}

	get_grid_copy := func() [][]string {
		g := make([][]string, len(grid))
		for i, row := range grid {
			r := make([]string, len(row))
			copy(r, row)
			g[i] = r
		}
		return g
	}

	get_guard := func() Guard {
		return Guard{row: guard_start[0], col: guard_start[1], dir: 0}
	}

	//
	// part 1: count X (guard path)
	//

	g := get_grid_copy()
	guard := get_guard()

	for true {
		g[guard.row][guard.col] = "X"

		nr := guard.row + dirs[guard.dir][0]
		nc := guard.col + dirs[guard.dir][1]

		if nr < 0 || nr >= m || nc < 0 || nc >= n {
			break
		}

		if g[nr][nc] == "#" {
			guard.dir = (guard.dir + 1) % len(dirs)
		} else {
			guard.row = nr
			guard.col = nc
		}
	}

	guard_path := [][2]int{}

	x_count := 0
	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if g[r][c] == "X" {
				x_count += 1
				guard_path = append(guard_path, [2]int{r, c})
			}
		}
	}

	fmt.Println("X count:", x_count)

	//
	// part 2: count loops after adding one obstacle
	//

	loops_count := 0

	for _, path := range guard_path {
		if path[0] == guard_start[0] && path[1] == guard_start[1] {
			continue
		}

		g = get_grid_copy()
		g[path[0]][path[1]] = "O"
		guard = get_guard()
		if checkLoops(g, guard) {
			loops_count += 1
		}
	}

	fmt.Println("loops after adding obstacle:", loops_count)
}

var dirs = [4][2]int{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

type Guard struct {
	row, col, dir int
}

func printGrid(g [][]string) {
	for _, row := range g {
		fmt.Println(row)
	}
}

func str(val int) string {
	return strconv.Itoa(val)
}

func checkLoops(g [][]string, guard Guard) bool {
	// mark obstacles as visited from certain direction
	// if approaching an obstacle twice from the same direction = loop

	m := len(g)
	n := len(g[0])

	for true {
		nr := guard.row + dirs[guard.dir][0]
		nc := guard.col + dirs[guard.dir][1]

		if nr < 0 || nr >= m || nc < 0 || nc >= n {
			return false
		}
		if isVisited(g[nr][nc], guard.dir) {
			return true
		}

		if isObstacle(g[nr][nc]) {
			g[nr][nc] = addVisited(g[nr][nc], guard.dir)
			guard.dir = (guard.dir + 1) % len(dirs)
		} else {
			guard.row = nr
			guard.col = nc
		}
	}

	return false
}

func isObstacle(val string) bool {
	f := strings.Split(val, "")
	return f[0] == "#" || f[0] == "O"
}

func isVisited(val string, dir int) bool {
	vals := strings.Split(val, "")
	if vals[0] != "#" && vals[0] != "O" {
		return false
	}
	for _, v := range vals {
		if v == str(dir) {
			return true
		}
	}
	return false
}

func addVisited(val string, dir int) string {
	v := strings.Split(val, "")
	v = append(v, str(dir))
	return strings.Join(v, "")
}

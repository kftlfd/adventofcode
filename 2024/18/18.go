package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	// GRID_SIZE   = 7
	// BYTES_COUNT = 12

	GRID_SIZE   = 71
	BYTES_COUNT = 1024

	R = 0
	C = 1
)

var DIRS = [4][2]int{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

func getInput(file *os.File) [][2]int {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	bytes := [][2]int{}

	for scanner.Scan() {
		line := scanner.Text()

		vals := strings.Split(line, ",")
		if len(vals) != 2 {
			panic(fmt.Errorf("bad input: %v", line))
		}

		x, err := strconv.Atoi(vals[0])
		if err != nil {
			panic(err)
		}

		y, err := strconv.Atoi(vals[1])
		if err != nil {
			panic(err)
		}

		bytes = append(bytes, [2]int{x, y})
	}

	return bytes
}

func getGrid(bytes [][2]int) [][]string {
	grid := [][]string{}
	for r := 0; r < GRID_SIZE; r++ {
		grid = append(grid, make([]string, GRID_SIZE))
	}

	for r := 0; r < GRID_SIZE; r++ {
		for c := 0; c < GRID_SIZE; c++ {
			grid[r][c] = "."
		}
	}

	for _, b := range bytes {
		x := b[0]
		y := b[1]
		grid[y][x] = "#"
	}

	return grid
}

//
//
//

func bfs(grid [][]string) int {
	end := GRID_SIZE - 1

	visited := [][]bool{}
	for r := 0; r < GRID_SIZE; r++ {
		visited = append(visited, make([]bool, GRID_SIZE))
	}
	visited[0][0] = true

	q := [][2]int{{0, 0}}
	i := 0

	for len(q) > 0 {
		nxt_q := [][2]int{}

		for j := 0; j < len(q); j++ {
			cur := q[j]

			if cur[R] == end && cur[C] == end {
				return i
			}

			for _, dir := range DIRS {
				nr := cur[R] + dir[R]
				nc := cur[C] + dir[C]

				if nr < 0 || nr >= GRID_SIZE || nc < 0 || nc >= GRID_SIZE || grid[nr][nc] == "#" || visited[nr][nc] {
					continue
				}

				visited[nr][nc] = true

				nxt_q = append(nxt_q, [2]int{nr, nc})
			}
		}

		q = nxt_q
		i += 1
	}

	return -1
}

//
// part 2
//

func binarySearchFirstUnreachable(bytes [][2]int) int {
	lo := 0
	hi := len(bytes)

	for lo < hi {
		mid := lo + (hi-lo)/2

		grid := getGrid(bytes[:mid+1])
		is_reachable := bfs(grid) > -1

		if is_reachable {
			lo = mid + 1
		} else {
			hi = mid
		}
	}

	return lo
}

//
//
//

func main() {
	input_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	bytes := getInput(input_file)

	// part 1
	grid := getGrid(bytes[:BYTES_COUNT])
	min_path := bfs(grid)
	fmt.Println("min path:", min_path)

	// part 2
	byte_idx := binarySearchFirstUnreachable(bytes)
	byte_val := bytes[byte_idx]
	fmt.Printf("byte that makes grid exit unreachable: %v,%v\n", byte_val[0], byte_val[1])
}

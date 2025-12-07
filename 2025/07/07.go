package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseInput() []string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		inp = append(inp, line)
	}

	return inp
}

func solvePart1(input []string) {
	grid := [][]string{}
	for _, inpRow := range input {
		grid = append(grid, strings.Split(inpRow, ""))
	}
	for i, v := range grid[0] {
		if v == "S" {
			grid[0][i] = "1"
		}
	}

	m := len(grid)
	n := len(grid[0])
	totalSplits := 0

	for r := 1; r < m; r++ {
		for c := range n {
			prev := grid[r-1][c]
			cur := grid[r][c]
			if prev == "." || prev == "^" {
				continue
			}
			if cur == "^" {
				totalSplits += 1
				grid[r][c-1] = prev
				grid[r][c+1] = prev
			} else {
				grid[r][c] = prev
			}
		}
	}

	fmt.Println("Part 1:", totalSplits)
}

func solvePart2(input []string) {
	// convert to grid of ints:
	// "." = 0
	// "^" = -1
	// positive ints - number of beams at this position
	grid := [][]int{}
	for _, inpLine := range input {
		inpRow := strings.Split(inpLine, "")
		curRow := []int{}
		for _, val := range inpRow {
			switch val {
			case "S":
				curRow = append(curRow, 1)
			case ".":
				curRow = append(curRow, 0)
			case "^":
				curRow = append(curRow, -1)
			default:
				panic(fmt.Errorf("invalid input: %s", val))
			}
		}
		grid = append(grid, curRow)
	}

	m := len(input)
	n := len(input[0])

	// fill beams paths
	for r := 1; r < m; r++ {
		for c := range n {
			prev := grid[r-1][c]
			cur := grid[r][c]

			if prev < 1 {
				continue
			}

			if cur != -1 {
				grid[r][c] += prev
				continue
			}

			grid[r][c-1] += prev
			grid[r][c+1] += prev
		}
	}

	// add total beams paths in last row
	totalPaths := 0
	for _, paths := range grid[len(grid)-1] {
		if paths > 0 {
			totalPaths += paths
		}
	}

	fmt.Println("Part 2:", totalPaths)
}

func main() {
	input := parseInput()
	solvePart1(input)
	solvePart2(input)
}

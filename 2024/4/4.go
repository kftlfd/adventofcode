package main

import (
	"bufio"
	"fmt"
	"os"
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
		row := strings.Split(fileScanner.Text(), "")
		grid = append(grid, row)
	}

	// for _, row := range grid {
	// 	fmt.Printf("%+v\n", row)
	// }

	m := len(grid)
	n := len(grid[0])

	// fmt.Println(m, n)

	ansXmas := 0
	ansMas := 0

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			ansXmas += checkXMAS(grid, r, c)
			ansMas += checkMAS(grid, r, c)
		}
	}

	fmt.Println("total XMAS:", ansXmas)
	fmt.Println("total MAS:", ansMas)
}

var DIRS = [8][2]int{
	{-1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
	{1, 0},
	{1, -1},
	{0, -1},
	{-1, -1},
}

var TARGET = [4]string{"X", "M", "A", "S"}

func checkXMAS(grid [][]string, r int, c int) int {
	m := len(grid)
	n := len(grid[0])

	if grid[r][c] != TARGET[0] {
		return 0
	}

	count := 0

	for _, dir := range DIRS {
		ok := true
		for i := 1; i < 4; i++ {
			cur_r := r + dir[0]*i
			cur_c := c + dir[1]*i
			if cur_r < 0 || cur_r >= m ||
				cur_c < 0 || cur_c >= n ||
				grid[cur_r][cur_c] != TARGET[i] {
				ok = false
				break
			}
		}
		if ok {
			count += 1
		}
	}

	return count
}

func checkMAS(grid [][]string, r int, c int) int {
	m := len(grid)
	n := len(grid[0])
	if grid[r][c] != "A" ||
		r < 1 || r >= m-1 ||
		c < 1 || c >= n-1 {
		return 0
	}

	d1 := strings.Join([]string{grid[r-1][c-1], grid[r+1][c+1]}, "")
	d2 := strings.Join([]string{grid[r-1][c+1], grid[r+1][c-1]}, "")

	if (d1 == "MS" || d1 == "SM") && (d2 == "MS" || d2 == "SM") {
		return 1
	}

	return 0
}

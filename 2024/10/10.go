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

	input := [][]int{}

	for fileScanner.Scan() {
		row := strings.Split(fileScanner.Text(), "")
		row_ints := make([]int, len(row))
		for i, val := range row {
			val_int, err := strconv.Atoi(val)
			if err != nil {
				panic(fmt.Errorf("not a number: %v", val))
			}
			row_ints[i] = val_int
		}
		input = append(input, row_ints)
	}

	// for _, row := range input {
	// 	fmt.Println(row)
	// }

	m := len(input)
	n := len(input[0])

	trailheads_score := 0

	score2 := 0
	rating := 0

	scoreDfs := 0
	ratingDfs := 0

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			trailheads_score += getTrailheadScore(input, r, c)

			sc2, rat2 := getTrailheadScoreAndRating(input, r, c)
			score2 += sc2
			rating += rat2

			sc3, rat3 := getScoreRatingDFS(input, r, c)
			scoreDfs += sc3
			ratingDfs += rat3
		}
	}

	fmt.Println("trailheads score:", trailheads_score)
	fmt.Println("score2:", score2, "rating:", rating)
	fmt.Println("scoreDfs:", scoreDfs, "ratingDfs:", ratingDfs)
}

var DIRS = [4][2]int{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

func getTrailheadScore(grid [][]int, r, c int) int {
	if grid[r][c] != 0 {
		return 0
	}

	m := len(grid)
	n := len(grid[0])

	visited := make([][]bool, m)
	for row := 0; row < m; row++ {
		visited[row] = make([]bool, n)
	}
	visited[r][c] = true

	score := 0

	q := [][2]int{{r, c}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		for _, dir := range DIRS {
			nr := cur[0] + dir[0]
			nc := cur[1] + dir[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue
			}
			if visited[nr][nc] {
				continue
			}
			if grid[nr][nc] != grid[cur[0]][cur[1]]+1 {
				continue
			}

			visited[nr][nc] = true

			if grid[nr][nc] == 9 {
				score += 1
			} else {
				q = append(q, [2]int{nr, nc})
			}
		}
	}

	return score
}

func getTrailheadScoreAndRating(grid [][]int, r, c int) (int, int) {
	if grid[r][c] != 0 {
		return 0, 0
	}

	m := len(grid)
	n := len(grid[0])

	ends := make([][]int, m)
	for row := 0; row < m; row++ {
		ends[row] = make([]int, n)
	}

	rating := 0

	q := [][2]int{{r, c}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		for _, dir := range DIRS {
			nr := cur[0] + dir[0]
			nc := cur[1] + dir[1]
			if nr < 0 || nr >= m || nc < 0 || nc >= n {
				continue
			}
			if grid[nr][nc] != grid[cur[0]][cur[1]]+1 {
				continue
			}

			if grid[nr][nc] == 9 {
				ends[nr][nc] = 1
				rating += 1
			} else {
				q = append(q, [2]int{nr, nc})
			}
		}
	}

	score := 0
	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			score += ends[row][col]
		}
	}

	return score, rating
}

func dfs(grid, ends [][]int, r, c, m, n int) {
	for _, dir := range DIRS {
		nr := r + dir[0]
		nc := c + dir[1]
		if nr < 0 || nr >= m || nc < 0 || nc >= n || grid[nr][nc] != grid[r][c]+1 {
			continue
		}

		if grid[nr][nc] == 9 {
			ends[nr][nc] += 1
		} else {
			dfs(grid, ends, nr, nc, m, n)
		}
	}
}

func getScoreRatingDFS(grid [][]int, r, c int) (score, rating int) {
	if grid[r][c] != 0 {
		return 0, 0
	}

	m := len(grid)
	n := len(grid[0])

	ends := make([][]int, m)
	for row := 0; row < m; row++ {
		ends[row] = make([]int, n)
	}

	dfs(grid, ends, r, c, m, n)

	score = 0
	rating = 0

	for row := 0; row < m; row++ {
		for col := 0; col < n; col++ {
			if ends[row][col] > 0 {
				score += 1
				rating += ends[row][col]
			}
		}
	}

	return score, rating
}

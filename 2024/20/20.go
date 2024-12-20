package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	R = 0
	C = 1
)

var DIRS = [4][2]int{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

type Racetrack struct {
	grid          [][]string
	m, n          int
	start, end    [2]int
	shortest_time int
	cheat_savings map[int]int
}

func parseInput(file *os.File) Racetrack {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	grid := [][]string{}

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		grid = append(grid, line)
	}

	m := len(grid)
	n := len(grid[0])

	racetrack := Racetrack{grid: grid, m: m, n: n, cheat_savings: make(map[int]int)}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == "S" {
				racetrack.start = [2]int{r, c}
			}
			if grid[r][c] == "E" {
				racetrack.end = [2]int{r, c}
			}
		}
	}

	return racetrack
}

//
//
//

func (rt *Racetrack) getShortestTime() int {
	q := [][2]int{rt.start}

	visited := [][]bool{}
	for r := 0; r < rt.m; r++ {
		visited = append(visited, make([]bool, rt.n))
	}
	visited[rt.start[R]][rt.start[C]] = true

	time := 0

	for len(q) > 0 {
		cur := q[0]
		q = q[1:]

		if cur[R] == rt.end[R] && cur[C] == rt.end[C] {
			break
		}

		for _, dir := range DIRS {
			nr := cur[R] + dir[R]
			nc := cur[C] + dir[C]

			if nr < 0 || nr >= rt.m || nc < 0 || nc >= rt.n || rt.grid[nr][nc] == "#" || visited[nr][nc] {
				continue
			}

			visited[nr][nc] = true
			q = append(q, [2]int{nr, nc})
		}

		time += 1
	}

	rt.shortest_time = time
	return time
}

func (rt *Racetrack) cheatDFS(r, c, time, cheats_left int, cheats_enabled bool, visited *[][]bool) {
	if time >= rt.shortest_time {
		return
	}

	if r == rt.end[R] && c == rt.end[C] {
		saved := rt.shortest_time - time
		rt.cheat_savings[saved] += 1
		return
	}

	if cheats_enabled {
		cheats_left -= 1
	}

	if rt.grid[r][c] == "#" && cheats_left < 1 {
		return
	}

	for _, dir := range DIRS {
		nr := r + dir[R]
		nc := c + dir[C]

		if nr < 0 || nr >= rt.m || nc < 0 || nc >= rt.n || (*visited)[nr][nc] {
			continue
		}

		if rt.grid[nr][nc] == "#" {
			if cheats_left > 0 {
				(*visited)[nr][nc] = true
				rt.cheatDFS(nr, nc, time+1, cheats_left, true, visited)
				(*visited)[nr][nc] = false
			}
		} else {
			(*visited)[nr][nc] = true
			rt.cheatDFS(nr, nc, time+1, cheats_left, cheats_enabled, visited)
			(*visited)[nr][nc] = false
		}
	}
}

func (rt *Racetrack) getCheatsSavings(cheats_allowed int) map[int]int {
	visited := [][]bool{}
	for r := 0; r < rt.m; r++ {
		visited = append(visited, make([]bool, rt.n))
	}
	visited[rt.start[R]][rt.start[C]] = true

	rt.cheatDFS(rt.start[R], rt.start[C], 0, cheats_allowed, false, &visited)
	return rt.cheat_savings
}

//
//
//

func (rt *Racetrack) getDistancesFrom(r, c int) [][]int {
	distances := [][]int{}
	for rr := 0; rr < rt.m; rr++ {
		row := make([]int, rt.n)
		for i := 0; i < rt.n; i++ {
			row[i] = math.MaxInt
		}
		distances = append(distances, row)
	}

	visited := [][]bool{}
	for row := 0; row < rt.m; row++ {
		visited = append(visited, make([]bool, rt.n))
	}
	visited[r][c] = true

	q := [][2]int{{r, c}}

	for dist := 0; len(q) > 0; dist++ {
		for i := 0; i < len(q); i++ {
			cur := q[0]
			q = q[1:]

			distances[cur[R]][cur[C]] = dist

			for _, dir := range DIRS {
				nr := cur[R] + dir[R]
				nc := cur[C] + dir[C]

				if nr < 0 || nr >= rt.m || nc < 0 || nc >= rt.n || rt.grid[nr][nc] == "#" || visited[nr][nc] {
					continue
				}

				visited[nr][nc] = true
				q = append(q, [2]int{nr, nc})
			}
		}
	}

	return distances
}

func abs(v int) int {
	if v < 0 {
		return -v
	}
	return v
}

func (rt *Racetrack) getShortcutsCount(cheats_allowed int, min_savings int) map[int]int {
	shortcuts := make(map[int]int)
	from_start := rt.getDistancesFrom(rt.start[R], rt.start[C])
	from_end := rt.getDistancesFrom(rt.end[R], rt.end[C])
	shortest := from_start[rt.end[R]][rt.end[C]]

	explore_shortcuts := func(start_row, start_col int) {
		// check every tile reachable with cheats (aka within manhattan distance of cheats_allowed)
		for rd := -cheats_allowed; rd <= cheats_allowed; rd++ {
			diff := cheats_allowed - abs(rd)
			for cd := -diff; cd <= diff; cd++ {
				end_row := start_row + rd
				end_col := start_col + cd
				shortcut_len := abs(end_row-start_row) + abs(end_col-start_col) // manhattan distance

				if end_row >= 0 && end_row < rt.m &&
					end_col >= 0 && end_col < rt.n &&
					rt.grid[end_row][end_col] != "#" {
					time := from_start[start_row][start_col] + shortcut_len + from_end[end_row][end_col]
					saved := shortest - time
					if saved >= min_savings {
						shortcuts[saved] += 1
					}
				}
			}
		}
	}

	q := [][2]int{rt.start}
	visited := [][]bool{}
	for row := 0; row < rt.m; row++ {
		visited = append(visited, make([]bool, rt.n))
	}
	visited[rt.start[R]][rt.start[C]] = true
	for len(q) > 0 {
		for i := 0; i < len(q); i++ {
			cur := q[0]
			q = q[1:]
			cr := cur[R]
			cc := cur[C]

			explore_shortcuts(cr, cc)

			for _, dir := range DIRS {
				nr := cr + dir[R]
				nc := cc + dir[C]
				if nr < 0 || nr >= rt.m || nc < 0 || nc >= rt.n ||
					rt.grid[nr][nc] == "#" || visited[nr][nc] {
					continue
				}
				visited[nr][nc] = true
				q = append(q, [2]int{nr, nc})
			}
		}
	}

	return shortcuts
}

//
//
//

func main() {
	input_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer input_file.Close()

	racetrack := parseInput(input_file)

	// part 1
	cheats_allowed_1 := 2
	min_saving_1 := 100

	min_time_no_cheating := racetrack.getShortestTime()
	fmt.Println("shortest time without cheating:", min_time_no_cheating)

	// cheats_savings := racetrack.getCheatsSavings(cheats_allowed_1)
	cheats_savings := racetrack.getShortcutsCount(cheats_allowed_1, min_saving_1)
	// fmt.Printf("\n%+v\n", cheats_savings)

	good_cheats_count := 0
	for k, v := range cheats_savings {
		if k >= min_saving_1 {
			good_cheats_count += v
		}
	}
	fmt.Printf("number of cheats (len %v) that save at least %vps: %v\n", cheats_allowed_1, min_saving_1, good_cheats_count)

	// part 2
	cheats_allowed_2 := 20
	min_saving_2 := 100

	shortcuts := racetrack.getShortcutsCount(20, 50)

	shortcuts_count := 0
	for k, v := range shortcuts {
		if k >= min_saving_2 {
			shortcuts_count += v
		}
	}
	// fmt.Printf("\n%+v\n", shortcuts)
	fmt.Printf("number of cheats (len %v) that save at least %vps: %v\n", cheats_allowed_2, min_saving_2, shortcuts_count)
}

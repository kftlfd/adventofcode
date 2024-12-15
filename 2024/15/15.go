package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	R     = 0
	C     = 1
	UP    = 0
	DOWN  = 1
	LEFT  = 2
	RIGHT = 3
)

var DIRS = [4][2]int{
	{-1, 0}, // ^ 0
	{1, 0},  // V 1
	{0, -1}, // < 2
	{0, 1},  // > 3
}

func parseMove(m string) int {
	switch m {
	case "^":
		return UP
	case "v":
		return DOWN
	case "<":
		return LEFT
	case ">":
		return RIGHT
	}
	panic(fmt.Errorf("move not recognized: %v", m))
}

func getInput(file *os.File) (grid [][]string, moves []int) {
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)

	grid = [][]string{}
	moves = []int{}
	is_moves := false

	for reader.Scan() {
		line := reader.Text()

		if len(line) < 1 {
			is_moves = true
			continue
		}

		if is_moves {
			for _, m := range strings.Split(line, "") {
				moves = append(moves, parseMove(m))
			}
			continue
		}

		grid = append(grid, strings.Split(line, ""))
	}

	return grid, moves
}

func getWideGrid(grid [][]string) [][]string {
	grid_wide := [][]string{}

	for _, row := range grid {
		row_wide := []string{}
		for _, c := range row {
			switch c {
			case "#":
				row_wide = append(row_wide, "#", "#")
			case ".":
				row_wide = append(row_wide, ".", ".")
			case "O":
				row_wide = append(row_wide, "[", "]")
			case "@":
				row_wide = append(row_wide, "@", ".")
			default:
				panic(fmt.Errorf("tile not recognized: %v", c))
			}
		}
		grid_wide = append(grid_wide, row_wide)
	}

	return grid_wide
}

func printGrid(grid [][]string) {
	for _, row := range grid {
		fmt.Println(row)
	}
}

func main() {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	grid, moves := getInput(inp_file)
	grid_wide := getWideGrid(grid)

	wh := newWarehouse(grid)
	whw := newWarehouseWide(grid_wide)

	// printGrid(whw.grid)

	for _, m := range moves {
		wh.move(m)
		whw.move(m)
	}

	// printGrid(whw.grid)

	fmt.Println("coordinates sum:", wh.coordinatesSum())
	fmt.Println("coordinates sum (wide):", whw.coordinatesSum())
}

func findRobot(grid [][]string) (r, c int) {
	for r, row := range grid {
		for c, cell := range row {
			if cell == "@" {
				return r, c
			}
		}
	}
	panic(fmt.Errorf("robot not found"))
}

type Warehouse struct {
	grid                   [][]string
	m, n, robot_r, robot_c int
}

func newWarehouse(grid [][]string) Warehouse {
	rr, rc := findRobot(grid)
	return Warehouse{grid, len(grid), len(grid[0]), rr, rc}
}

func (w *Warehouse) move(m int) {
	dir := DIRS[m]
	nr := w.robot_r + dir[R]
	nc := w.robot_c + dir[C]

	if w.grid[nr][nc] == "#" {
		return
	}

	if w.grid[nr][nc] == "." {
		w.grid[w.robot_r][w.robot_c] = "."
		w.grid[nr][nc] = "@"
		w.robot_r = nr
		w.robot_c = nc
		return
	}

	if w.grid[nr][nc] != "O" {
		printGrid(w.grid)
		panic(fmt.Errorf("grid tile not recognized: %v [%v, %v]", w.grid[nr][nc], nr, nc))
	}

	// find free tile in direction
	found_free_tile := false
	f_r := nr + dir[R]
	f_c := nc + dir[C]

findFree:
	for true {
		switch w.grid[f_r][f_c] {
		case ".":
			found_free_tile = true
			break findFree
		case "O":
			f_r += dir[R]
			f_c += dir[C]
		case "#":
			break findFree
		default:
			panic(fmt.Errorf("tile not recognized: %v [%v, %v]", w.grid[f_r][f_c], f_r, f_c))
		}
	}

	if !found_free_tile {
		return
	}

	w.grid[f_r][f_c] = "O"
	w.grid[w.robot_r][w.robot_c] = "."
	w.grid[nr][nc] = "@"
	w.robot_r = nr
	w.robot_c = nc
}

func (w *Warehouse) coordinatesSum() int {
	coordinates_sum := 0
	for r := 0; r < w.m; r++ {
		for c := 0; c < w.n; c++ {
			if w.grid[r][c] == "O" {
				coordinates_sum += 100*r + c
			}
		}
	}
	return coordinates_sum
}

//
// part 2
//

type WarehouseWide struct {
	grid                   [][]string
	m, n, robot_r, robot_c int
}

func newWarehouseWide(grid [][]string) WarehouseWide {
	rr, rc := findRobot(grid)
	return WarehouseWide{grid, len(grid), len(grid[0]), rr, rc}
}

func (w *WarehouseWide) move(m int) {
	dir := DIRS[m]
	nr := w.robot_r + dir[R]
	nc := w.robot_c + dir[C]

	if w.grid[nr][nc] == "#" {
		return
	}

	if w.grid[nr][nc] == "." {
		w.grid[w.robot_r][w.robot_c] = "."
		w.grid[nr][nc] = "@"
		w.robot_r = nr
		w.robot_c = nc
		return
	}

	if w.grid[nr][nc] != "[" && w.grid[nr][nc] != "]" {
		panic(fmt.Errorf("tile not recognized: %v", w.grid[nr][nc]))
	}

	if m == LEFT || m == RIGHT {
		if m == LEFT && w.grid[nr][nc] != "]" {
			panic(fmt.Errorf("bad state: moving left not into ']'"))
		}
		if m == RIGHT && w.grid[nr][nc] != "[" {
			panic(fmt.Errorf("bad state: moving right not into '['"))
		}

		found_free_tile := false
		f_c := nc + dir[C]

	findFree:
		for true {
			switch w.grid[nr][f_c] {
			case ".":
				found_free_tile = true
				break findFree
			case "[":
				f_c += dir[C]
			case "]":
				f_c += dir[C]
			case "#":
				break findFree
			default:
				panic(fmt.Errorf("tile not recognized: %v [%v, %v]", w.grid[nr][f_c], nr, f_c))
			}
		}

		if !found_free_tile {
			return
		}

		if m == LEFT {
			i := 0
			for true {
				if w.grid[nr][f_c+i] == "@" {
					break
				}
				if i%2 == 0 {
					w.grid[nr][f_c+i] = "["
				} else {
					w.grid[nr][f_c+i] = "]"
				}
				i += 1
			}
		} else {
			i := 0
			for true {
				if w.grid[nr][f_c-i] == "@" {
					break
				}
				if i%2 == 0 {
					w.grid[nr][f_c-i] = "]"
				} else {
					w.grid[nr][f_c-i] = "["
				}
				i += 1
			}
		}

		w.grid[w.robot_r][w.robot_c] = "."
		w.grid[nr][nc] = "@"
		w.robot_r = nr
		w.robot_c = nc
		return
	}

	// up, down
	canMove := true
	if w.grid[nr][nc] == "[" {
		canMove = w.checkCanMove(m, nr, nc)
	} else if w.grid[nr][nc] == "]" {
		canMove = w.checkCanMove(m, nr, nc-1)
	} else {
		panic(fmt.Errorf("move: bad state"))
	}

	if !canMove {
		return
	}

	if w.grid[nr][nc] == "[" {
		w.tryMove(m, nr, nc)
	} else if w.grid[nr][nc] == "]" {
		w.tryMove(m, nr, nc-1)
	}

	w.grid[w.robot_r][w.robot_c] = "."
	w.grid[nr][nc] = "@"
	w.robot_r = nr
	w.robot_c = nc
}

func (w *WarehouseWide) checkCanMove(m, r, c int) bool {
	if m != UP && m != DOWN {
		panic(fmt.Errorf("checkCanMove: direction is not UP or DOWN"))
	}
	if w.grid[r][c] != "[" {
		panic(fmt.Errorf("checkCanMove: bad position"))
	}
	nr := r + DIRS[m][R]

	if w.grid[nr][c] == "." && w.grid[nr][c+1] == "." {
		return true
	}

	if w.grid[nr][c] == "#" || w.grid[nr][c+1] == "#" {
		return false
	}

	if w.grid[nr][c] == "[" {
		return w.checkCanMove(m, nr, c)
	}

	if !((w.grid[nr][c] == "." && w.grid[nr][c+1] == "[") ||
		(w.grid[nr][c] == "]" && w.grid[nr][c+1] == ".") ||
		(w.grid[nr][c] == "]" && w.grid[nr][c+1] == "[")) {
		panic(fmt.Errorf("checkCanMove: bad state: next block move"))
	}

	can_move := true

	if w.grid[nr][c] == "]" {
		can_move = can_move && w.checkCanMove(m, nr, c-1)
	}
	if w.grid[nr][c+1] == "[" {
		can_move = can_move && w.checkCanMove(m, nr, c+1)
	}

	return can_move
}

func (w *WarehouseWide) tryMove(m, r, c int) {
	if m != UP && m != DOWN {
		panic(fmt.Errorf("tryMove: direction is not UP or DOWN"))
	}
	if w.grid[r][c] != "[" {
		panic(fmt.Errorf("tryMove: bad position"))
	}
	nr := r + DIRS[m][R]

	if w.grid[nr][c] == "#" || w.grid[nr][c+1] == "#" {
		panic(fmt.Errorf("tryMove: trying to move into wall"))
	}

	if w.grid[nr][c] == "." && w.grid[nr][c+1] == "." {
		w.grid[r][c] = "."
		w.grid[r][c+1] = "."
		w.grid[nr][c] = "["
		w.grid[nr][c+1] = "]"
		return
	}

	if w.grid[nr][c] == "[" {
		w.tryMove(m, nr, c)
	}

	if w.grid[nr][c] == "]" {
		w.tryMove(m, nr, c-1)
	}

	if w.grid[nr][c+1] == "[" {
		w.tryMove(m, nr, c+1)
	}

	w.grid[r][c] = "."
	w.grid[r][c+1] = "."
	w.grid[nr][c] = "["
	w.grid[nr][c+1] = "]"
}

func (w *WarehouseWide) coordinatesSum() int {
	coordinates_sum := 0
	for r := 0; r < w.m; r++ {
		for c := 0; c < w.n; c++ {
			if w.grid[r][c] == "[" {
				coordinates_sum += 100*r + c
			}
		}
	}
	return coordinates_sum
}

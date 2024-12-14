package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	X           = 0
	Y           = 1
	GRID_HEIGHT = 103
	GRID_WIDTH  = 101
	SECONDS     = 100
)

func main() {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	reader := bufio.NewScanner(inp_file)
	reader.Split(bufio.ScanLines)

	robots := []Robot{}

	for reader.Scan() {
		robots = append(robots, parseRobot(reader.Text()))
	}

	grid := [][]int{}
	for h := 0; h < GRID_HEIGHT; h++ {
		grid = append(grid, make([]int, GRID_WIDTH))
	}

	for _, robot := range robots {
		x := (((robot.start[X] + robot.v[X]*SECONDS) % GRID_WIDTH) + GRID_WIDTH) % GRID_WIDTH
		y := (((robot.start[Y] + robot.v[Y]*SECONDS) % GRID_HEIGHT) + GRID_HEIGHT) % GRID_HEIGHT
		grid[y][x] += 1
	}

	grid_x_mid := GRID_WIDTH / 2
	grid_y_mid := GRID_HEIGHT / 2

	quadrants := []int{
		gridSum(grid, 0, grid_x_mid, 0, grid_y_mid),
		gridSum(grid, grid_x_mid+1, GRID_WIDTH, 0, grid_y_mid),
		gridSum(grid, 0, grid_x_mid, grid_y_mid+1, GRID_HEIGHT),
		gridSum(grid, grid_x_mid+1, GRID_WIDTH, grid_y_mid+1, GRID_HEIGHT),
	}

	safety_factor := 1
	for _, q := range quadrants {
		safety_factor *= q
	}

	fmt.Println("safety factor:", safety_factor)

	inputScanner := bufio.NewScanner(os.Stdin)
	// i := 11
	i := 7687
	for true {
		printGridAfterSeconds(robots, i)
		i += 101
		inputScanner.Scan()
	}
}

type Robot struct {
	start, v [2]int
}

func parseInt(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		panic(err)
	}
	return val
}

func parseRobot(str string) Robot {
	row := strings.Split(str, " ")
	if len(row) != 2 {
		panic(fmt.Errorf("bad robot input: %v", row))
	}

	p := strings.Split(row[0], ",")
	v := strings.Split(row[1], ",")

	px := parseInt(p[0][2:])
	py := parseInt(p[1])
	vx := parseInt(v[0][2:])
	vy := parseInt(v[1])

	return Robot{start: [2]int{px, py}, v: [2]int{vx, vy}}
}

func gridSum(grid [][]int, x_start, x_end, y_start, y_end int) int {
	sum := 0
	for y := y_start; y < y_end; y++ {
		for x := x_start; x < x_end; x++ {
			sum += grid[y][x]
		}
	}
	return sum
}

func printGridAfterSeconds(robots []Robot, seconds int) {
	grid := [][]int{}
	for r := 0; r < GRID_HEIGHT; r++ {
		grid = append(grid, make([]int, GRID_WIDTH))
	}

	for _, robot := range robots {
		x := (((robot.start[X] + robot.v[X]*seconds) % GRID_WIDTH) + GRID_WIDTH) % GRID_WIDTH
		y := (((robot.start[Y] + robot.v[Y]*seconds) % GRID_HEIGHT) + GRID_HEIGHT) % GRID_HEIGHT
		grid[y][x] += 1
	}

	fmt.Printf("\n\nafter %v seconds\n", seconds)
	for r := 0; r < GRID_HEIGHT; r++ {
		for c := 0; c < GRID_WIDTH; c++ {
			if grid[r][c] == 0 {
				fmt.Printf(".")
			} else {
				fmt.Printf("0")
			}
		}
		fmt.Printf("\n")
	}
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	inp_path := os.Args[1]
	inp_file, err := os.Open(inp_path)
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	reader := bufio.NewScanner(inp_file)
	reader.Split(bufio.ScanLines)

	grid := [][]string{}

	for reader.Scan() {
		grid = append(grid, strings.Split(reader.Text(), ""))
	}

	// for _, row := range grid {
	// 	fmt.Println(row)
	// }

	m := len(grid)
	n := len(grid[0])

	garden := Garden{garden: grid, m: m, n: n}

	regions := newRegions(m, n)

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if regions.regions[r][c] > 0 {
				continue
			}
			markRegion(&garden, r, c, &regions)
		}
	}

	// fmt.Println("---stats---")
	// for i, stat := range regions.stats {
	// 	fmt.Println(i, stat)
	// }
	// fmt.Println("---regions---")
	// for _, row := range regions.regions {
	// 	fmt.Println(row)
	// }

	fencing_price := 0
	fencing_sides_price := 0
	for _, stat := range regions.stats {
		fencing_price += stat[0] * stat[1]
		fencing_sides_price += stat[0] * stat[2]
	}
	fmt.Println("regions fencing price:", fencing_price)
	fmt.Println("sides fencing price:", fencing_sides_price)
}

type Garden struct {
	garden [][]string
	m, n   int
}

type Regions struct {
	stats   [][3]int // stats[i] = [area_i, perimeter_i, sides_i]
	regions [][]int
	m, n    int
	last_i  int
}

func newRegions(m, n int) Regions {
	regions_stats := [][3]int{{0, 0, 0}}

	regions_grid := [][]int{}
	for r := 0; r < m; r++ {
		regions_grid = append(regions_grid, make([]int, n))
	}

	regions := Regions{stats: regions_stats, regions: regions_grid, m: m, n: n, last_i: 0}

	return regions
}

type dfsAns struct{ area, perimeter int }

func markRegion(g *Garden, row, col int, r *Regions) {
	cur_region := r.last_i + 1
	r.last_i = cur_region

	ans := dfsAns{}

	dfs(g, row, col, r, g.garden[row][col], cur_region, &ans)

	sides := regionSides(r, cur_region)

	r.stats = append(r.stats, [3]int{ans.area, ans.perimeter, sides})
}

var DIRS = [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}

func dfs(g *Garden, row, col int, r *Regions, cur_plant string, cur_region int, ans *dfsAns) {
	r.regions[row][col] = cur_region

	visited_neibs := 0
	for _, dir := range DIRS {
		nr := row + dir[0]
		nc := col + dir[1]
		if nr < 0 || nr >= g.m || nc < 0 || nc >= g.n {
			continue
		}
		if r.regions[nr][nc] == cur_region {
			visited_neibs += 1
		}
	}

	ans.area += 1
	ans.perimeter += 4 - (visited_neibs * 2)

	for _, dir := range DIRS {
		nr := row + dir[0]
		nc := col + dir[1]
		if nr < 0 || nr >= g.m || nc < 0 || nc >= g.n {
			continue
		}
		if g.garden[nr][nc] != cur_plant {
			continue
		}
		if r.regions[nr][nc] != cur_region {
			dfs(g, nr, nc, r, cur_plant, cur_region, ans)
		}
	}
}

func regionSides(r *Regions, cur_region int) int {
	top := 0
	bottom := 0
	left := 0
	right := 0

	for row := 0; row < r.m; row++ {
		top_last := -2
		bottom_last := -2

		for col := 0; col < r.n; col++ {
			if r.regions[row][col] != cur_region {
				continue
			}

			is_top := row == 0 || r.regions[row-1][col] != cur_region
			is_bottom := row == r.n-1 || r.regions[row+1][col] != cur_region

			if is_top {
				if top_last != col-1 {
					top += 1
				}
				top_last = col
			}
			if is_bottom {
				if bottom_last != col-1 {
					bottom += 1
				}
				bottom_last = col
			}
		}
	}

	for col := 0; col < r.n; col++ {
		left_last := -2
		right_last := -2

		for row := 0; row < r.m; row++ {
			if r.regions[row][col] != cur_region {
				continue
			}

			is_left := col == 0 || r.regions[row][col-1] != cur_region
			is_right := col == r.n-1 || r.regions[row][col+1] != cur_region

			if is_left {
				if left_last != row-1 {
					left += 1
				}
				left_last = row
			}
			if is_right {
				if right_last != row-1 {
					right += 1
				}
				right_last = row
			}
		}
	}

	return top + bottom + left + right
}

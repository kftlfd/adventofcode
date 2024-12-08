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

	input := [][]string{}
	for fileScanner.Scan() {
		row := strings.Split(fileScanner.Text(), "")
		input = append(input, row)
	}

	m := len(input)
	n := len(input[0])

	antennas := make(map[string][][2]int)
	for r, row := range input {
		for c, val := range row {
			if val == "." {
				continue
			}
			locations, ok := antennas[val]
			if !ok {
				locations = [][2]int{}
			}
			locations = append(locations, [2]int{r, c})
			antennas[val] = locations
		}
	}

	antinodes := [][]int{}
	for r := 0; r < m; r++ {
		antinodes = append(antinodes, make([]int, n))
	}

	for _, locations := range antennas {
		for i, l1 := range locations {
			for j := i + 1; j < len(locations); j++ {
				l2 := locations[j]
				for _, antinode := range getAntinodes(l1, l2) {
					r, c := antinode[0], antinode[1]
					if r < 0 || r >= m || c < 0 || c >= n {
						continue
					}
					antinodes[r][c] = 1
				}
			}
		}
	}

	unique_antinode_locations_count := 0
	for _, row := range antinodes {
		for _, val := range row {
			unique_antinode_locations_count += val
		}
	}
	fmt.Println("number of unique antinode locations:", unique_antinode_locations_count)

	// part 2

	antinodes = [][]int{}
	for r := 0; r < m; r++ {
		antinodes = append(antinodes, make([]int, n))
	}

	for _, locations := range antennas {
		for i, l1 := range locations {
			for j := i + 1; j < len(locations); j++ {
				l2 := locations[j]
				for _, antinode := range getAntinodes2(l1, l2, m, n) {
					antinodes[antinode[0]][antinode[1]] = 1
				}
			}
		}
	}

	unique_antinode_locations_count_2 := 0
	for _, row := range antinodes {
		for _, val := range row {
			unique_antinode_locations_count_2 += val
		}
	}
	fmt.Println("number of unique antinode w resonance:", unique_antinode_locations_count_2)

}

func getAntinodes(l1, l2 [2]int) [2][2]int {
	dr := l2[0] - l1[0]
	dc := l2[1] - l1[1]
	return [2][2]int{
		{l1[0] - dr, l1[1] - dc},
		{l2[0] + dr, l2[1] + dc},
	}
}

func getAntinodes2(l1, l2 [2]int, m, n int) [][2]int {
	ans := [][2]int{}
	dr := l2[0] - l1[0]
	dc := l2[1] - l1[1]

	i := 0
	for true {
		nr := l1[0] - (dr * i)
		nc := l1[1] - (dc * i)
		if nr < 0 || nc < 0 || nr >= m || nc >= n {
			break
		}
		ans = append(ans, [2]int{nr, nc})
		i += 1
	}

	i = 0
	for true {
		nr := l2[0] + (dr * i)
		nc := l2[1] + (dc * i)
		if nr < 0 || nc < 0 || nr >= m || nc >= n {
			break
		}
		ans = append(ans, [2]int{nr, nc})
		i += 1
	}

	return ans
}

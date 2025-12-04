package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func parseInput() [][]string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		inp = append(inp, strings.Split(line, ""))
	}

	return inp
}

func main() {
	input := parseInput()

	// Part 1
	accessible := 0
	for r, row := range input {
		for c, cur := range row {
			if cur != "@" {
				continue
			}

			touching := -1

			for rr := r - 1; rr <= r+1; rr++ {
				if rr < 0 || rr >= len(input) {
					continue
				}
				for cc := c - 1; cc <= c+1; cc++ {
					if cc < 0 || cc >= len(row) {
						continue
					}
					if input[rr][cc] == "@" {
						touching += 1
					}
				}
			}

			if touching < 4 {
				accessible += 1
			}
		}
	}
	fmt.Println("Part 1:", accessible)

	// Part 2
	removed := 0
	changes := true
	for changes {
		changes = false
		toRemove := [][2]int{}
		for r, row := range input {
			for c, cur := range row {
				if cur != "@" {
					continue
				}
				touching := 0
				for rr := r - 1; rr <= r+1; rr++ {
					if rr < 0 || rr >= len(input) {
						continue
					}
					for cc := c - 1; cc <= c+1; cc++ {
						if cc < 0 || cc >= len(row) {
							continue
						}
						if rr == r && cc == c {
							continue
						}
						if input[rr][cc] == "@" {
							touching += 1
						}
					}
				}
				if touching < 4 {
					toRemove = append(toRemove, [2]int{r, c})
				}
			}
		}
		if len(toRemove) > 0 {
			for _, roll := range toRemove {
				input[roll[0]][roll[1]] = "."
			}
			changes = true
			removed += len(toRemove)
		}
	}
	fmt.Println("Part 2:", removed)
}

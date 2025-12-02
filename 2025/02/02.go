package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput() [][2]int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := [][2]int{}

	for scanner.Scan() {
		line := scanner.Text()
		ranges := strings.Split(line, ",")

		for _, cur := range ranges {
			cur_range := strings.Split(cur, "-")

			from, err := strconv.Atoi(cur_range[0])
			if err != nil {
				panic(err)
			}

			to, err := strconv.Atoi(cur_range[1])
			if err != nil {
				panic((err))
			}

			inp = append(inp, [2]int{from, to})
		}
	}

	return inp
}

func findInvalidIdsP1(from, to int) []int {
	out := []int{}
	end := to + 1
	for cur := from; cur < end; cur++ {
		idString := strconv.Itoa(cur)
		n := len(idString) / 2
		if idString[:n] == idString[n:] {
			out = append(out, cur)
		}
	}
	return out
}

func checkRepeats(str string) bool {
	l := len(str)
	maxN := l / 2
	for n := 1; n <= maxN; n++ {
		if l%n != 0 {
			continue
		}
		repeatFound := true
		for i := n * 2; i <= l; i += n {
			if str[i-n*2:i-n] != str[i-n:i] {
				repeatFound = false
				break
			}
		}
		if repeatFound {
			return true
		}
	}
	return false
}

func findInvalidIdsP2(from, to int) []int {
	out := []int{}
	end := to + 1

	for cur := from; cur < end; cur++ {
		idString := strconv.Itoa(cur)
		if checkRepeats(idString) {
			out = append(out, cur)
		}
	}

	return out
}

func main() {
	input := parseInput()

	// Part 1
	invalidSum := 0
	for _, pair := range input {
		invalidIds := findInvalidIdsP1(pair[0], pair[1])
		for _, id := range invalidIds {
			invalidSum += id
		}
	}
	fmt.Println("Part 1:", invalidSum)

	// Part 2
	invalidSum = 0
	for _, pair := range input {
		invalidIds := findInvalidIdsP2(pair[0], pair[1])
		for _, id := range invalidIds {
			invalidSum += id
		}
	}
	fmt.Println("Part 2:", invalidSum)
}

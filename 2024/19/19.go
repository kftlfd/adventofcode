package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Towels struct {
	list                    []string
	is_possible_memo        map[string]bool
	arrangements_count_memo map[string]int
}

func getInput(file *os.File) (Towels, []string) {
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	memo := make(map[string]bool)
	memo[""] = true

	memo_2 := make(map[string]int)

	towels := Towels{list: make([]string, 0), is_possible_memo: memo, arrangements_count_memo: memo_2}
	designs := []string{}

	is_designs := false

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < 1 {
			is_designs = true
			continue
		}

		if is_designs {
			designs = append(designs, line)
			continue
		}

		towels.list = strings.Split(line, ", ")
	}

	return towels, designs
}

//
// part 1
//

func (t *Towels) isDesignPossible(design string) bool {
	is_possible, ok := t.is_possible_memo[design]
	if ok {
		return is_possible
	}

	for _, towel := range t.list {
		n := len(towel)
		if n > len(design) {
			continue
		}
		if design[:n] == towel && t.isDesignPossible(design[n:]) {
			t.is_possible_memo[design] = true
			return true
		}
	}

	t.is_possible_memo[design] = false
	return false
}

//
// part 2
//

func (t *Towels) arrangementsCount(design string) int {
	arrs, ok := t.arrangements_count_memo[design]
	if ok {
		return arrs
	}

	design_arrangements := 0

	for _, towel := range t.list {
		n := len(towel)
		if n > len(design) || design[:n] != towel {
			continue
		}

		if n == len(design) {
			design_arrangements += 1
		} else {
			design_arrangements += t.arrangementsCount(design[n:])
		}
	}

	t.arrangements_count_memo[design] = design_arrangements
	return design_arrangements
}

//
//
//

func main() {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	towels, designs := getInput(inp_file)

	possible_designs_count := 0
	possible_arrangements_count := 0

	for _, d := range designs {
		if towels.isDesignPossible(d) {
			possible_designs_count += 1
		}
		possible_arrangements_count += towels.arrangementsCount(d)
	}

	fmt.Println("designs possible:", possible_designs_count)
	fmt.Println("arrangements count:", possible_arrangements_count)
}

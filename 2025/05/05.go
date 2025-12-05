package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type parsedInput struct {
	freshIdsRanges [][2]int
	ingredientsIds []int
}

func parseInput() parsedInput {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	ranges := [][2]int{}
	ids := []int{}
	isIds := false

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			isIds = true
			continue
		}

		if isIds {
			curId, err := strconv.Atoi(line)
			if err != nil {
				panic(err)
			}
			ids = append(ids, curId)
			continue
		}

		rangeStrArr := strings.Split(line, "-")
		if len(rangeStrArr) != 2 {
			panic(fmt.Errorf("invalid ids range"))
		}
		from, err := strconv.Atoi(rangeStrArr[0])
		if err != nil {
			panic(err)
		}
		to, err := strconv.Atoi(rangeStrArr[1])
		if err != nil {
			panic(err)
		}
		ranges = append(ranges, [2]int{from, to})
	}

	return parsedInput{freshIdsRanges: ranges, ingredientsIds: ids}
}

func main() {
	input := parseInput()

	// Part 1
	freshCount := 0
	for _, id := range input.ingredientsIds {
		for _, curRange := range input.freshIdsRanges {
			if id >= curRange[0] && id <= curRange[1] {
				freshCount += 1
				break
			}
		}
	}
	fmt.Println("Part 1:", freshCount)

	// Part 2
	totalFreshCount := 0
	// sort ranges
	sortedIds := append([][2]int{}, input.freshIdsRanges...)
	sort.Slice(sortedIds, func(i, j int) bool {
		return sortedIds[i][0] < sortedIds[j][0]
	})
	// remove overlaps
	cleanRanges := [][2]int{sortedIds[0]}
	for i := 1; i < len(sortedIds); i++ {
		curRange := sortedIds[i]
		last := len(cleanRanges) - 1
		if curRange[0] > cleanRanges[last][1] {
			cleanRanges = append(cleanRanges, curRange)
		} else if curRange[1] > cleanRanges[last][1] {
			cleanRanges[last][1] = curRange[1]
		}
	}
	// sum up counts in ranges
	for _, curRange := range cleanRanges {
		totalFreshCount += curRange[1] - curRange[0] + 1
	}
	fmt.Println("Part 2:", totalFreshCount)
}

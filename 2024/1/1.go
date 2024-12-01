package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	// prepare two slices
	list1 := []int{}
	list2 := []int{}

	// read input file by rows
	// parse values from each column and add to the slices
	inp, err := os.ReadFile(os.Args[1])
	if err != nil {
		panic(err)
	}
	rows := strings.Split(string(inp), "\n")
	for _, row := range rows {
		vals := strings.Split(row, "   ")
		if len(vals) != 2 {
			continue
		}
		val1, err := strconv.Atoi(vals[0])
		val2, err := strconv.Atoi(vals[1])
		if err != nil {
			panic(err)
		}
		list1 = append(list1, val1)
		list2 = append(list2, val2)
	}

	// sort slices
	sort.Slice(list1, func(i, j int) bool {
		return list1[i] < list1[j]
	})
	sort.Slice(list2, func(i, j int) bool {
		return list2[i] < list2[j]
	})

	// calculate total abs difference of lists
	totalDiff := 0
	for i, val1 := range list1 {
		diff := list2[i] - val1
		if diff < 0 {
			diff *= -1
		}
		totalDiff += diff
	}

	fmt.Println("distances diff:\t\t", totalDiff)

	// part 2
	cnt := make(map[int]int)
	similarityScore := 0
	for _, num1 := range list1 {
		count, ok := cnt[num1]
		if !ok {
			count = 0
			for _, num2 := range list2 {
				if num1 == num2 {
					count += 1
				}
			}
			cnt[num1] = count
		}
		similarityScore += num1 * count
	}

	fmt.Println("similarity score:\t", similarityScore)
}

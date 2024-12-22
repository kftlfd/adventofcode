package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func parseInput() []int {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	nums := []int{}

	for scanner.Scan() {
		num, err := strconv.Atoi(scanner.Text())
		if err != nil {
			panic(err)
		}
		nums = append(nums, num)
	}

	return nums
}

func mix(a, b int) int {
	return a ^ b
}

func prune(a int) int {
	return a % 16777216
}

func evolve(num int) int {
	num = prune(mix(num, num*64))
	num = prune(mix(num, num/32))
	num = prune(mix(num, num*2048))
	return num
}

func evolveTimes(num, times int) int {
	for i := 0; i < times; i++ {
		num = evolve(num)
	}
	return num
}

func getPrice(num int) int {
	return num % 10
}

func getKey(change int) int {
	return change + 9
}

func getMemo[T int | bool]() [19][19][19][19]T {
	return [19][19][19][19]T{}
}

func evolveAndRecordPriceAtChanges(num, times int, bananas *[19][19][19][19]int, seen *[19][19][19][19]bool) {
	prev_price := getPrice(num)
	price := prev_price

	changes_sequence := [4]int{}

	// first 4 changes
	for i := 0; i < 4; i++ {
		num = evolve(num)
		prev_price = price
		price = getPrice(num)
		change := price - prev_price
		changes_sequence[i] = change
	}
	k1 := getKey(changes_sequence[0])
	k2 := getKey(changes_sequence[1])
	k3 := getKey(changes_sequence[2])
	k4 := getKey(changes_sequence[3])
	(*seen)[k1][k2][k3][k4] = true
	(*bananas)[k1][k2][k3][k4] += price

	// remaining changes
	for i := 4; i < times; i++ {
		num = evolve(num)
		prev_price = price
		price = getPrice(num)
		change := price - prev_price

		changes_sequence[0] = changes_sequence[1]
		changes_sequence[1] = changes_sequence[2]
		changes_sequence[2] = changes_sequence[3]
		changes_sequence[3] = change

		k1 := getKey(changes_sequence[0])
		k2 := getKey(changes_sequence[1])
		k3 := getKey(changes_sequence[2])
		k4 := getKey(changes_sequence[3])

		if !(*seen)[k1][k2][k3][k4] {
			(*seen)[k1][k2][k3][k4] = true
			(*bananas)[k1][k2][k3][k4] += price
		}
	}
}

func main() {
	// part 1
	nums := parseInput()
	evolved_sum := 0
	for _, num := range nums {
		evolved := evolveTimes(num, 2000)
		evolved_sum += evolved
	}
	fmt.Println("sum:", evolved_sum)

	// part 2
	bananas := getMemo[int]()
	seen := getMemo[bool]()
	for _, num := range nums {
		evolveAndRecordPriceAtChanges(num, 2000, &bananas, &seen)
		// reset seen
		for i := 0; i < 19; i++ {
			for j := 0; j < 19; j++ {
				for k := 0; k < 19; k++ {
					for l := 0; l < 19; l++ {
						seen[i][j][k][l] = false
					}
				}
			}
		}
	}
	max_bananas := 0
	for i := 0; i < 19; i++ {
		for j := 0; j < 19; j++ {
			for k := 0; k < 19; k++ {
				for l := 0; l < 19; l++ {
					max_bananas = max(max_bananas, bananas[i][j][k][l])
				}
			}
		}
	}
	fmt.Println("max bananas:", max_bananas)
}

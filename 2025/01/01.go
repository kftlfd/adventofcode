package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func parseInput() []int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := []int{}

	for scanner.Scan() {
		line := scanner.Text()

		sign := 1
		if line[:1] == "L" {
			sign = -1
		}

		num, err := strconv.Atoi(line[1:])
		if err != nil {
			panic(err)
		}

		inp = append(inp, sign*num)
	}

	return inp
}

func main() {
	input := parseInput()

	// Part 1

	cur := 50
	atZero := 0

	for _, num := range input {
		cur = (cur + num + 100) % 100
		if cur == 0 {
			atZero += 1
		}
	}

	fmt.Println("Part 1:", atZero)

	// Part 2

	cur = 50
	atZero = 0

	for _, num := range input {
		// numAbs := num
		// move := 1
		// if num < 0 {
		// 	numAbs *= -1
		// 	move = -1
		// }

		// for i := 0; i < numAbs; i++ {
		// 	cur = (cur + move + 100) % 100
		// 	if cur == 0 {
		// 		atZero += 1
		// 	}
		// }

		fullRotations := num / 100
		if fullRotations < 0 {
			fullRotations *= -1
		}

		atZero += fullRotations

		num = num % 100

		nxtCur := (cur + num + 100) % 100

		landAtZero := nxtCur == 0

		passZero := cur > 0 && ((cur+num > 100) || (cur+num < 0))

		if landAtZero || passZero {
			atZero += 1
		}

		cur = nxtCur
	}

	fmt.Println("Part 2:", atZero)
}

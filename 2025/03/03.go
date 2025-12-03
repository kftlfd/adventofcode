package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func parseInput() [][]int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := [][]int{}

	for scanner.Scan() {
		line := scanner.Text()
		inp = append(inp, batteryBankStringToIntSlice(line))
	}

	return inp
}

func batteryBankStringToIntSlice(str string) []int {
	n := len(str)
	out := make([]int, n)
	for i, numChar := range strings.Split(str, "") {
		numInt, err := strconv.Atoi(numChar)
		if err != nil {
			panic(err)
		}
		out[i] = numInt
	}
	return out
}

func findTwoMaxBatteries(batteries []int) int {
	n := len(batteries)
	if n < 2 {
		panic(fmt.Errorf("batteries bank len < 2"))
	}

	a := batteries[0]
	b := batteries[1]

	for i := 2; i < n; i++ {
		c := batteries[i]
		if b > a {
			a = b
			b = c
		} else if c > b {
			b = c
		}
	}

	return a*10 + b
}

func intSliceToInt(arr []int) int {
	out := 0
	n := len(arr)
	for i := range n {
		out = (out * 10) + arr[i]
	}
	return out
}

func find12MaxBatteries(batteries []int) int {
	n := len(batteries)
	if n < 12 {
		panic(fmt.Errorf("batteries len < 12"))
	}

	out := make([]int, 12)
	for i := range 12 {
		out[i] = batteries[i]
	}

	for i := 12; i < n; i++ {
		num := batteries[i]

		shiftInMiddle := false
		for j := range 11 {
			if out[j+1] > out[j] {
				newOut := []int{}
				newOut = append(newOut, out[:j]...)
				newOut = append(newOut, out[j+1:]...)
				newOut = append(newOut, num)
				out = newOut
				shiftInMiddle = true
				break
			}
		}
		if shiftInMiddle {
			continue
		}

		if num > out[11] {
			out[11] = num
		}
	}

	return intSliceToInt(out)
}

func main() {
	input := parseInput()

	// Part 1
	totalJoltage := 0
	for _, bank := range input {
		totalJoltage += findTwoMaxBatteries(bank)
	}
	fmt.Println("Part 1:", totalJoltage)

	// Part 2
	totalJoltage = 0
	for _, bank := range input {
		totalJoltage += find12MaxBatteries(bank)
	}
	fmt.Println("Part 2:", totalJoltage)
}

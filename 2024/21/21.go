package main

import (
	"bufio"
	// "day21/try1"
	// "day21/try2"
	// "day21/try3"
	// "day21/try4"
	"day21/try5"
	"fmt"
	"os"
)

func parseInput() []string {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	scanner := bufio.NewScanner(inp_file)
	scanner.Split(bufio.ScanLines)

	numerical_codes := []string{}

	for scanner.Scan() {
		numerical_codes = append(numerical_codes, scanner.Text())
	}

	return numerical_codes
}

func main() {
	num_codes := parseInput()

	// fmt.Println("complexity try 1:", try1.GetComplexity(num_codes, 1))
	// fmt.Println("complexity try 2:", try2.GetComplexity(num_codes, 1))
	// fmt.Println("complexity try 3:", try3.GetComplexity(num_codes, 2))
	// fmt.Println("complexity try 4:", try4.GetComplexity(num_codes, 2))

	fmt.Println("part 1:", try5.GetComplexity(num_codes, 2))
	fmt.Println("part 2:", try5.GetComplexity(num_codes, 25))
}

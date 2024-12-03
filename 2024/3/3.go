package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	inpPath := os.Args[1]
	inpFile, err := os.ReadFile(inpPath)
	if err != nil {
		panic(err)
	}
	re, err := regexp.Compile(`mul\(\d{1,3},\d{1,3}\)`)
	if err != nil {
		panic(err)
	}

	matches := re.FindAllString(string(inpFile), -1)

	if matches == nil {
		fmt.Println("no matches")
		return
	}

	// fmt.Println(matches)

	var mulSum int64

	for _, mul := range matches {
		vals := strings.Split(mul[4:len(mul)-1], ",")
		if len(vals) != 2 {
			panic(fmt.Errorf("invalid mul: %v", mul))
		}
		int1, err1 := strconv.Atoi(vals[0])
		int2, err2 := strconv.Atoi(vals[1])
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("invalid mul: %v: %v, %v", mul, err1, err2))
		}
		mulSum += int64(int1) * int64(int2)
	}

	fmt.Println("uncorrupted mul sum:", mulSum)

	// part 2

	re, err = regexp.Compile(`mul\(\d{1,3},\d{1,3}\)|don't\(\)|do\(\)`)
	if err != nil {
		panic(err)
	}

	matches = re.FindAllString(string(inpFile), -1)

	if matches == nil {
		fmt.Println("no matches")
		return
	}

	// fmt.Println(matches)

	var mulSum2 int64
	disabled := false

	for _, mul := range matches {
		if mul == "do()" {
			disabled = false
			continue
		}
		if mul == "don't()" {
			disabled = true
			continue
		}
		if disabled {
			continue
		}

		vals := strings.Split(mul[4:len(mul)-1], ",")
		if len(vals) != 2 {
			panic(fmt.Errorf("invalid mul: %v", mul))
		}
		int1, err1 := strconv.Atoi(vals[0])
		int2, err2 := strconv.Atoi(vals[1])
		if err1 != nil || err2 != nil {
			panic(fmt.Errorf("invalid mul: %v: %v, %v", mul, err1, err2))
		}
		mulSum2 += int64(int1) * int64(int2)
	}

	fmt.Println("enabled mul sum:", mulSum2)
}

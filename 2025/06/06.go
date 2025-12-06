package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type parsedInput struct {
	lineLen int
	lines   []string
}

func parseInput() parsedInput {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	lineLen := 0
	lines := []string{}

	for scanner.Scan() {
		line := scanner.Text()
		if lineLen == 0 {
			lineLen = len(line)
		}
		if len(line) != lineLen {
			panic(fmt.Errorf("different length of input lines"))
		}
		lines = append(lines, line)
	}

	return parsedInput{lineLen: lineLen, lines: lines}
}

type mathProblem struct {
	nums []int
	op   string
}

func doMath(problems []mathProblem) int {
	totalSum := 0

	for _, probl := range problems {
		curOp := probl.op

		if curOp != "+" && curOp != "*" {
			panic(fmt.Errorf("op not recognized: %s", curOp))
		}

		curResult := 0
		if curOp == "*" {
			curResult = 1
		}

		for _, num := range probl.nums {
			if curOp == "+" {
				curResult += num
			} else {
				curResult *= num
			}
		}

		totalSum += curResult
	}

	return totalSum
}

func getPart1Problems(inp parsedInput) []mathProblem {
	// parse input by rows
	columns := 0
	nums := [][]int{}
	ops := []string{}
	for _, line := range inp.lines {
		fields := strings.Fields(line)
		if columns == 0 {
			columns = len(fields)
		}
		if len(fields) != columns {
			panic(fmt.Errorf("different fields count in input lines"))
		}

		if fields[0] == "+" || fields[0] == "*" {
			ops = fields
			continue
		}

		curNums := []int{}
		for _, numStr := range fields {
			curNum, err := strconv.Atoi(numStr)
			if err != nil {
				panic(err)
			}
			curNums = append(curNums, curNum)
		}
		nums = append(nums, curNums)
	}

	problems := []mathProblem{}

	for i := range columns {
		curOp := ops[i]
		curNums := []int{}
		for _, numLine := range nums {
			curNums = append(curNums, numLine[i])
		}
		problems = append(problems, mathProblem{nums: curNums, op: curOp})
	}

	return problems
}

func getPart2Problems(inp parsedInput) []mathProblem {
	// parse input by columns
	columns := [][]string{}
	lineCount := len(inp.lines)
	for c := range inp.lineLen {
		curCol := []string{}
		for r := range lineCount {
			curCol = append(curCol, inp.lines[r][c:c+1])
		}
		columns = append(columns, curCol)
	}
	// add extra separator-column to the end
	columns = append(columns, []string{" "})

	problems := []mathProblem{}

	curNums := []int{}
	curOp := "."
	for _, col := range columns {
		// fmt.Println(col)

		isSeparator := len(strings.TrimSpace(strings.Join(col, ""))) < 1

		if isSeparator {
			// ignore consecutive separators
			if curOp == "." && len(curNums) < 1 {
				continue
			}
			// check problem validity
			if curOp == "." {
				panic(fmt.Errorf("op not found for problem"))
			}
			if len(curNums) < 1 {
				panic(fmt.Errorf("nums not found for problem"))
			}
			// add curProbl to problems
			problems = append(problems, mathProblem{nums: curNums, op: curOp})
			// reset curProblem
			curNums = []int{}
			curOp = "."
			continue
		}

		lastChar := col[len(col)-1]
		lastCharOp := lastChar == "+" || lastChar == "*"
		if lastCharOp {
			if curOp != "." {
				panic(fmt.Errorf("multiple ops found in problem"))
			}
			curOp = lastChar
		}

		var numStr string
		if lastCharOp {
			numStr = strings.TrimSpace(strings.Join(col[:lineCount-1], ""))
		} else {
			numStr = strings.TrimSpace(strings.Join(col, ""))
		}
		num, err := strconv.Atoi(numStr)
		if err != nil {
			panic(err)
		}
		curNums = append(curNums, num)
	}

	return problems
}

func main() {
	input := parseInput()

	// Part 1
	part1Sum := doMath(getPart1Problems(input))
	fmt.Println("Part 1:", part1Sum)

	// Part 2
	part2Sum := doMath(getPart2Problems(input))
	fmt.Println("Part 2:", part2Sum)
}

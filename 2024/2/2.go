package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	inpPath := os.Args[1]
	inpFile, err := os.Open(inpPath)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(inpFile)
	fileScanner.Split(bufio.ScanLines)

	safeReports := 0
	safeReportsWithTolerate := 0

	for fileScanner.Scan() {
		rowS := strings.Split(fileScanner.Text(), " ")
		if len(rowS) < 2 {
			continue
		}

		row := make([]int, len(rowS))
		for i, val := range rowS {
			intVal, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			row[i] = intVal
		}

		if isSafe(row, 0) {
			safeReports += 1
		}
		if isSafe(row, 1) {
			safeReportsWithTolerate += 1
		}
	}

	fmt.Println("safe reports: ", safeReports)
	fmt.Println("with tolerate:", safeReportsWithTolerate)
}

func isSafe(report []int, tolerate int) bool {
	return isSafeSign(report, 1, tolerate) || isSafeSign(report, -1, tolerate)
}

func isSafeSign(report []int, sign int, tolerate int) bool {
	// edge case: remove the first level
	if tolerate > 0 && isSafeSign(report[1:], sign, tolerate-1) {
		return true
	}
	if len(report) < 1 {
		return true
	}

	safe := true
	prev := report[0]
	for i := 1; i < len(report); i++ {
		diff := (report[i] - prev) * sign
		if diff >= 1 && diff <= 3 {
			prev = report[i]
		} else if tolerate > 0 {
			tolerate -= 1
		} else {
			safe = false
			break
		}
	}
	return safe
}

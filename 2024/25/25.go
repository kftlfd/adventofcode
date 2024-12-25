package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type KeyLock struct {
	height int
	pins   []int
}

func parseInput() (locks, keys []KeyLock) {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	scanner := bufio.NewScanner(inp_file)
	scanner.Split(bufio.ScanLines)

	locks = []KeyLock{}
	keys = []KeyLock{}

	parseKeyLock := func(pend [][]string) {
		height := len(pend)
		width := len(pend[0])

		pins := []int{}

		if pend[0][0] == "#" {
			// is lock
			for col := 0; col < width; col++ {
				for row := 1; row < height; row++ {
					if pend[row][col] == "." {
						pins = append(pins, row)
						break
					}
				}
			}
			locks = append(locks, KeyLock{height: height, pins: pins})
			return
		}

		// is key
		for col := 0; col < len(pend[0]); col++ {
			for row := 0; row < len(pend); row++ {
				if pend[row][col] == "#" {
					pins = append(pins, height-row)
					break
				}
			}
		}
		keys = append(keys, KeyLock{height: height, pins: pins})
	}

	pending := [][]string{}

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			if len(pending) > 0 {
				parseKeyLock(pending)
				pending = [][]string{}
			}
			continue
		}

		pending = append(pending, strings.Split(line, ""))
	}

	if len(pending) > 0 {
		parseKeyLock(pending)
		pending = [][]string{}
	}

	return locks, keys
}

func canFit(lock, key KeyLock) bool {
	if lock.height != key.height || len(lock.pins) != len(key.pins) {
		return false
	}
	for i := 0; i < len(lock.pins); i++ {
		if lock.pins[i]+key.pins[i] > lock.height {
			return false
		}
	}
	return true
}

func main() {
	locks, keys := parseInput()

	pairs := 0
	for _, lock := range locks {
		for _, key := range keys {
			if canFit(lock, key) {
				pairs += 1
			}
		}
	}
	fmt.Println("pairs:", pairs)
}

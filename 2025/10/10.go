package main

import (
	"aoc2025-10/flags"
	"aoc2025-10/machine"
	"aoc2025-10/part1"
	"aoc2025-10/part2"
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

func parseInput() []machine.Machine {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	machines := []machine.Machine{}

	for scanner.Scan() {
		line := scanner.Text()
		m, err := machine.ParseMachine(line)
		if err != nil {
			fmt.Println(fmt.Errorf("invalid input: %s", err.Error()))
			os.Exit(1)
		}
		machines = append(machines, m)
	}

	return machines
}

func main() {
	flag.Parse()
	showTime := *flags.Time

	start := time.Now()
	machines := parseInput()

	// Part 1
	p1Start := time.Now()
	totalPressesToStart := part1.GetTotalBtnPressesToStartMachines(machines)
	fmt.Printf("Part 1: %d", totalPressesToStart)
	if showTime {
		fmt.Printf(" -- %v\n", time.Since(p1Start))
	} else {
		fmt.Printf("\n")
	}

	// Part 2
	p2Start := time.Now()
	totalPressesForJoltage := part2.GetMinBtnPressesToConfigureJoltage(machines)
	fmt.Printf("Part 2: %d", totalPressesForJoltage)
	if showTime {
		fmt.Printf(" -- %v\n", time.Since(p2Start))
	} else {
		fmt.Printf("\n")
	}

	if showTime {
		fmt.Println("total time:", time.Since(start))
	}
}

package main

import (
	"flag"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"
)

func main() {
	flagInputPath := flag.String("i", "", "path to input file")
	flagTimes := flag.Bool("t", false, "show times")
	flagDebug := flag.Bool("d", false, "print additional debug info")
	flag.Parse()

	filePath := *flagInputPath
	if filePath == "" {
		fmt.Println("input file path is required")
		os.Exit(1)
	}

	inp_file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	start := time.Now()
	input := ParseInput(inp_file)

	s := Solver{
		withTimes:     *flagTimes,
		withDebugInfo: *flagDebug,
		input:         input,
	}

	if s.withTimes {
		defer func() {
			fmt.Printf("\ntotal time: %v\n", time.Since(start))
		}()
	}

	fmt.Printf("--- part 1 ---\n")
	s.part1()

	fmt.Printf("\n--- part 2 ---\n")
	s.part2()
}

type Solver struct {
	withTimes     bool
	withDebugInfo bool
	input         ParsedInput
}

func (s *Solver) part1() {
	start := time.Now()
	defer func() {
		if s.withTimes {
			fmt.Printf("time: %v\n", time.Since(start))
		}
	}()

	ev, err := NewEvaluator(s.input.gates)
	if err != nil {
		fmt.Printf("FAIL: %v\n", err)
		return
	}

	res := ev.Evaluate(s.input.state)

	fmt.Printf("z: %v\n", res.GetOutput())

}

func (s *Solver) part2() {
	start := time.Now()
	defer func() {
		if s.withTimes {
			fmt.Printf("time: %v\n", time.Since(start))
		}
	}()

	swaps := FindSwaps(s.input.gates, s.input.maxZ, s.withDebugInfo)

	if !VerifySwaps(s.input.gates, s.input.maxZ, swaps) {
		fmt.Printf("FAIL: failed to find swaps that fix the circuit\n")
		return
	}

	swappedGates := []string{}

	for _, swap := range swaps {
		swappedGates = append(swappedGates, swap[0], swap[1])
	}

	slices.Sort(swappedGates)

	fmt.Printf("swapped: %v\n", strings.Join(swappedGates, ","))
}

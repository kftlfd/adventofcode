package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParsedInput struct {
	state GatesState
	gates []Gate
	maxZ  int
}

func ParseInput(inp io.Reader) (input ParsedInput) {
	state := GatesState{}
	gates := []Gate{}

	scanner := bufio.NewScanner(inp)
	scanner.Split(bufio.ScanLines)

	mode := "STATE"

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case line == "":
			mode = "GATES"

		case mode == "STATE":
			gate, val, parseErr := parseStateLine(line)
			if parseErr != nil {
				panic(fmt.Errorf("bad input: %w", parseErr))
			}
			state[gate] = val > 0

		case mode == "GATES":
			gate, parseErr := parseGateLine(line)
			if parseErr != nil {
				panic(fmt.Errorf("bad input: %w", parseErr))
			}
			gates = append(gates, gate)
		}
	}

	err := scanner.Err()
	if err != nil {
		panic(err)
	}

	maxZ := getMaxZ(gates)

	return ParsedInput{state: state, gates: gates, maxZ: maxZ}
}

func parseStateLine(line string) (gate string, val int, err error) {
	vals := strings.Split(line, ": ")
	if len(vals) != 2 {
		err = fmt.Errorf("invalid state line: %v", line)
		return
	}

	val, err = strconv.Atoi(vals[1])
	if err != nil {
		err = fmt.Errorf("invlid gate state: %w", err)
		return
	}

	gate = vals[0]
	return
}

func parseGateLine(line string) (gate Gate, err error) {
	vals := strings.Split(line, " ")
	if len(vals) != 5 {
		err = fmt.Errorf("invalid gate line: %v", line)
		return
	}

	var op GateOp
	switch vals[1] {
	case "AND":
		op = AND
	case "OR":
		op = OR
	case "XOR":
		op = XOR
	default:
		err = fmt.Errorf("invlid op: %v", vals[1])
		return
	}

	gate = Gate{
		in1: vals[0],
		op:  op,
		in2: vals[2],
		out: vals[4],
	}
	return
}

func getMaxZ(gates []Gate) int {
	max := 0
	for _, gate := range gates {
		out := gate.out
		if out[0] != 'z' {
			continue
		}
		n, err := strconv.Atoi(out[1:])
		if err != nil {
			panic(fmt.Errorf("invalid z wire: %v", gate.out))
		}
		if n > max {
			max = n
		}
	}
	return max
}

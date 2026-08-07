package main

import (
	"fmt"
	"strconv"
)

type Gate struct {
	in1 string
	op  GateOp
	in2 string
	out string
}

type GateOp string

const (
	AND GateOp = "AND"
	OR  GateOp = "OR"
	XOR GateOp = "XOR"
)

type GatesState map[string]bool

func (s GatesState) Reset() {
	clear(s)
}

func (s GatesState) SetInputs(x, y, maxZ int) error {
	if x+y >= (1 << (maxZ + 1)) {
		return fmt.Errorf("result won't be fully evaluated")
	}

	s.Reset()

	g := func(pref string, bit int) string {
		return fmt.Sprintf("%v%02d", pref, bit)
	}

	for bit := 0; bit < maxZ; bit++ {
		mask := 1 << bit
		if x&mask == mask {
			s[g("x", bit)] = true
		}
		if y&mask == mask {
			s[g("y", bit)] = true
		}
	}

	return nil
}

func (s GatesState) GetOutput() int {
	result := 0

	for g, v := range s {
		if v && g[0] == 'z' {
			bit, err := strconv.Atoi(g[1:])
			if err != nil {
				panic(fmt.Errorf("invalid z gate: %v", g))
			}
			result |= 1 << bit
		}
	}

	return result
}

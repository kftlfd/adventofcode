package main

import (
	"fmt"
	"slices"
)

type WireRole int

const (
	WX WireRole = iota
	WY
	WPartialSum
	WPartialCarry
	WCarryPropagate
	WZ
	WCarryOut
)

type Identity struct {
	Bit  int
	Role WireRole
}

func id(bit int, role WireRole) Identity {
	return Identity{Bit: bit, Role: role}
}

type Subcircuit [7]string

type Circuit []Subcircuit

type Fixer struct {
	gates   []Gate
	maxZ    int
	circuit Circuit
	swaps   [][2]string
	issues  []string
}

func (f *Fixer) resetCircuit() {
	f.circuit = make(Circuit, f.maxZ)
	f.issues = []string{}
}

func (f *Fixer) addIssue(bit int, msg string) {
	f.issues = append(f.issues, fmt.Sprintf("%v: %v", bit, msg))
}

func (f *Fixer) findGateByInputs(in1 string, op GateOp, in2 string) *Gate {
	for i, gate := range f.gates {
		if gate.op == op && ((gate.in1 == in1 && gate.in2 == in2) ||
			(gate.in1 == in2 && gate.in2 == in1)) {
			return &f.gates[i]
		}
	}
	return nil
}

func (f *Fixer) findAdjacentWiresByGateOp(wire string, op GateOp) []string {
	res := []string{}
	for _, gate := range f.gates {
		if gate.op == op {
			switch {
			case gate.in1 == wire:
				res = append(res, gate.in2)
			case gate.in2 == wire:
				res = append(res, gate.in1)
			}
		}
	}
	return res
}

func (f *Fixer) swapOutputs(a, b string) {
	f.swaps = append(f.swaps, [2]string{a, b})
	for i := range f.gates {
		switch f.gates[i].out {
		case a:
			f.gates[i].out = b
		case b:
			f.gates[i].out = a
		}
	}
}

func (f *Fixer) assign(wire string, id Identity) (newChanges bool) {
	cur := f.circuit[id.Bit][id.Role]
	if cur != "" && cur != wire {
		f.swapOutputs(f.circuit[id.Bit][id.Role], wire)
		return true
	}
	f.circuit[id.Bit][id.Role] = wire
	return false
}

func (f *Fixer) checkBit0() (newChanges bool) {
	bit := 0

	x := fmt.Sprintf("x%02d", bit)
	y := fmt.Sprintf("y%02d", bit)
	z := fmt.Sprintf("z%02d", bit)

	if f.assign(x, id(bit, WX)) {
		return true
	}
	if f.assign(y, id(bit, WY)) {
		return true
	}
	if f.assign(z, id(bit, WZ)) {
		return true
	}

	finalSum := f.findGateByInputs(x, XOR, y)
	if finalSum != nil {
		if f.assign(finalSum.out, id(bit, WZ)) {
			return true
		}
	} else {
		f.addIssue(bit, "no final_sum")
	}

	carryOut := f.findGateByInputs(x, AND, y)
	if carryOut != nil {
		if f.assign(carryOut.out, id(bit, WCarryOut)) {
			return true
		}
	} else {
		f.addIssue(bit, "no carry_out")
	}

	return false
}

func (f *Fixer) checkBitN(bit int) (newChanges bool) {
	if bit < 1 || bit >= f.maxZ {
		panic("invalid bit")
	}

	x := fmt.Sprintf("x%02d", bit)
	y := fmt.Sprintf("y%02d", bit)
	z := fmt.Sprintf("z%02d", bit)

	if f.assign(x, id(bit, WX)) {
		return true
	}
	if f.assign(y, id(bit, WY)) {
		return true
	}
	if f.assign(z, id(bit, WZ)) {
		return true
	}

	if bit == f.maxZ-1 {
		lastCarry := fmt.Sprintf("z%02d", f.maxZ)
		if f.assign(lastCarry, id(bit, WCarryOut)) {
			return true
		}
	}

	partialSum := f.findGateByInputs(x, XOR, y)
	if partialSum != nil {
		if f.assign(partialSum.out, id(bit, WPartialSum)) {
			return true
		}
	} else {
		f.addIssue(bit, "no partial_sum")
	}

	partialCarry := f.findGateByInputs(x, AND, y)
	if partialCarry != nil {
		if f.assign(partialCarry.out, id(bit, WPartialCarry)) {
			return true
		}
	} else {
		f.addIssue(bit, "no partial_carry")
	}

	pSum := f.circuit[bit][WPartialSum]
	carryIn := f.circuit[bit-1][WCarryOut]

	foundFinalOut := false

	if pSum != "" && carryIn != "" {
		finalOut := f.findGateByInputs(pSum, XOR, carryIn)
		if finalOut != nil {
			foundFinalOut = true
			if f.assign(finalOut.out, id(bit, WZ)) {
				return true
			}
		}
	}

	if !foundFinalOut {
		f.addIssue(bit, "no final_out")

		if pSum != "" {
			for _, wire := range f.findAdjacentWiresByGateOp(pSum, XOR) {
				if f.assign(wire, id(bit-1, WCarryOut)) {
					return true
				}
			}
		}

		if carryIn != "" {
			for _, wire := range f.findAdjacentWiresByGateOp(carryIn, XOR) {
				if f.assign(wire, id(bit, WPartialSum)) {
					return true
				}
			}
		}
	}

	foundCarryProp := false

	if pSum != "" && carryIn != "" {
		carryProp := f.findGateByInputs(pSum, AND, carryIn)
		if carryProp != nil {
			foundCarryProp = true
			if f.assign(carryProp.out, id(bit, WCarryPropagate)) {
				return true
			}
		}
	}

	if !foundCarryProp {
		f.addIssue(bit, "no carry_propagate")

		if pSum != "" {
			for _, wire := range f.findAdjacentWiresByGateOp(pSum, AND) {
				if f.assign(wire, id(bit-1, WCarryOut)) {
					return true
				}
			}
		}

		if carryIn != "" {
			for _, wire := range f.findAdjacentWiresByGateOp(carryIn, AND) {
				if f.assign(wire, id(bit, WPartialSum)) {
					return true
				}
			}
		}
	}

	carryProp := f.circuit[bit][WCarryPropagate]
	pCarry := f.circuit[bit][WPartialCarry]

	foundCarryOut := false

	if carryProp != "" && pCarry != "" {
		carryOut := f.findGateByInputs(carryProp, OR, pCarry)
		if carryOut != nil {
			foundCarryOut = true
			if f.assign(carryOut.out, id(bit, WCarryOut)) {
				return true
			}
		}
	}

	if !foundCarryOut {
		f.addIssue(bit, "no carry_out")

		if carryProp != "" {
			f.addIssue(bit, "empty partial_carry")
			for _, wire := range f.findAdjacentWiresByGateOp(carryProp, OR) {
				if f.assign(wire, id(bit, WPartialCarry)) {
					return true
				}
			}
		}

		if pCarry != "" {
			f.addIssue(bit, "empty carry_propagate")
			for _, wire := range f.findAdjacentWiresByGateOp(pCarry, OR) {
				if f.assign(wire, id(bit, WCarryPropagate)) {
					return true
				}
			}
		}
	}

	return false
}

func FindSwaps(gates []Gate, maxZ int, printDebugInfo bool) [][2]string {
	f := &Fixer{gates: slices.Clone(gates), maxZ: maxZ}

	changes := true
	for changes {
		f.resetCircuit()

		changes = f.checkBit0()
		if changes {
			continue
		}
		for bit := 1; bit < f.maxZ; bit++ {
			changes = f.checkBitN(bit)
			if changes {
				break
			}
		}
	}

	if printDebugInfo {
		fmt.Printf("\n-- issues --\n")
		for _, issue := range f.issues {
			fmt.Println(issue)
		}

		fmt.Printf("\n-- swaps --\n")
		for _, swap := range f.swaps {
			fmt.Printf("%v\n", swap)
		}

		fmt.Printf("\n-- circuit --\n")
		for _, sc := range f.circuit {
			for _, wire := range sc {
				fmt.Printf("%-5s", wire)
			}
			fmt.Printf("\n")
		}
	}

	return f.swaps
}

func VerifySwaps(gates []Gate, maxZ int, swaps [][2]string) bool {
	g := slices.Clone(gates)

	for _, swap := range swaps {
		a, b := swap[0], swap[1]

		for i := range g {
			switch g[i].out {
			case a:
				g[i].out = b
			case b:
				g[i].out = a
			}
		}
	}

	f := &Fixer{gates: g, maxZ: maxZ}

	f.resetCircuit()
	f.checkBit0()
	for bit := 1; bit < f.maxZ; bit++ {
		f.checkBitN(bit)
	}
	if len(f.issues) > 0 {
		return false
	}

	return true
}

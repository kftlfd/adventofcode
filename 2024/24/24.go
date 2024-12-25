package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

const (
	AND = "AND"
	XOR = "XOR"
	OR  = "OR"
)

func getOp(val string) string {
	switch val {
	case AND:
		return AND
	case XOR:
		return XOR
	case OR:
		return OR
	}
	panic(fmt.Errorf("bad op: %v", val))
}

func parseInput() (state map[string]int, instructions [][4]string) {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	scanner := bufio.NewScanner(inp_file)
	scanner.Split(bufio.ScanLines)

	state_map := map[string]int{}
	instructions = [][4]string{}

	is_instructions := false

	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			is_instructions = true
			continue
		}

		if is_instructions {
			vals := strings.Split(line, " ")
			if len(vals) != 5 {
				panic(fmt.Errorf("bad input: %v", line))
			}

			op := getOp(vals[1])
			instructions = append(instructions, [4]string{op, vals[0], vals[2], vals[4]})
			continue
		}

		vals := strings.Split(line, ": ")
		if len(vals) != 2 {
			panic(fmt.Errorf("bad input: %v", line))
		}

		val, err := strconv.Atoi(vals[1])
		if err != nil {
			panic(err)
		}
		state_map[vals[0]] = val
	}

	return state_map, instructions
}

type InstrCycle struct {
	state           *map[string]int
	intructions_ran []int
}

func applyOp(op string, v1, v2 int) int {
	switch op {
	case XOR:
		return v1 ^ v2
	case AND:
		return v1 & v2
	case OR:
		return v1 | v2
	}
	panic(fmt.Errorf("unknown instruction: %v", op))
}

func runInstructions(instructions [][4]string, state map[string]int) ([]InstrCycle, error) {
	done := make([]bool, len(instructions))
	all_done := false

	cycles := []InstrCycle{}

	for !all_done {
		all_done = true
		changes := false
		ran := []int{}
		for i, instr := range instructions {
			if done[i] {
				continue
			}

			v1, ok1 := state[instr[1]]
			v2, ok2 := state[instr[2]]

			if !ok1 || !ok2 {
				all_done = false
				continue
			}
			changes = true

			out := instr[3]
			val := applyOp(instr[0], v1, v2)

			state[out] = val
			done[i] = true
			ran = append(ran, i)
		}

		if !changes {
			return cycles, fmt.Errorf("instructions cannot be completed")
		}

		cycle_state := map[string]int{}
		for k, v := range state {
			cycle_state[k] = v
		}
		cur_cycle := InstrCycle{state: &cycle_state, intructions_ran: ran}
		cycles = append(cycles, cur_cycle)
	}

	return cycles, nil
}

func getStateNumber(prefix string, state map[string]int) int {
	num := 0

	for k, v := range state {
		if k[:1] != prefix || v == 0 {
			continue
		}

		pos, err := strconv.Atoi(k[1:])
		if err != nil {
			panic(err)
		}

		num += 1 << pos
	}

	return num
}

func Part1() {
	state, instructions := parseInput()

	_, err := runInstructions(instructions, state)
	if err != nil {
		panic(err)
	}

	z := getStateNumber("z", state)
	fmt.Println("z:", z)
}

//
//
//

type Machine struct {
	starting_state       map[string]int
	instructions         [][4]string
	x, y, correct_result int
}

func copyState(s map[string]int) map[string]int {
	copy := map[string]int{}
	for k, v := range s {
		copy[k] = v
	}
	return copy
}

func isPositionSet(num int, pos int) bool {
	mask := 1 << pos
	return num&mask == mask
}

func bitAtPosition(num int, pos int) int {
	if isPositionSet(num, pos) {
		return 1
	}
	return 0
}

func withCorrections(instructions [][4]string, corrections [][2]string) [][4]string {
	corrections_map := map[string]string{}
	for _, corr := range corrections {
		corrections_map[corr[0]] = corr[1]
		corrections_map[corr[1]] = corr[0]
	}

	corrected := [][4]string{}

	for _, instr := range instructions {
		out := instr[3]
		swapped, swap_out := corrections_map[out]
		if swap_out {
			out = swapped
		}
		corrected = append(corrected, [4]string{instr[0], instr[1], instr[2], out})
	}

	return corrected
}

type Wrong struct {
	op        string
	names     []string
	vals      []int
	should_be int
}

func (m *Machine) tryCorrections(corrections [][2]string) (ok bool, wrong []Wrong) {
	state := copyState(m.starting_state)
	instructions := m.instructions
	z := m.correct_result

	_, err := runInstructions(withCorrections(instructions, corrections), state)
	if err != nil {
		return false, []Wrong{}
	}

	res := getStateNumber("z", state)

	for _, instr := range instructions {
		if instr[3][:1] != "z" {
			continue
		}

		pos, err := strconv.Atoi(instr[3][1:])
		if err != nil {
			panic(err)
		}

		should_be := bitAtPosition(z, pos)
		have := bitAtPosition(res, pos)

		if have != should_be {
			wrong = append(wrong, Wrong{
				op:        instr[0],
				names:     []string{instr[1], instr[2], instr[3]},
				vals:      []int{state[instr[1]], state[instr[2]], state[instr[3]]},
				should_be: should_be,
			})
		}
	}

	return true, wrong
}

func nameInCorrections(corrections [][2]string, name string) bool {
	for _, row := range corrections {
		if row[0] == name || row[1] == name {
			return true
		}
	}
	return false
}

var check_fails = 0

func (m *Machine) findCorrectionsDfs(corrections [][2]string) (bool, [][2]string) {
	finishes, wrong := m.tryCorrections(corrections)
	if finishes && len(wrong) < 1 {
		return true, corrections
	}
	if !finishes || len(corrections) >= 4 {
		check_fails += 1
		if check_fails%200 == 0 {
			fmt.Printf("fails: %v\n", check_fails)
		}
		return false, corrections
	}

	for i, w := range wrong {
		// swap outputs
		out1name := w.names[2]
		out1v := w.vals[2]
		if !nameInCorrections(corrections, out1name) {
			for j := i + 1; j < len(wrong); j++ {
				nxt_w := wrong[j]
				out2name := nxt_w.names[2]
				out2v := nxt_w.vals[2]

				if nameInCorrections(corrections, out2name) {
					continue
				}
				if w.should_be != out2v || nxt_w.should_be != out1v {
					continue
				}

				upd_corr := append(corrections, [2]string{out1name, out2name})
				ok, ww := m.findCorrectionsDfs(upd_corr)
				if ok && len(ww) < 1 {
					return true, upd_corr
				}
			}
		}

		getSwapCandidates := func(target_val int) []string {
			candidates := []string{}
			for j := i + 1; j < len(wrong); j++ {
				ww := wrong[j]

				if ww.vals[1] == target_val {
					if applyOp(ww.op, target_val^1, ww.vals[2]) == ww.should_be && !nameInCorrections(corrections, ww.names[1]) {
						candidates = append(candidates, ww.names[1])
					}
				}

				if ww.vals[2] == target_val {
					if applyOp(ww.op, target_val^1, ww.vals[2]) == ww.should_be && !nameInCorrections(corrections, ww.names[2]) {
						candidates = append(candidates, ww.names[2])
					}
				}
			}
			return candidates
		}

		// swap inputs: XOR
		if w.op == XOR {
			for _, cand := range getSwapCandidates(w.vals[1] ^ 1) {
				upd_corr := append(corrections, [2]string{w.names[1], cand})
				ok, ww := m.findCorrectionsDfs(upd_corr)
				if ok && len(ww) < 1 {
					return true, upd_corr
				}
			}
			for _, cand := range getSwapCandidates(w.vals[2] ^ 1) {
				upd_corr := append(corrections, [2]string{w.names[2], cand})
				ok, ww := m.findCorrectionsDfs(upd_corr)
				if ok && len(ww) < 1 {
					return true, upd_corr
				}
			}
		}

		// swap inputs: AND, OR
		if w.vals[1] == 0 {
			for _, cand := range getSwapCandidates(1) {
				upd_corr := append(corrections, [2]string{w.names[1], cand})
				ok, ww := m.findCorrectionsDfs(upd_corr)
				if ok && len(ww) < 1 {
					return true, upd_corr
				}
			}
		}
		if w.vals[2] == 0 {
			for _, cand := range getSwapCandidates(1) {
				upd_corr := append(corrections, [2]string{w.names[2], cand})
				ok, ww := m.findCorrectionsDfs(upd_corr)
				if ok && len(ww) < 1 {
					return true, upd_corr
				}
			}
		}
	}

	return false, corrections
}

func (m *Machine) findCorrections() ([][2]string, error) {
	ok, corrections := m.findCorrectionsDfs([][2]string{})
	if ok {
		return corrections, nil
	}
	return [][2]string{}, fmt.Errorf("corrections not found")
}

type AlphArr []string

func (a AlphArr) Len() int           { return len(a) }
func (a AlphArr) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a AlphArr) Less(i, j int) bool { return a[i] < a[j] }
func formatCorrectionsAns(corrections [][2]string) string {
	arr := []string{}
	for _, pair := range corrections {
		arr = append(arr, pair[0], pair[1])
	}
	sort.Sort(AlphArr(arr))
	return strings.Join(arr, ",")
}

func part2() {
	state, instructions := parseInput()
	x := getStateNumber("x", state)
	y := getStateNumber("y", state)
	z := x + y

	machine := Machine{starting_state: state, instructions: instructions, x: x, y: y, correct_result: z}

	ans, err := machine.findCorrections()
	if err != nil {
		fmt.Println("corrections not found")
	} else {
		fmt.Printf("corrections: %v\n", formatCorrectionsAns(ans))
	}
}

//
//
//

func main() {
	fmt.Printf("--- part 1 ---\n")
	Part1()

	// fmt.Printf("\n--- part 2 ---\n")
	// part2()
}

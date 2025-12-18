package machine

import (
	"aoc2025-10/matrix"
	"fmt"
	"strconv"
	"strings"
)

type Machine struct {
	LightDiagram []bool
	Buttons      [][]int
	Joltage      []int
	str          string
}

func (m Machine) GetString() string { return m.str }

func ParseMachine(line string) (Machine, error) {
	m := Machine{str: line}

	fields := strings.Fields(line)
	if len(fields) < 3 {
		return m, fmt.Errorf("fields < 3")
	}

	diagram := []bool{}
	buttons := [][]int{}
	joltage := []int{}

	lightsStr := fields[0]
	if len(lightsStr) < 3 {
		return m, fmt.Errorf("lightDiagram len < 3")
	}
	if lightsStr[0] != '[' || lightsStr[len(lightsStr)-1] != ']' {
		return m, fmt.Errorf("invalid light diagram format")
	}
	for i := 1; i < len(lightsStr)-1; i++ {
		cur := lightsStr[i : i+1]
		switch cur {
		case ".":
			diagram = append(diagram, false)
		case "#":
			diagram = append(diagram, true)
		default:
			return m, fmt.Errorf("invalid light diagram format")
		}
	}

	buttonsStr := fields[1 : len(fields)-1]
	if len(buttonsStr) < 1 {
		return m, fmt.Errorf("no buttons")
	}
	for _, btnStr := range buttonsStr {
		if btnStr[0] != '(' || btnStr[len(btnStr)-1] != ')' {
			return m, fmt.Errorf("invalid button format")
		}
		numsStrArr := strings.Split(btnStr[1:len(btnStr)-1], ",")
		nums := []int{}
		for _, numStr := range numsStrArr {
			num, err := strconv.Atoi(numStr)
			if err != nil {
				return m, fmt.Errorf("invali button format, expected ints")
			}
			nums = append(nums, num)
		}
		buttons = append(buttons, nums)
	}

	joltageStr := fields[len(fields)-1]
	if len(joltageStr) < 3 {
		return m, fmt.Errorf("invalid jotage string")
	}
	if joltageStr[0] != '{' || joltageStr[len(joltageStr)-1] != '}' {
		return m, fmt.Errorf("invalid joltage format")
	}
	numsStrArr := strings.SplitSeq(joltageStr[1:len(joltageStr)-1], ",")
	for numStr := range numsStrArr {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return m, fmt.Errorf("invali joltage format, expected ints")
		}
		joltage = append(joltage, num)
	}

	m.LightDiagram = diagram
	m.Buttons = buttons
	m.Joltage = joltage
	return m, nil
}

func (m Machine) Print() {
	fmt.Printf("%+v\n", m)
}

func (m Machine) ToBtnJoltageAugmentedMatrix() matrix.AugmMatrix2d {
	mat := make(matrix.AugmMatrix2d, len(m.Joltage))
	n := len(m.Buttons) + 1
	for j, joltage := range m.Joltage {
		mat[j] = make(matrix.Matrix1d, n)
		mat[j][n-1] = joltage
	}
	for j, btn := range m.Buttons {
		for _, pos := range btn {
			mat[pos][j] = 1
		}
	}
	return mat
}

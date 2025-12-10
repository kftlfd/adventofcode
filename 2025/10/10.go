package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Machine struct {
	lightDiagram []bool
	buttons      [][]int
	joltage      []int
	masks        MachineMasks
}

func parseInput() []Machine {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	machines := []Machine{}

	for scanner.Scan() {
		line := scanner.Text()
		m, err := parseMachine(line)
		if err != nil {
			fmt.Println(fmt.Errorf("invalid input: %s", err.Error()))
			os.Exit(1)
		}
		machines = append(machines, m)
	}

	return machines
}

func parseMachine(line string) (Machine, error) {
	m := Machine{}

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
	if lightsStr[:1] != "[" || lightsStr[len(lightsStr)-1:] != "]" {
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
		if btnStr[:1] != "(" || btnStr[len(btnStr)-1:] != ")" {
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
	if joltageStr[:1] != "{" || joltageStr[len(joltageStr)-1:] != "}" {
		return m, fmt.Errorf("invalid joltage format")
	}
	numsStrArr := strings.Split(joltageStr[1:len(joltageStr)-1], ",")
	for _, numStr := range numsStrArr {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return m, fmt.Errorf("invali joltage format, expected ints")
		}
		joltage = append(joltage, num)
	}

	m.lightDiagram = diagram
	m.buttons = buttons
	m.joltage = joltage
	m.masks = getMasks(diagram, buttons)
	return m, nil
}

type MachineMasks struct {
	diagram int
	buttons []int
}

func boolArrToMask(arr []bool) int {
	mask := 0
	for i := len(arr) - 1; i >= 0; i-- {
		mask = (mask << 1)
		if arr[i] {
			mask += 1
		}
	}
	return mask
}

func buttonToMask(n int, btn []int) int {
	arr := make([]bool, n)
	for _, pos := range btn {
		arr[pos] = true
	}
	return boolArrToMask(arr)
}

func getMasks(lightDiagram []bool, buttons [][]int) MachineMasks {
	n := len(lightDiagram)
	diagram := boolArrToMask(lightDiagram)
	btns := []int{}
	for _, btn := range buttons {
		btns = append(btns, buttonToMask(n, btn))
	}
	return MachineMasks{diagram: diagram, buttons: btns}
}

func (m Machine) getFewestPressesToStart() int {
	results := []int{0}
	nxtResults := []int{}
	curLen := 1

	for true {
		for _, state := range results {
			for _, btnMask := range m.masks.buttons {
				nxtState := state ^ btnMask
				if nxtState == m.masks.diagram {
					return curLen
				}
				nxtResults = append(nxtResults, nxtState)
			}
		}
		results, nxtResults = nxtResults, []int{}
		curLen += 1
	}

	return -1
}

func main() {
	start := time.Now()
	machines := parseInput()

	// Part 1
	totalPressesToStart := 0
	for _, m := range machines {
		presses := m.getFewestPressesToStart()
		totalPressesToStart += presses
	}
	fmt.Println("Part 1:", totalPressesToStart, time.Since(start))
}

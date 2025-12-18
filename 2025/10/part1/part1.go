package part1

import "aoc2025-10/machine"

type machineMasks struct {
	diagram int
	buttons []int
}

func getMasks(m machine.Machine) machineMasks {
	lightDiagram := m.LightDiagram
	buttons := m.Buttons

	n := len(lightDiagram)
	diagram := boolArrToMask(lightDiagram)
	btns := []int{}
	for _, btn := range buttons {
		btns = append(btns, buttonToMask(n, btn))
	}
	return machineMasks{diagram: diagram, buttons: btns}
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

func tryBtnMask(buttons []int, buttonsMask int) (buttonPresses, result int) {
	for i := range len(buttons) {
		cur := 1 << i
		pressButton := cur&buttonsMask == cur
		if pressButton {
			result ^= buttons[i]
			buttonPresses += 1
		}
	}
	return buttonPresses, result
}

func getMinPresses(m machine.Machine) int {
	masks := getMasks(m)
	buttonsN := len(m.Buttons)
	minPresses := buttonsN
	for btnMask := range (1 << buttonsN) - 1 {
		presses, res := tryBtnMask(masks.buttons, btnMask)
		if res == masks.diagram {
			minPresses = min(minPresses, presses)
		}
	}
	return minPresses
}

func GetTotalBtnPressesToStartMachines(machines []machine.Machine) int {
	totalPressesToStart := 0
	for _, m := range machines {
		totalPressesToStart += getMinPresses(m)
	}
	return totalPressesToStart
}

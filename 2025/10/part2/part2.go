package part2

import (
	"aoc2025-10/flags"
	"aoc2025-10/machine"
)

func GetMinBtnPressesToConfigureJoltage(machines []machine.Machine) int {
	if *flags.P2SimpleDFS {
		return dfsGetMinBtnPressesToConfigureJoltage(machines)
	}

	return matrixSolverGetMinBtnPressesToConfigureJoltage(machines)
}

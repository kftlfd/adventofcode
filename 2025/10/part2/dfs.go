package part2

import (
	"aoc2025-10/flags"
	"aoc2025-10/machine"
	"fmt"
	"math"
	"sync"
)

type dfsMemo map[string]bool

func getDfsMemoKey(state []int, presses int) string {
	return fmt.Sprintf("%d/%+v", presses, state)
}

func dfs(memo dfsMemo, state []int, presses int, joltage []int, buttons [][]int, minPresses *int) {
	if presses >= *minPresses {
		return
	}

	memoKey := getDfsMemoKey(state, presses)
	if inMemo := memo[memoKey]; inMemo {
		return
	}
	memo[memoKey] = true

	for _, btn := range buttons {
		nxtPresses := presses + 1
		nxtState := append([]int{}, state...)

		overshoot := false
		for pos, v := range btn {
			nxtState[pos] += v
			if nxtState[pos] > joltage[pos] {
				overshoot = true
				break
			}
		}
		if overshoot {
			continue
		}

		foundTarget := true
		for pos, v := range nxtState {
			if joltage[pos] != v {
				foundTarget = false
				break
			}
		}
		if foundTarget {
			*minPresses = min(*minPresses, nxtPresses)
			return
		}

		dfs(memo, nxtState, nxtPresses, joltage, buttons, minPresses)
	}
}

func getMinPressesToConfigureJoltage(m machine.Machine) int {
	joltage := m.Joltage
	buttons := make([][]int, len(m.Buttons))

	for i, btn := range m.Buttons {
		cur := make([]int, len(joltage))
		for _, pos := range btn {
			cur[pos] = 1
		}
		buttons[i] = cur
	}

	memo := dfsMemo{}
	minPresses := math.MaxInt
	initState := make([]int, len(joltage))

	dfs(memo, initState, 0, joltage, buttons, &minPresses)

	return minPresses
}

func dfsGetMinBtnPressesToConfigureJoltage(machines []machine.Machine) int {
	debug := *flags.Debug

	type result struct {
		machineIdx int
		presses    int
	}

	in := make(chan int)
	out := make(chan result)

	var wg sync.WaitGroup
	var wgRes sync.WaitGroup

	for range len(machines) {
		wg.Go(func() {
			for machineIdx := range in {
				presses := getMinPressesToConfigureJoltage(machines[machineIdx])
				out <- result{machineIdx: machineIdx, presses: presses}
			}
		})
	}

	wgRes.Go(func() {
		for i := range machines {
			in <- i
		}
		close(in)
		wg.Wait()
		close(out)
	})

	totalPressesForJoltage := 0
	cur := 0
	for res := range out {
		cur++
		if debug {
			fmt.Printf("%d %d/%d: %d\n", cur, res.machineIdx+1, len(machines), res.presses)
		}
		totalPressesForJoltage += res.presses
	}

	wgRes.Wait()

	return totalPressesForJoltage
}

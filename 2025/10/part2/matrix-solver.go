package part2

import (
	"aoc2025-10/flags"
	"aoc2025-10/machine"
	"fmt"
	"sync"
)

func matrixSolverGetMinBtnPressesToConfigureJoltage(machines []machine.Machine) int {
	debug := *flags.Debug

	totalPresses := 0
	solved := 0
	failed := 0

	type result struct {
		i, val int
		err    error
	}

	res := make(chan result)

	var wg sync.WaitGroup

	for i, m := range machines {
		wg.Go(func() {
			mat := m.ToBtnJoltageAugmentedMatrix().ToGaussJordanReduction()
			presses, err := mat.GetMinPositiveSolutionSum()
			res <- result{i: i, val: presses, err: err}
		})
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	resI := 1
	for r := range res {
		presses, err := r.val, r.err

		printDebug := func(err error) {
			m := machines[r.i]
			mat := m.ToBtnJoltageAugmentedMatrix()
			redMat := mat.ToGaussJordanReduction()
			solver := redMat.GetSolver()
			fmt.Println("")
			fmt.Println("Error:", err.Error())
			m.Print()
			fmt.Println("")
			fmt.Println("--mat--")
			mat.Print()
			fmt.Println("")
			fmt.Println("--reduced--")
			redMat.Print()
			fmt.Println("")
			fmt.Println("--constraints--")
			fmt.Println(solver.GetFreeVarsCount(), "/", solver.GetConstraints())
			fmt.Println("")
			fmt.Println("---")
			fmt.Println("")
		}

		if debug {
			freeVarsCount := machines[r.i].ToBtnJoltageAugmentedMatrix().ToGaussJordanReduction().GetSolver().GetFreeVarsCount()
			fmt.Printf("%d %d/%d: (free=%d) -> ", resI, r.i, len(machines), freeVarsCount)
			resI++
		}
		if err != nil {
			if debug {
				failed += 1
				fmt.Printf("-\n")
				printDebug(err)
			}
			return 0
		} else {
			if debug {
				solved += 1
				fmt.Printf("%v\n", presses)
			}
			totalPresses += presses
		}
	}

	if debug {
		fmt.Printf("total: %v   ok: %v   fails: %v\n", len(machines), solved, failed)
	}

	return totalPresses
}

package matrix

import (
	"fmt"
	"math"
)

type reducedAugmMatrix2dSolver struct {
	mat           ReducedAugmMatrix2d
	m, n          int
	constraints   []int
	leadingCols   []int
	freeVarsCount int
}

func (s reducedAugmMatrix2dSolver) GetFreeVarsCount() int { return s.freeVarsCount }
func (s reducedAugmMatrix2dSolver) GetConstraints() []int { return s.constraints }

func (s reducedAugmMatrix2dSolver) getMinPositiveSolutionSum() (int, error) {
	if s.freeVarsCount == 0 {
		sol := s.getSoluionStarter()
		err := s.fillPositiveConstrainedVars(sol)
		if err != nil {
			return 0, err
		}
		return getSliceSum(sol), nil
	}

	freeVarsVals := make([]int, s.freeVarsCount)
	memo := make(dfsMemo)
	// limit the search space -- optimizing for my problem input
	maxVal := 400 / s.freeVarsCount

	return dfs(s, memo, freeVarsVals, maxVal)
}

func (s reducedAugmMatrix2dSolver) getSoluionStarter() []int {
	return make([]int, s.n-1)
}

func (s reducedAugmMatrix2dSolver) fillFreeVars(solution []int, freeVarsVals []int) {
	idx := 0
	for i := 0; i < len(solution) && idx < len(freeVarsVals); i++ {
		if s.constraints[i] == -1 {
			solution[i] = freeVarsVals[idx]
			idx++
		}
	}
}

func (s reducedAugmMatrix2dSolver) fillPositiveConstrainedVars(solution []int) error {
	n := s.n - 1
	for i := range n {
		valRow := s.constraints[i]
		if valRow == -1 {
			// skip free vars
			continue
		}
		valCol := s.leadingCols[valRow]

		val := s.mat[valRow][valCol]
		if val == 0 {
			return fmt.Errorf("zero value")
		}

		// get init ans from the last col
		ans := s.mat[valRow][n]

		// substract free vars
		for col := valCol + 1; col < n; col++ {
			nxtVal := s.mat[valRow][col]
			ans -= nxtVal * solution[col]
		}

		if ans%val != 0 {
			return fmt.Errorf("not an int solution")
		}
		ans = ans / val
		if ans < 0 {
			return fmt.Errorf("non-positive solution")
		}

		solution[i] = ans
	}
	return nil
}

type memoVal struct {
	val int
	err error
}

type dfsMemo map[string]memoVal

func getMemoKey(freeVarsVal []int) string {
	return fmt.Sprintf("%+v", freeVarsVal)
}

func dfs(s reducedAugmMatrix2dSolver, memo dfsMemo, freeVarsVals []int, maxVal int) (int, error) {
	memoKey := getMemoKey(freeVarsVals)
	if val, ok := memo[memoKey]; ok {
		return val.val, val.err
	}

	sol := s.getSoluionStarter()
	s.fillFreeVars(sol, freeVarsVals)
	err := s.fillPositiveConstrainedVars(sol)

	curSolValid := true
	if err != nil {
		curSolValid = false
	}
	curSolSum := getSliceSum(sol)

	minSolValid := curSolValid
	minValidSolSum := math.MaxInt
	if curSolValid {
		minValidSolSum = curSolSum
	}

	for idx := range freeVarsVals {
		nxtFreeVarsVals := copySlice(freeVarsVals)
		nxtFreeVarsVals[idx] += 1
		if nxtFreeVarsVals[idx] > maxVal {
			continue
		}

		nxtSolSum, err := dfs(s, memo, nxtFreeVarsVals, maxVal)
		if err == nil {
			minValidSolSum = min(minValidSolSum, nxtSolSum)
			minSolValid = true
		}
	}

	if minSolValid {
		memo[memoKey] = memoVal{val: minValidSolSum, err: nil}
		return minValidSolSum, nil
	} else {
		err := fmt.Errorf("positive solution not found")
		memo[memoKey] = memoVal{val: 0, err: err}
		return 0, err
	}
}

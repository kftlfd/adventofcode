package matrix

import (
	"fmt"
)

type ReducedAugmMatrix2d []Matrix1d

func (mat ReducedAugmMatrix2d) Print() {
	for _, row := range mat {
		fmt.Println(row)
	}
}

func (mat ReducedAugmMatrix2d) GetSolver() reducedAugmMatrix2dSolver {
	m := len(mat)
	n := len(mat[0])

	constraints := make([]int, n-1)
	for col := range n - 1 {
		colValsCount := 0
		for row := range m {
			if mat[row][col] != 0 {
				constraints[col] = row
				colValsCount += 1
			}
		}
		if colValsCount > 1 {
			constraints[col] = -1
		}
	}

	leadingCols := make([]int, m)
	for row := range m {
		for col := range n - 1 {
			if mat[row][col] != 0 {
				leadingCols[row] = col
				break
			}
		}
	}

	freeVarsCount := 0
	for _, row := range constraints {
		if row == -1 {
			freeVarsCount += 1
		}
	}

	return reducedAugmMatrix2dSolver{
		mat: mat, m: m, n: n,
		constraints:   constraints,
		leadingCols:   leadingCols,
		freeVarsCount: freeVarsCount,
	}
}

func (mat ReducedAugmMatrix2d) GetMinPositiveSolutionSum() (int, error) {
	solver := mat.GetSolver()
	return solver.getMinPositiveSolutionSum()
}

package matrix

import (
	"fmt"
)

type Matrix1d []int

type Matrix2d []Matrix1d

type AugmMatrix2d []Matrix1d

func (mat AugmMatrix2d) Print() {
	for _, row := range mat {
		fmt.Println(row)
	}
}

func (mat AugmMatrix2d) copy() AugmMatrix2d {
	m := len(mat)
	n := len(mat[0])
	matCopy := make(AugmMatrix2d, m)
	for i := range m {
		matCopy[i] = make(Matrix1d, n)
		copy(matCopy[i], mat[i])
	}
	return matCopy
}

func (mat AugmMatrix2d) ToGaussJordanReduction() ReducedAugmMatrix2d {
	red := mat.copy()
	m := len(red)
	n := len(red[0])

	curRow := 0
	curCol := 0
	for curCol < n-1 && curRow < m {
		// find row with leading non-zero value
		swapRow := -1
		for row := curRow; row < m; row++ {
			if red[row][curCol] != 0 {
				swapRow = row
				break
			}
		}
		if swapRow == -1 {
			curCol++
			continue
		}

		// move row up
		if curRow != swapRow {
			red[curRow], red[swapRow] = red[swapRow], red[curRow]
		}

		// eliminate other values in current col
		base := red[curRow][curCol]
		for row := range m {
			factor := red[row][curCol]
			if row != curRow && factor != 0 {
				for col := range n {
					red[row][col] = red[row][col]*base - red[curRow][col]*factor
				}
			}
		}

		curRow++
		curCol++
	}

	// remove bottom all-zeroes rows
	lastRow := m - 1
	for lastRow > -1 {
		zeroes := 0
		for col := range n {
			if red[lastRow][col] == 0 {
				zeroes += 1
			}
		}
		if zeroes == n {
			lastRow -= 1
		} else {
			break
		}
	}

	return ReducedAugmMatrix2d(red[:lastRow+1])
}

package try2

import (
	"fmt"
	"math"
	"strconv"
)

func getNumFromCode(code string) int {
	num, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		panic(err)
	}
	return num
}

const (
	R = 0
	C = 1
)

const (
	UP    = 0
	DOWN  = 1
	LEFT  = 2
	RIGHT = 3
	A     = 4
)

func getNumericCoordinates(code string) [2]int {
	/*
		7 8 9
		4 5 6
		1 2 3
		  0 A
	*/
	switch code {
	case "7":
		return [2]int{0, 0}
	case "8":
		return [2]int{0, 1}
	case "9":
		return [2]int{0, 2}
	case "4":
		return [2]int{1, 0}
	case "5":
		return [2]int{1, 1}
	case "6":
		return [2]int{1, 2}
	case "1":
		return [2]int{2, 0}
	case "2":
		return [2]int{2, 1}
	case "3":
		return [2]int{2, 2}
	case "0":
		return [2]int{3, 1}
	case "A":
		return [2]int{3, 2}
	default:
		panic(fmt.Errorf("unrecognized numeric code: %v", code))
	}
}

func dfs(code string, goal_i, cur_r, cur_c int, path []int, res *[][]int) {
	if goal_i >= len(code) {
		(*res) = append((*res), path)
		return
	}

	coords := getNumericCoordinates(code[goal_i : goal_i+1])
	goal_r := coords[0]
	goal_c := coords[1]

	if cur_r == goal_r && cur_c == goal_c {
		path = append(path, A)
		dfs(code, goal_i+1, cur_r, cur_c, path, res)
		return
	}

	getMove := func(move, times int) []int {
		moves := []int{}
		for i := 0; i < times; i++ {
			moves = append(moves, move)
		}
		return moves
	}

	if cur_r < goal_r && !(cur_c == 0 && goal_r == 3) {
		dfs(code, goal_i, goal_r, cur_c, append(path, getMove(DOWN, goal_r-cur_r)...), res)
	}
	if cur_c > goal_c && !(cur_r == 3 && goal_c == 0) {
		dfs(code, goal_i, cur_r, goal_c, append(path, getMove(LEFT, cur_c-goal_c)...), res)
	}
	if cur_r > goal_r {
		dfs(code, goal_i, goal_r, cur_c, append(path, getMove(UP, cur_r-goal_r)...), res)
	}
	if cur_c < goal_c {
		dfs(code, goal_i, cur_r, goal_c, append(path, getMove(RIGHT, goal_c-cur_c)...), res)
	}
}

func getPossibleDirCodesForNumCode(num_code string) [][]int {
	dir_codes := [][]int{}

	dfs(num_code, 0, 3, 2, []int{}, &dir_codes)

	return dir_codes
}

func getMoves(from, to int) []int {
	switch from {
	case A:
		switch to {
		case RIGHT:
			return []int{DOWN}
		case UP:
			return []int{LEFT}
		case DOWN:
			return []int{DOWN, LEFT}
		case LEFT:
			return []int{DOWN, LEFT, LEFT}
		}

	case DOWN:
		if to == A {
			return []int{RIGHT, UP}
		}
		return []int{to}

	case UP:
		switch to {
		case A:
			return []int{RIGHT}
		case DOWN:
			return []int{DOWN}
		default:
			return []int{DOWN, to}
		}

	case LEFT:
		return append([]int{RIGHT}, getMoves(DOWN, to)...)

	case RIGHT:
		if to == A {
			return []int{UP}
		}
		return append([]int{LEFT}, getMoves(DOWN, to)...)
	}

	panic(fmt.Errorf("getMoves: unreachable"))
}

func expandDirCode(dir_code []int, position, depth int) (expanded_len int) {
	moves := []int{}

	pos := position
	for _, cur_code := range dir_code {
		if cur_code == pos {
			moves = append(moves, A)
			continue
		}

		cur_moves := getMoves(pos, cur_code)

		moves = append(moves, cur_moves...)
		moves = append(moves, A)
		pos = cur_code
	}

	if depth > 0 {
		expanded_len = 0
		pos := position
		for _, code := range moves {
			exp_len := expandDirCode([]int{code}, pos, depth-1)
			expanded_len += exp_len
			pos = code
		}
		return expanded_len
	}

	// for _, c := range moves {
	// 	fmt.Printf("%v", formatCode(c))
	// }

	expanded_len = len(moves)
	return expanded_len
}

func getMinComplexityForNumCode(code string, depth int) int {
	min_complexity := math.MaxInt

	num_codes := getPossibleDirCodesForNumCode(code)
	// fmt.Println(code)
	// for _, code := range num_codes {
	// 	fmt.Println(formatDirCode(code))
	// }
	// fmt.Println("..")

	for _, num_code := range num_codes {
		cur_complexity := expandDirCode(num_code, A, depth)
		min_complexity = min(min_complexity, cur_complexity)
	}

	return min_complexity
}

//
//
//

func GetComplexity(num_codes []string, depth int) int {
	complx_2 := 0

	for _, code := range num_codes {
		num := getNumFromCode(code)
		complx_2 += getMinComplexityForNumCode(code, depth) * num
	}

	// fmt.Println("total complexity 2:", complx_2)
	return complx_2
}

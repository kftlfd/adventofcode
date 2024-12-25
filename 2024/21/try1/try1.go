package try1

import (
	"fmt"
	"strconv"
	"strings"
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

func dirCodeFromNumCode(numerical_code string) []int {
	r := 3
	c := 2

	moves := []int{}

	for _, code := range strings.Split(numerical_code, "") {
		coords := getNumericCoordinates(code)

		if r == 3 {
			for r > coords[R] {
				moves = append(moves, UP)
				r -= 1
			}
			for c > coords[C] {
				moves = append(moves, LEFT)
				c -= 1
			}
		} else {
			for c > coords[C] {
				moves = append(moves, LEFT)
				c -= 1
			}
			for r > coords[R] {
				moves = append(moves, UP)
				r -= 1
			}
		}

		for c < coords[C] {
			moves = append(moves, RIGHT)
			c += 1
		}
		for r < coords[R] {
			moves = append(moves, DOWN)
			r += 1
		}

		moves = append(moves, A)
	}

	return moves
}

func dirCodeFromDirCode(code []int) []int {
	/*
		  ^ A
		< V >
	*/

	moves := []int{}

	pending := [4]int{}

	move := func(code int) {
		moves = append(moves, code)
	}

	press := func(code int) {
		for i := 0; i < pending[code]; i++ {
			moves = append(moves, A)
		}
		pending[code] = 0
	}

	for _, cur_code := range code {
		if cur_code != A {
			pending[cur_code] += 1
			continue
		}

		// starting at A, applying pending, ending back at A

		if pending[LEFT] > 0 {
			move(LEFT)
			press(UP)

			move(LEFT)
			move(DOWN)
			press(LEFT)

			move(RIGHT)
			press(DOWN)

			move(RIGHT)
			press(RIGHT)

			move(UP)
			moves = append(moves, A)
			continue
		}

		if pending[DOWN] > 0 || (pending[UP] > 0 && pending[RIGHT] > 0) {
			move(LEFT)
			press(UP)

			move(DOWN)
			press(DOWN)

			move(RIGHT)
			press(RIGHT)

			move(UP)
			moves = append(moves, A)
			continue
		}

		if pending[UP] > 0 {
			move(LEFT)
			press(UP)
			move(RIGHT)
			moves = append(moves, A)
			continue
		}

		if pending[RIGHT] > 0 {
			move(DOWN)
			press(RIGHT)
			move(UP)
			moves = append(moves, A)
			continue
		}

		moves = append(moves, A)
	}

	return moves
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

//
//
//

func GetComplexity(num_codes []string, depth int) int {
	complexity := 0

	for _, code := range num_codes {
		seq_num := dirCodeFromNumCode(code)
		seq_len := expandDirCode(seq_num, A, depth)
		num := getNumFromCode(code)
		complexity += seq_len * num
		// fmt.Printf("code: %v  seq_len: %v  num: %v\n", code, seq_len, num)
	}

	// fmt.Println("total complexity 1:", complexity)
	return complexity
}

func PrintCode(code string) {
	c := code
	num_c := dirCodeFromNumCode(c)
	dir_1 := dirCodeFromDirCode(num_c)
	dir_2 := dirCodeFromDirCode(dir_1)
	fmt.Println(c)
	fmt.Println(FormatDirCode(num_c))
	fmt.Println(FormatDirCode(dir_1))
	fmt.Println(FormatDirCode(dir_2))
}

func FormatCode(c int) string {
	switch c {
	case UP:
		return "^"
	case DOWN:
		return "v"
	case LEFT:
		return "<"
	case RIGHT:
		return ">"
	case A:
		return "A"
	}
	panic(fmt.Errorf("unknown code"))
}

func FormatDirCode(dir_code []int) []string {
	repr := []string{}
	for _, c := range dir_code {
		switch c {
		case UP:
			repr = append(repr, "^")
		case DOWN:
			repr = append(repr, "v")
		case LEFT:
			repr = append(repr, "<")
		case RIGHT:
			repr = append(repr, ">")
		case A:
			repr = append(repr, "A")
		}
	}
	return repr
}

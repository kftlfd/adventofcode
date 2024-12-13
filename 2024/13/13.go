package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	X = 0
	Y = 1
)

func main() {
	inp_path := os.Args[1]
	inp_file, err := os.Open(inp_path)
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	reader := bufio.NewScanner(inp_file)
	reader.Split(bufio.ScanLines)

	machines := []Machine{}
	cur_machine := Machine{}
	conversion_error := 10000000000000
	i := 0

	for reader.Scan() {
		row := strings.Split(reader.Text(), " ")
		if len(row) < 3 {
			continue
		}

		if i == 2 {
			// Prize
			if len(row) != 3 {
				panic(fmt.Errorf("bad prize row: %v", row))
			}
			x, errx := strconv.Atoi(row[1][2 : len(row[1])-1])
			y, erry := strconv.Atoi(row[2][2:])
			if errx != nil || erry != nil {
				panic(fmt.Errorf("bad prize inputs: %v %v %v", row, errx, erry))
			}

			cur_machine.prize = [2]int{x, y}
			cur_machine.prize2 = [2]int{x + conversion_error, y + conversion_error}
			machines = append(machines, cur_machine)
			cur_machine = Machine{}
		} else {
			// Button
			if len(row) != 4 {
				panic(fmt.Errorf("bad button row: %v", row))
			}
			x, errx := strconv.Atoi(row[2][2 : len(row[2])-1])
			y, erry := strconv.Atoi(row[3][2:])
			if errx != nil || erry != nil {
				panic(fmt.Errorf("bad prize inputs: %v %v %v", row, errx, erry))
			}
			if x < 0 || y < 0 {
				fmt.Println("have negative x/y")
			}

			if i == 0 {
				cur_machine.a = [2]int{x, y}
			} else {
				cur_machine.b = [2]int{x, y}
			}
		}

		i = (i + 1) % 3
	}

	// for _, m := range machines {
	// 	fmt.Println(m)
	// }

	min_tokens_for_prizes := 0
	tokens_corrected := 0
	for _, m := range machines {
		tokens := m.minTokens1()
		// tokens := m.minTokens2()
		if tokens > -1 {
			min_tokens_for_prizes += tokens
		}
		tok_mat := minTokensMat(m.a, m.b, m.prize2)
		if tok_mat > -1 {
			tokens_corrected += tok_mat
		}
	}
	fmt.Println("min tokens for prizes:", min_tokens_for_prizes)
	fmt.Println("min tokens corrected:", tokens_corrected)
}

type Machine struct {
	a, b, prize, prize2 [2]int
}

func (m Machine) minTokens1() int {
	min_moves := 1
	max_moves := 200

	for moves := min_moves; moves <= max_moves; moves++ {
		for a_n := 0; a_n <= moves; a_n++ {
			b_n := moves - a_n
			x := m.a[0]*a_n + m.b[0]*b_n
			y := m.a[1]*a_n + m.b[1]*b_n
			if x == m.prize[0] && y == m.prize[1] {
				tokens := a_n*3 + b_n
				return tokens
			}
		}
	}

	return -1
}

func boolToInt(b bool) int {
	i := 0
	if b {
		i = 1
	}
	return i
}

func intCeil(a, b int) int {
	return (a / b) + boolToInt(a%b != 0)
}

func (m Machine) minTokens2() int {
	a_moves := max(intCeil(m.prize2[0], m.a[0]), intCeil(m.prize2[1], m.a[1]))
	b_moves := max(intCeil(m.prize2[0], m.b[0]), intCeil(m.prize2[1], m.b[1]))
	min_moves := a_moves
	max_moves := b_moves
	if a_moves > b_moves {
		min_moves = b_moves
		max_moves = a_moves
	}

	ans := -1
	for moves := min_moves; moves <= max_moves; moves++ {
		for a_n := 0; a_n <= moves; a_n++ {
			b_n := moves - a_n
			x := m.a[0]*a_n + m.b[0]*b_n
			y := m.a[1]*a_n + m.b[1]*b_n
			if x == m.prize[0] && y == m.prize[1] {
				tokens := a_n*3 + b_n
				if ans == -1 {
					ans = tokens
				} else {
					ans = min(ans, tokens)
				}
			}
		}
	}
	return ans
}

func minTokensMat(A, B, P [2]int) int {
	// x_is_int := (B[X]*P[Y]-B[Y]*P[X])%(A[Y]*B[X]-A[X]*B[Y]) == 0
	// y_is_int := (A[Y]*P[X]-A[X]*P[Y])%(A[Y]*B[X]-A[X]*B[Y]) == 0

	// if !x_is_int || !y_is_int {
	// 	return -1
	// }

	// x := (B[X]*P[Y] - B[Y]*P[X]) / (A[Y]*B[X] - A[X]*B[Y])
	// y := (A[Y]*P[X] - A[X]*P[Y]) / (A[Y]*B[X] - A[X]*B[Y])

	// return 3*x + y

	x_top := B[X]*P[Y] - B[Y]*P[X]
	y_top := A[Y]*P[X] - A[X]*P[Y]
	div := A[Y]*B[X] - A[X]*B[Y]

	x_is_int := x_top%div == 0
	y_is_int := y_top%div == 0

	if !x_is_int || !y_is_int {
		return -1
	}

	x := x_top / div
	y := y_top / div

	return 3*x + y
}

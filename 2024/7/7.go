package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	inpPath := os.Args[1]
	inpFile, err := os.Open(inpPath)
	if err != nil {
		panic(err)
	}

	fileScanner := bufio.NewScanner(inpFile)
	fileScanner.Split(bufio.ScanLines)

	eqs := []Eq{}

	for fileScanner.Scan() {
		row := fileScanner.Text()
		arr := strings.Split(row, ": ")
		if len(arr) != 2 {
			panic(fmt.Errorf("invalid input: %s", row))
		}

		target, err := strconv.Atoi(arr[0])
		if err != nil {
			panic(fmt.Errorf("bad target: %s", row))
		}

		val_strings := strings.Split(arr[1], " ")
		values := make([]int, len(val_strings))
		for i, val := range val_strings {
			val_int, err := strconv.Atoi(val)
			if err != nil {
				panic(fmt.Errorf("bad val: %s\t row: %s", val, row))
			}
			values[i] = val_int
		}

		eqs = append(eqs, Eq{target: target, values: values})
	}

	// for _, v := range eqs {
	// 	fmt.Println(v)
	// }

	valid_sum := 0
	valid_w_concat_sum := 0
	for _, eq := range eqs {
		if eq.isPossible() {
			valid_sum += eq.target
		}
		if eq.isPossible2() {
			valid_w_concat_sum += eq.target
		}
	}
	fmt.Println("valid sum:", valid_sum)
	fmt.Println("valid sum w concat:", valid_w_concat_sum)
}

type Eq struct {
	target int
	values []int
}

func (e *Eq) isPossible() bool {
	n := len(e.values)
	q := [][2]int{
		{e.values[0], 1},
	}

	for len(q) > 0 {
		q_val := q[0]
		q = q[1:]
		cur_val := q_val[0]
		cur_i := q_val[1]
		if cur_i >= n {
			if cur_val == e.target {
				return true
			}
			continue
		}
		q = append(q, [2]int{cur_val + e.values[cur_i], cur_i + 1})
		q = append(q, [2]int{cur_val * e.values[cur_i], cur_i + 1})
	}

	return false
}

func (e *Eq) isPossible2() bool {
	n := len(e.values)
	q := [][2]int{
		{e.values[0], 1},
	}

	for len(q) > 0 {
		q_val := q[0]
		q = q[1:]
		cur_val := q_val[0]
		cur_i := q_val[1]
		if cur_i >= n {
			if cur_val == e.target {
				return true
			}
			continue
		}
		q = append(q, [2]int{cur_val + e.values[cur_i], cur_i + 1})
		q = append(q, [2]int{cur_val * e.values[cur_i], cur_i + 1})
		q = append(q, [2]int{concat(cur_val, e.values[cur_i]), cur_i + 1})
	}

	return false
}

func concat(v1, v2 int) int {
	val, err := strconv.Atoi(
		strings.Join(
			[]string{
				strconv.Itoa(v1), strconv.Itoa(v2),
			},
			"",
		),
	)
	if err != nil {
		panic(fmt.Errorf("cant concat numbers: %v %v", v1, v2))
	}
	return val
}

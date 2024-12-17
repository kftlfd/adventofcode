package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type Computer struct {
	A, B, C, instruction int
}

func (c *Computer) run(program []int) []int {
	out := []int{}
	program_len := len(program)

	for true {
		if c.instruction < 0 {
			panic(fmt.Errorf("negative instruction index"))
		}
		if c.instruction >= program_len {
			// halt
			break
		}
		if c.instruction+1 >= program_len {
			panic(fmt.Errorf("operand index out of bounds"))
		}

		opcode := program[c.instruction]
		operand := program[c.instruction+1]

		switch opcode {
		case 0:
			c.adv_0(operand)
		case 1:
			c.bxl_1(operand)
		case 2:
			c.bst_2(operand)
		case 3:
			c.jnz_3(operand)
		case 4:
			c.bxc_4()
		case 5:
			out = append(out, c.out_5(operand))
		case 6:
			c.bdv_6(operand)
		case 7:
			c.cdv_7(operand)
		default:
			panic(fmt.Errorf("unrecognized opcode"))
		}

		c.instruction += 2
	}

	return out
}

func (c *Computer) getComboOperand(operand int) int {
	if operand < 0 || operand > 6 {
		panic(fmt.Errorf("bad combo operand: %v", operand))
	}

	if operand < 4 {
		return operand
	}

	switch operand {
	case 4:
		return c.A
	case 5:
		return c.B
	case 6:
		return c.C
	}

	panic(fmt.Errorf("unreachable"))
}

func (c *Computer) dv(combo int) int {
	operand := c.getComboOperand(combo)
	return c.A / int(math.Pow(2, float64(operand)))
}

func (c *Computer) adv_0(combo int) {
	c.A = c.dv(combo)
}

func (c *Computer) bxl_1(literal int) {
	c.B = c.B ^ literal
}

func (c *Computer) bst_2(combo int) {
	c.B = c.getComboOperand(combo) % 8
}

func (c *Computer) jnz_3(literal int) {
	if c.A == 0 {
		return
	}
	c.instruction = literal - 2
}

func (c *Computer) bxc_4() {
	c.B = c.B ^ c.C
}

func (c *Computer) out_5(combo int) int {
	return c.getComboOperand(combo) % 8
}

func (c *Computer) bdv_6(combo int) {
	c.B = c.dv(combo)
}

func (c *Computer) cdv_7(combo int) {
	c.C = c.dv(combo)
}

//
//
//

func parseInput(file *os.File) (Computer, []int) {
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)

	computer := Computer{}
	program := []int{}

	is_program := false
	reg := 0

	for reader.Scan() {
		line := reader.Text()

		if len(line) < 1 {
			is_program = true
			continue
		}

		if is_program {
			vals := strings.Split(line, " ")
			if len(vals) != 2 {
				panic(fmt.Errorf("bad input: %v", line))
			}
			vals = strings.Split(vals[1], ",")
			for _, val := range vals {
				val_int, err := strconv.Atoi(val)
				if err != nil {
					panic(err)
				}
				program = append(program, val_int)
			}
			continue
		}

		vals := strings.Split(line, " ")
		if len(vals) != 3 {
			panic(fmt.Errorf("bad input: %v", line))
		}
		reg_val, err := strconv.Atoi(vals[2])
		if err != nil {
			panic(err)
		}
		switch reg {
		case 0:
			computer.A = reg_val
		case 1:
			computer.B = reg_val
		case 2:
			computer.C = reg_val
		}
		reg += 1
	}

	return computer, program
}

func main() {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	computer, program := parseInput(inp_file)

	fmt.Println("--- part 1 ---")
	fmt.Printf("computer: %+v\n", computer)
	fmt.Printf("program:  %+v\n", program)

	output := computer.run(program)

	out := []string{}
	for _, val := range output {
		out = append(out, strconv.Itoa(val))
	}

	out_string := strings.Join(out, ",")

	fmt.Printf("out:      %v\n", out_string)

	//
	// part 2
	//

	/*
		Not a general solution, solving only for this particular program: 2,4,1,6,7,5,4,4,1,7,0,3,5,5,3,0

		Key points after analyzing the program workings:
			1. there is only one jump instructions at the end (3,0), which means this runs like a do-while loop (while A > 0)
			2. there is only one out instruction (5,5), which prints the B%8
			3. we have only one instruction that modifies A -> (0,3): A = A / 8
			4. all other instructions set B and C depending on the current value of A

		So, essentially this program does this:
			do {
				print( f(A) )
				A = A / 8
			} while (A > 0)
		where:
			f(A) = (((A % 8) ^ 6) ^ (A / (2**((A%8)^6))) ^ 7) % 8

		We know that this program should output itself, so it should halt after 16 runs => after 16 runs A / 8 == 0
		we use integer division, so the final value of A (haltA) is in the range [0, 7]
		and the f(haltA) == 0 (the final value in program)

		Now we can go backwards in the program values and see if the f(A) eaquals that value,
		if yes, we update A = A * 8 + [0, 7], and for each check the next(previous value in program)
	*/

	fmt.Printf("\n--- part 2 ---\n")

	prog := []int{2, 4, 1, 6, 7, 5, 4, 4, 1, 7, 0, 3, 5, 5, 3, 0}

	fmt.Printf("program  : %+v\n", prog)

	resA := math.MaxInt
	for haltA := 0; haltA <= 7; haltA++ {
		dfs(haltA, len(prog)-1, prog, &resA)
	}

	if resA != math.MaxInt {
		fmt.Printf("A to self copy: %v\n\n", resA)
		c := Computer{A: resA, B: 0, C: 0, instruction: 0}
		fmt.Printf("computer: %+v\n", c)
		fmt.Println("inp:", prog)
		fmt.Println("out:", c.run(prog))
	} else {
		fmt.Println("A to self copy: not found")
	}
}

func dfs(A, out_i int, prog []int, ans *int) {
	B := (A % 8) ^ 6
	C := A / int(math.Pow(2, float64(B)))
	X := (B ^ C ^ 7) % 8

	if X != prog[out_i] {
		return
	}

	if out_i == 0 {
		(*ans) = min(*ans, A)
		return
	}

	for A_d := 0; A_d < 8; A_d++ {
		dfs(A*8+A_d, out_i-1, prog, ans)
	}
}

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

	rules := make(Rules)
	updates := [][]string{}

	rules2 := make(Rules2)

	isUpdates := false
	for fileScanner.Scan() {
		row := fileScanner.Text()
		if len(row) < 1 {
			isUpdates = true
			continue
		}

		if isUpdates {
			updates = append(updates, strings.Split(row, ","))
			continue
		}

		ruleRow := strings.Split(row, "|")
		if (len(ruleRow)) != 2 {
			panic(fmt.Errorf("invalid rule: %v", row))
		}

		rule0 := rules.GetOrCreate(ruleRow[0])
		rule0.before = append(rule0.before, ruleRow[1])

		rule1 := rules.GetOrCreate(ruleRow[1])
		rule1.after = append(rule1.after, ruleRow[0])

		rules2.SetBefore(ruleRow[0], ruleRow[1])
	}

	// fmt.Printf("%+v \n%+v \n", rules, updates)
	// for k, v := range rules {
	// 	fmt.Println(k, *v)
	// }

	// for k, v := range rules2 {
	// 	fmt.Printf("%v: %+v\n", k, v)
	// }

	correctCount := 0
	correctMidSums := 0
	correctedCount := 0
	correctedMidSums := 0

	r2_correctCount := 0
	r2_correctMidSums := 0
	r2_correctedCount := 0
	r2_correctedMidSums := 0

	for _, update := range updates {
		if isCorrect(update, rules) {
			correctCount += 1
			correctMidSums += getMiddleNum(update)
		} else {
			correctedCount += 1
			correctedMidSums += getMiddleNum(correct(update, rules))
		}

		if rules2.IsCorrect(update) {
			r2_correctCount += 1
			r2_correctMidSums += getMiddleNum(update)
		} else {
			r2_correctedCount += 1
			r2_correctedMidSums += getMiddleNum(rules2.ToCorrected(update))
		}
	}

	fmt.Println("correct updates:", correctCount)
	fmt.Println("middles sum:", correctMidSums)
	fmt.Println("corrected updates:", correctedCount)
	fmt.Println("corrected middles sum:", correctedMidSums)

	fmt.Println("r2_correct updates:", r2_correctCount)
	fmt.Println("r2_middles sum:", r2_correctMidSums)
	fmt.Println("r2_corrected updates:", r2_correctedCount)
	fmt.Println("r2_corrected middles sum:", r2_correctedMidSums)
}

type Rules map[string]*Rule

type Rule struct {
	before, after []string
}

func (r *Rules) GetOrCreate(key string) *Rule {
	rule, ok := (*r)[key]
	if !ok {
		rule = new(Rule)
		(*r)[key] = rule
	}
	return rule
}

func (r *Rules) isBefore(v1, v2 string) bool {
	r1 := r.GetOrCreate(v1)
	r2 := r.GetOrCreate(v2)
	return !isIn(v1, r2.before) && !isIn(v2, r1.after)
}

func isCorrect(update []string, rules Rules) bool {
	for i, val := range update {
		r, ok := rules[val]
		if !ok {
			continue
		}

		ok = true

		for _, prev := range update[:i] {
			if isIn(prev, r.before) {
				ok = false
				break
			}
		}
		if !ok {
			return false
		}

		for _, nxt := range update[i+1:] {
			if isIn(nxt, r.after) {
				ok = false
				break
			}
		}
		if !ok {
			return false
		}
	}
	return true
}

func isIn(val string, arr []string) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

func getMiddleNum(update []string) int {
	if len(update)%2 != 1 {
		panic(fmt.Errorf("update of even len: %+v", update))
	}

	midEl := update[len(update)/2]

	val, err := strconv.Atoi(midEl)
	if err != nil {
		panic(err)
	}

	return val
}

func correct(incorrectUpdate []string, rules Rules) []string {
	n := len(incorrectUpdate)
	update := make([]string, n)
	copy(update, incorrectUpdate)
	swaps := true

	for swaps {
		swaps = false
		for i, v1 := range update {
			for j := i + 1; j < n; j++ {
				v2 := update[j]
				if rules.isBefore(v1, v2) {
					continue
				} else {
					update[i], update[j] = update[j], update[i]
					swaps = true
				}
			}
		}
	}

	return update
}

// solution 2

type Rules2 map[string][]string

func (r *Rules2) SetBefore(v1, v2 string) {
	arr := (*r)[v1]
	arr = append(arr, v2)
	(*r)[v1] = arr
}

func (r *Rules2) MustBeBefore(v1, v2 string) bool {
	arr := (*r)[v1]
	for _, val := range arr {
		if val == v2 {
			return true
		}
	}
	return false
}

func (r *Rules2) IsCorrect(update []string) bool {
	n := len(update)
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			if r.MustBeBefore(update[j], update[i]) {
				return false
			}
		}
	}
	return true
}

func (r *Rules2) ToCorrected(incorrectUpdate []string) []string {
	n := len(incorrectUpdate)
	update := make([]string, n)
	copy(update, incorrectUpdate)
	swaps := true
	for swaps {
		swaps = false
		for i := 0; i < n; i++ {
			for j := i + 1; j < n; j++ {
				if r.MustBeBefore(update[j], update[i]) {
					update[i], update[j] = update[j], update[i]
					swaps = true
				}
			}
		}
	}
	return update
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	inpPath := os.Args[1]
	inpFile, err := os.Open(inpPath)
	if err != nil {
		panic(err)
	}
	defer inpFile.Close()

	fileScanner := bufio.NewScanner(inpFile)
	fileScanner.Split(bufio.ScanLines)

	stones := []string{}

	for fileScanner.Scan() {
		row := strings.Split(fileScanner.Text(), " ")
		for _, val := range row {
			stones = append(stones, val)
		}
	}

	fmt.Println("stones:", stones)

	// inMemoryTransforms(stones, 25)

	// inFilesTransforms(stones, 25)

	dfsMemoTransforms(stones, 75)

	fmt.Println(time.Since(start))
}

//
// in memory
//

func inMemoryTransforms(stones []string, transforms int) {
	for i := 0; i < transforms; i++ {
		stones = blink(stones)
	}
	fmt.Printf("[in memory] after %v blinks: %v stones\n", transforms, len(stones))
}

func blink(stones []string) []string {
	nxt := []string{}

	for _, stone := range stones {
		if stone == "0" {
			nxt = append(nxt, "1")
			continue
		}

		if len(stone)%2 == 0 {
			s1 := stone[:len(stone)/2]
			num2, err := strconv.Atoi(stone[len(stone)/2:])
			if err != nil {
				panic(fmt.Errorf("not a number: %v, when splitting: %v", stone[len(stone)/2:], stone))
			}
			s2 := strconv.Itoa(num2)
			nxt = append(nxt, s1, s2)
			continue
		}

		num, err := strconv.Atoi(stone)
		if err != nil {
			panic(fmt.Errorf("not a number: %v", stone))
		}
		nxt = append(nxt, strconv.Itoa(num*2024))
	}

	return nxt
}

//
// in files
//

func inFilesTransforms(stones []string, transforms int) {
	a := "a.txt"
	b := "b.txt"

	a_f, err := os.Create(a)
	if err != nil {
		panic(err)
	}
	defer a_f.Close()

	writer := bufio.NewWriter(a_f)
	for _, stone := range stones {
		if _, err := writer.WriteString(stone + "\n"); err != nil {
			panic(err)
		}
	}
	if err := writer.Flush(); err != nil {
		panic(err)
	}

	n := transforms
	for i := 0; i < n; i++ {
		blinkFile(a, b)
		a, b = b, a
	}

	stones_count := 0
	out_file, err := os.Open(a)
	if err != nil {
		panic(err)
	}
	defer out_file.Close()
	out_reader := bufio.NewScanner(out_file)
	out_reader.Split(bufio.ScanLines)
	for out_reader.Scan() {
		stones_count += 1
	}
	fmt.Println("[in files] after", n, "blinks:", stones_count)
}

func transform(stone string) []string {
	if stone == "0" {
		return []string{"1"}
	}

	if len(stone)%2 == 0 {
		s1 := stone[:len(stone)/2]
		num2, err := strconv.Atoi(stone[len(stone)/2:])
		if err != nil {
			panic(fmt.Errorf("not a number: %v, when splitting: %v", stone[len(stone)/2:], stone))
		}
		s2 := strconv.Itoa(num2)
		return []string{s1, s2}
	}

	num, err := strconv.Atoi(stone)
	if err != nil {
		panic(fmt.Errorf("not a number: %v", stone))
	}
	return []string{strconv.Itoa(num * 2024)}
}

func blinkFile(inp, out string) {
	inputFile, err := os.OpenFile(inp, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer inputFile.Close()

	outputFile, err := os.Create(out)
	if err != nil {
		panic(err)
	}
	defer outputFile.Close()

	reader := bufio.NewScanner(inputFile)
	reader.Split(bufio.ScanLines)

	writer := bufio.NewWriter(outputFile)

	writes := 0

	for reader.Scan() {
		stone := reader.Text()
		for _, nxt := range transform(stone) {
			if _, err := writer.WriteString(nxt + "\n"); err != nil {
				panic(err)
			}
			writes += 1
		}
		if writes > 10000 {
			if err := writer.Flush(); err != nil {
				panic(err)
			}
			writes = 0
		}
	}

	if writes > 0 {
		if err := writer.Flush(); err != nil {
			panic(err)
		}
	}
}

//
// dfs + memo
//

func getKey(stone string, i int) string {
	return strings.Join([]string{stone, "-", strconv.Itoa(i)}, "")
}

func dfs(stone string, i, n int, memo *map[string]int) int {
	if i >= n {
		return 1
	}

	key := getKey(stone, i)
	memo_val, ok := (*memo)[key]
	if ok {
		return memo_val
	}

	cur_val := 0
	for _, nxt := range transform(stone) {
		cur_val += dfs(nxt, i+1, n, memo)
	}

	(*memo)[key] = cur_val
	return cur_val
}

func dfsMemoTransforms(stones []string, transforms int) {
	count := 0

	memo := make(map[string]int)

	for _, stone := range stones {
		count += dfs(stone, 0, transforms, &memo)
	}

	fmt.Println("[dfs+memo] after", transforms, "blinks:", count)
}

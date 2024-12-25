package try3

import (
	"math"
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

var numpadKeyIdx = map[string]int{
	"0": 0,
	"1": 1,
	"2": 2,
	"3": 3,
	"4": 4,
	"5": 5,
	"6": 6,
	"7": 7,
	"8": 8,
	"9": 9,
	"A": 10,
}

// var numpadIdxKey = map[int]string{
// 	0:  "0",
// 	1:  "1",
// 	2:  "2",
// 	3:  "3",
// 	4:  "4",
// 	5:  "5",
// 	6:  "6",
// 	7:  "7",
// 	8:  "8",
// 	9:  "9",
// 	10: "A",
// }

var numpadCoords = [11][2]int{
	{3, 1}, // 0
	{2, 0}, // 1
	{2, 1}, // 2
	{2, 2}, // 3
	{1, 0}, // 4
	{1, 1}, // 5
	{1, 2}, // 6
	{0, 0}, // 7
	{0, 1}, // 8
	{0, 2}, // 9
	{3, 2}, // A
}

// paths from key to key
var numpadMap = [11][11][]string{}

func numpadDfs(from, to, cur_r, cur_c, target_r, target_c int, cur_path []string) {
	if cur_r == 3 && cur_c == 0 {
		// at the gap
		return
	}

	if cur_r == target_r && cur_c == target_c {
		numpadMap[from][to] = append(numpadMap[from][to], strings.Join(append(cur_path, "A"), ""))
		return
	}

	if cur_r < target_r {
		numpadDfs(from, to, cur_r+1, cur_c, target_r, target_c, append(cur_path, "v"))
	}
	if cur_r > target_r {
		numpadDfs(from, to, cur_r-1, cur_c, target_r, target_c, append(cur_path, "^"))
	}
	if cur_c < target_c {
		numpadDfs(from, to, cur_r, cur_c+1, target_r, target_c, append(cur_path, ">"))
	}
	if cur_c > target_c {
		numpadDfs(from, to, cur_r, cur_c-1, target_r, target_c, append(cur_path, "<"))
	}
}

func populateNumpadMap() {
	for from := 0; from < 11; from++ {
		for to := 0; to < 11; to++ {
			numpadMap[from][to] = make([]string, 0)
		}
	}
	for from := 0; from < 11; from++ {
		for to := 0; to < 11; to++ {
			from_coords := numpadCoords[from]
			to_coords := numpadCoords[to]
			numpadDfs(from, to, from_coords[0], from_coords[1], to_coords[0], to_coords[1], []string{})
		}
	}
}

var dirpadMap = [5][5][]string{}

var dirpadKeyIdx = map[string]int{
	"^": 0,
	"v": 1,
	"<": 2,
	">": 3,
	"A": 4,
}

// var dirpadIdxKey = map[int]string{
// 	0: "^",
// 	1: "v",
// 	2: "<",
// 	3: ">",
// 	4: "A",
// }

var dirpadCoords = [5][2]int{
	{0, 1}, // ^
	{1, 1}, // v
	{1, 0}, // <
	{1, 2}, // >
	{0, 2}, // A
}

func dirpadDfs(from, to, cur_r, cur_c, target_r, target_c int, cur_path []string) {
	if cur_r == 0 && cur_c == 0 {
		// at the gap
		return
	}

	if cur_r == target_r && cur_c == target_c {
		dirpadMap[from][to] = append(dirpadMap[from][to], strings.Join(append(cur_path, "A"), ""))
		return
	}

	if cur_r < target_r {
		dirpadDfs(from, to, cur_r+1, cur_c, target_r, target_c, append(cur_path, "v"))
	}
	if cur_r > target_r {
		dirpadDfs(from, to, cur_r-1, cur_c, target_r, target_c, append(cur_path, "^"))
	}
	if cur_c < target_c {
		dirpadDfs(from, to, cur_r, cur_c+1, target_r, target_c, append(cur_path, ">"))
	}
	if cur_c > target_c {
		dirpadDfs(from, to, cur_r, cur_c-1, target_r, target_c, append(cur_path, "<"))
	}
}

func populateDirpadMap() {
	for from := 0; from < 5; from++ {
		for to := 0; to < 5; to++ {
			dirpadMap[from][to] = make([]string, 0)
		}
	}
	for from := 0; from < 5; from++ {
		for to := 0; to < 5; to++ {
			from_coords := dirpadCoords[from]
			to_coords := dirpadCoords[to]
			dirpadDfs(from, to, from_coords[0], from_coords[1], to_coords[0], to_coords[1], []string{})
		}
	}
}

func inputSequencesDfs(code string, idx int, prev_key string, path []string, result *[][]string) {
	if idx >= len(code) {
		(*result) = append((*result), path)
		return
	}

	cur_key := code[idx : idx+1]

	// fmt.Printf("%v -> %v: %+v\n", prev_key, cur_key, numpadMap[numpadKeyIdx[prev_key]][numpadKeyIdx[cur_key]])

	for _, nxt_move := range numpadMap[numpadKeyIdx[prev_key]][numpadKeyIdx[cur_key]] {
		inputSequencesDfs(code, idx+1, cur_key, append(path, nxt_move), result)
	}
}

func getInputSequences(code string) [][]string {
	sequences := [][]string{}
	inputSequencesDfs(code, 0, "A", []string{}, &sequences)

	// fmt.Printf("%v:\n%+v\n", code, sequences)
	return sequences
}

func dirSequenceDfs(code string, idx int, prev_key string, path []string, result *[][]string) {
	if idx >= len(code) {
		(*result) = append((*result), path)
		return
	}

	cur_key := code[idx : idx+1]

	for _, nxt_move := range dirpadMap[dirpadKeyIdx[prev_key]][dirpadKeyIdx[cur_key]] {
		dirSequenceDfs(code, idx+1, cur_key, append(path, nxt_move), result)
	}
}

func getDirSequence(code string) [][]string {
	seq := [][]string{}
	dirSequenceDfs(code, 0, "A", []string{}, &seq)
	return seq
}

func getMinSeqLenAtDepth(code string, depth int) int {
	cache := make(Cache)

	input_seqs := getInputSequences(code)

	min_seq_len := math.MaxInt

	for _, seq := range input_seqs {
		min_seq_len = min(min_seq_len, shortestSeqLen(seq, depth, &cache))
	}

	return min_seq_len
}

type Cache map[int]*map[string]int

func (c *Cache) get(depth int, keys string) (int, bool) {
	kmap, ok := (*c)[depth]
	if !ok {
		return 0, false
	}
	val, ok := (*kmap)[keys]
	if !ok {
		return 0, false
	}
	return val, true
}

func (c *Cache) set(depth int, keys string, val int) {
	kmap, ok := (*c)[depth]
	if ok {
		(*kmap)[keys] = val
	} else {
		m := make(map[string]int)
		(*c)[depth] = &m
		kmap := (*c)[depth]
		(*kmap)[keys] = val
	}
}

func shortestSeqLen(keys []string, depth int, cache *Cache) int {
	if depth == 0 {
		keys_len := 0
		for _, key := range keys {
			keys_len += len(key)
		}
		return keys_len
	}

	if val, ok := cache.get(depth, strings.Join(keys, "")); ok {
		return val
	}

	seq_len := 0

	for _, key := range keys {
		key_seqs := getDirSequence(key)
		min_seq_len := math.MaxInt
		for _, seq := range key_seqs {
			min_seq_len = min(min_seq_len, shortestSeqLen(seq, depth-1, cache))
		}
		seq_len += min_seq_len
	}

	cache.set(depth, strings.Join(keys, ""), seq_len)

	return seq_len
}

//
//
//

func GetComplexity(num_codes []string, depth int) int {
	populateDirpadMap()
	// // fmt.Printf("%+v\n", dirpadMap)
	// fmt.Printf("\ndirpad\n")
	// for from, from_arr := range dirpadMap {
	// 	for to, to_arr := range from_arr {
	// 		fmt.Printf("%v -> %v: %+v\n", dirpadIdxKey[from], dirpadIdxKey[to], to_arr)
	// 	}
	// }

	populateNumpadMap()
	// // fmt.Printf("%+v\n", dirpadMap)
	// fmt.Printf("\nnumpad\n")
	// for from, from_arr := range numpadMap {
	// 	for to, to_arr := range from_arr {
	// 		fmt.Printf("%v -> %v: %+v\n", numpadIdxKey[from], numpadIdxKey[to], to_arr)
	// 	}
	// }

	// for _, seq := range getInputSequences("0279A") {
	// 	fmt.Println(seq)
	// 	for _, code := range seq {
	// 		dir_seq := getDirSequence(code)
	// 		fmt.Printf("%v: %+v\n", code, dir_seq)
	// 	}
	// }

	// fmt.Printf("\n\n")
	// fmt.Println(getDirSequence("<A"))

	complx_3 := 0

	for _, code := range num_codes {
		num := getNumFromCode(code)
		min_seq_len := getMinSeqLenAtDepth(code, 2)
		complx_3 += min_seq_len * num
	}

	// fmt.Println("total complexity 3:", complx_3)
	return complx_3
}

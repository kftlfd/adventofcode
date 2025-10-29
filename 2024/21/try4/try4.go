package try4

import (
	"math"
	"strconv"
)

const (
	UP    = 0
	RIGHT = 1
	DOWN  = 2
	LEFT  = 3
	A     = 4
)

var DIRKEYS = [5]string{"^", ">", "v", "<", "A"}

var DIRS = [4][2]int{
	{-1, 0},
	{0, 1},
	{1, 0},
	{0, -1},
}

var NUMPAD = [][]string{
	{"7", "8", "9"},
	{"4", "5", "6"},
	{"1", "2", "3"},
	{".", "0", "A"},
}

var DIRPAD = [][]string{
	{".", "^", "A"},
	{"<", "v", ">"},
}

type Paths map[string]map[string][][]string

var NUMPAD_PATHS = Paths{}

var DIRPAD_PATHS = Paths{}

func getPathsDfs(start, end string, r, c, end_r, end_c int, path []string, paths Paths, gap_r, gap_c int) {
	if r == gap_r && c == gap_c {
		return
	}
	if r == end_r && c == end_c {
		if _, ok := paths[start]; !ok {
			paths[start] = map[string][][]string{}
		}
		paths[start][end] = append(paths[start][end], append(path, "A"))
		return
	}
	if r > end_r {
		getPathsDfs(start, end, r-1, c, end_r, end_c, append(path, DIRKEYS[UP]), paths, gap_r, gap_c)
	}
	if r < end_r {
		getPathsDfs(start, end, r+1, c, end_r, end_c, append(path, DIRKEYS[DOWN]), paths, gap_r, gap_c)
	}
	if c > end_c {
		getPathsDfs(start, end, r, c-1, end_r, end_c, append(path, DIRKEYS[LEFT]), paths, gap_r, gap_c)
	}
	if c < end_c {
		getPathsDfs(start, end, r, c+1, end_r, end_c, append(path, DIRKEYS[RIGHT]), paths, gap_r, gap_c)
	}
}

func populatePaths(paths Paths, arr [][]string, gap_r, gap_c int) Paths {
	m := len(arr)
	n := len(arr[0])

	for start_r := 0; start_r < m; start_r++ {
		for start_c := 0; start_c < n; start_c++ {
			start := arr[start_r][start_c]
			for end_r := 0; end_r < m; end_r++ {
				for end_c := 0; end_c < n; end_c++ {
					end := arr[end_r][end_c]
					// if start == end {
					// 	continue
					// }
					getPathsDfs(start, end, start_r, start_c, end_r, end_c, []string{}, paths, gap_r, gap_c)
				}
			}
		}
	}
	return paths
}

//
//
//

func codeToDirPaths(code string) [][][]string {
	paths := [][][]string{}
	codeToDirPathsDfs(code, "A", 0, [][]string{}, &paths)
	return paths
}
func codeToDirPathsDfs(code, prev string, idx int, path [][]string, paths *[][][]string) {
	if idx >= len(code) {
		(*paths) = append((*paths), path)
		return
	}
	cur := prev
	nxt := code[idx : idx+1]
	for _, path_to_nxt := range NUMPAD_PATHS[cur][nxt] {
		codeToDirPathsDfs(code, nxt, idx+1, append(path, path_to_nxt), paths)
	}
}

func getSequencesForSegment(seg []string) [][][]string {
	sequences := [][][]string{}
	getSequencesForSegmentDfs(seg, "A", 0, [][]string{}, &sequences)
	return sequences
}
func getSequencesForSegmentDfs(seg []string, prev string, idx int, path [][]string, paths *[][][]string) {
	if idx >= len(seg) {
		(*paths) = append((*paths), path)
		return
	}
	cur := prev
	nxt := seg[idx]
	for _, path_to_nxt := range DIRPAD_PATHS[cur][nxt] {
		getSequencesForSegmentDfs(seg, nxt, idx+1, append(path, path_to_nxt), paths)
	}
}

func getNextSequences(seq [][]string) [][][]string {
	next_seqs := [][][]string{}
	getNextSequencesDfs(seq, 0, [][]string{}, &next_seqs)
	return next_seqs
}
func getNextSequencesDfs(seq [][]string, idx int, cur_seq [][]string, next_seqs *[][][]string) {
	if idx >= len(seq) {
		(*next_seqs) = append((*next_seqs), cur_seq)
		return
	}
	cur_segm := seq[idx]
	seqs_for_segm := getSequencesForSegment(cur_segm)
	for _, segm_seq := range seqs_for_segm {
		upd_cur_seq := append(cur_seq, segm_seq...)
		getNextSequencesDfs(seq, idx+1, upd_cur_seq, next_seqs)
	}
}

//
//
//

func getSequenceLen(seq [][]string, depth int) int {
	if depth == 0 {
		cur_len := 0
		for _, segm := range seq {
			cur_len += len(segm)
		}
		return cur_len
	}

	cur_len := math.MaxInt
	next_seqs := getNextSequences(seq)
	for _, nxt_seq := range next_seqs {
		cur_len = min(cur_len, getSequenceLen(nxt_seq, depth-1))
	}
	return cur_len
}

func getMinSeqLenForCode(code string, depth int) int {
	code_paths := codeToDirPaths(code)
	min_len := math.MaxInt
	for _, code_p := range code_paths {
		min_len = min(min_len, getSequenceLen(code_p, depth))
	}
	return min_len
}

//
//
//

func getCodeNum(code string) int {
	val, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		panic(err)
	}
	return val
}

//
//
//

func GetComplexity(num_codes []string, depth int) int {
	populatePaths(NUMPAD_PATHS, NUMPAD, 3, 0)
	populatePaths(DIRPAD_PATHS, DIRPAD, 0, 0)

	total_complexity := 0

	for _, code := range num_codes {
		num := getCodeNum(code)
		seq_len := getMinSeqLenForCode(code, 2)
		code_complexity := num * seq_len
		total_complexity += code_complexity
		// fmt.Println(code, seq_len, code_complexity)
	}

	return total_complexity
}

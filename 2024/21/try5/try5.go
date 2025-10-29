package try5

import (
	"math"
	"strconv"
)

// I gave up and looked for a tutorial
// thanks to u/ThatMakesMeM0ist
// https://www.reddit.com/r/adventofcode/comments/1hjx0x4/2024_day_21_quick_tutorial_to_solve_part_2_in
func GetComplexity(codes []string, depth int) int {
	total := 0

	cache := make(shortestSeqLenCache)

	for _, code := range codes {
		total += getCodeComplexity(code, depth, cache)
	}

	return total
}

func getCodeComplexity(code string, depth int, cache shortestSeqLenCache) int {
	numpadSeq := getSequencesForDirKeys(code, numpadMap)

	minSeqLen := math.MaxInt

	for _, seq := range numpadSeq {
		seqLen := getShortestSeqLen(seq, depth, cache)
		if seqLen < minSeqLen {
			minSeqLen = seqLen
		}
	}

	codeNum, err := strconv.Atoi(code[:len(code)-1])
	if err != nil {
		panic(err)
	}

	return minSeqLen * codeNum
}

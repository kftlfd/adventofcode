package try5

import (
	"math"
	"strings"
)

type shortestSeqLenCache = map[string]map[int]int

// recursively get get a len of shortest sequence of keys to press at level `depth` to get `keys`
func getShortestSeqLen(keys string, depth int, cache shortestSeqLenCache) int {
	if depth < 1 {
		return len(keys)
	}

	keysMap, keysMapExists := cache[keys]
	if keysMapExists {
		cacheVal, cacheValExists := keysMap[depth]
		if cacheValExists {
			return cacheVal
		}
	}

	total := 0

	subkeys := strings.SplitAfter(keys, "A")

	for _, subkey := range subkeys {
		sequences := getSequencesForDirKeys(subkey, dirpadMap)

		minSeqLen := math.MaxInt
		for _, seq := range sequences {
			seqLen := getShortestSeqLen(seq, depth-1, cache)
			if seqLen < minSeqLen {
				minSeqLen = seqLen
			}
		}

		total += minSeqLen
	}

	if !keysMapExists {
		cache[keys] = make(map[int]int)
		keysMap = cache[keys]
	}
	keysMap[depth] = total

	return total
}

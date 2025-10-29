package try5

import (
	"strings"
)

type keypadPathsMap = map[string]map[string][]string

// shortest paths between numpad keys
//
// eg. numpadMap["7"]["0"] = [">vvv", "v>vv", "vv>v"]
var numpadMap keypadPathsMap

func init() {
	numpad := [][]string{
		{"7", "8", "9"},
		{"4", "5", "6"},
		{"1", "2", "3"},
		{".", "0", "A"},
	}

	numpadKeys := []string{"A", "0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	numpadMap = createShortestPathsMap(numpad, numpadKeys)
}

// shortest paths between numpad keys
//
// eg. numpadMap["^"][">"] = ["v>", ">v"]
var dirpadMap keypadPathsMap

func init() {
	dirpad := [][]string{
		{".", "^", "A"},
		{"<", "v", ">"},
	}

	dirpadKeys := []string{"A", "<", ">", "^", "v"}

	dirpadMap = createShortestPathsMap(dirpad, dirpadKeys)
}

func createShortestPathsMap(keypad [][]string, keys []string) keypadPathsMap {
	keypadMap := make(keypadPathsMap)

	dirs := [4][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
	dirKeys := [4]string{">", "<", "v", "^"}

	keypadKeysCoords := make(map[string][2]int)

	for r := 0; r < len(keypad); r++ {
		for c := 0; c < len(keypad[r]); c++ {
			keypadKeysCoords[keypad[r][c]] = [2]int{r, c}
		}
	}

	containsZigZag := func(path []string) bool {
		if len(path) < 3 {
			return false
		}
		for i := 2; i < len(path); i++ {
			a := i - 2
			b := i - 1
			c := i
			if path[a] == path[c] && path[b] != path[a] {
				return true
			}
		}
		return false
	}

	getShortestPaths := func(start, end string) []string {
		paths := []string{}

		startIdx := keypadKeysCoords[start]
		startR := startIdx[0]
		startC := startIdx[1]

		endIdx := keypadKeysCoords[end]
		endR := endIdx[0]
		endC := endIdx[1]

		type QItem struct {
			r, c int
			path []string
		}
		q := []QItem{{r: startR, c: startC, path: []string{}}}
		seen := make(map[int]bool)
		seen[startR*10+startC] = true

		for len(q) > 0 {
			nxtQ := []QItem{}
			markSeen := []int{}
			reachedEnd := false

			for _, item := range q {
				if item.r == endR && item.c == endC {
					reachedEnd = true
					if !containsZigZag(item.path) {
						paths = append(paths, strings.Join(item.path, ""))
					}
					continue
				}
				for dirIdx := 0; dirIdx < 4; dirIdx++ {
					dirKey := dirKeys[dirIdx]
					dirR := dirs[dirIdx][0]
					dirC := dirs[dirIdx][1]
					nxtR := item.r + dirR
					nxtC := item.c + dirC
					if nxtR < 0 || nxtR >= len(keypad) || nxtC < 0 || nxtC >= len(keypad[nxtR]) {
						continue
					}
					if seen[nxtR*10+nxtC] {
						continue
					}
					if keypad[nxtR][nxtC] == "." {
						continue
					}
					nxtPath := []string{}
					nxtPath = append(nxtPath, item.path...)
					nxtPath = append(nxtPath, dirKey)
					nxtQ = append(nxtQ, QItem{r: nxtR, c: nxtC, path: nxtPath})
					markSeen = append(markSeen, nxtR*10+nxtC)
				}
			}

			for _, val := range markSeen {
				seen[val] = true
			}

			if reachedEnd {
				break
			}
			q = nxtQ
		}

		return paths
	}

	for start := 0; start < len(keys); start++ {
		for end := 0; end < len(keys); end++ {
			startKey := keys[start]
			endKey := keys[end]

			startMap, startMapExists := keypadMap[startKey]
			if !startMapExists {
				keypadMap[startKey] = make(map[string][]string)
				startMap = keypadMap[startKey]
			}

			startMap[endKey] = getShortestPaths(startKey, endKey)
		}
	}

	return keypadMap
}

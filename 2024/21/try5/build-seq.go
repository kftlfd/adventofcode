package try5

import "strings"

// for the given `keys` generate possible sequences of keys to press at one level up.
//
// eg. `<A` -> `['<v<A>>^A', '<v<A>^>A', 'v<<A>>^A', 'v<<A>^>A']`
func getSequencesForDirKeys(keys string, keypadMap keypadPathsMap) []string {
	result := []string{}
	buildSeq(keypadMap, strings.Split(keys, ""), 0, "A", []string{}, &result)
	return result
}

func buildSeq(keypadMap keypadPathsMap, keys []string, index int, prevKey string, currPath []string, result *[]string) {
	if index == len(keys) {
		*result = append(*result, strings.Join(currPath, ""))
		return
	}

	curKey := keys[index]
	paths := keypadMap[prevKey][curKey]

	for _, path := range paths {
		nxtPath := []string{}
		nxtPath = append(nxtPath, currPath...)
		nxtPath = append(nxtPath, path, "A")
		buildSeq(keypadMap, keys, index+1, curKey, nxtPath, result)
	}
}

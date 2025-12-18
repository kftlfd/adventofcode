package matrix

func getSliceSum(arr []int) (sum int) {
	for _, val := range arr {
		sum += val
	}
	return sum
}

func copySlice(arr []int) []int {
	cp := make([]int, len(arr))
	copy(cp, arr)
	return cp
}

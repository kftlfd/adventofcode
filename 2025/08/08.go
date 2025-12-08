package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

func parseInput() [][3]int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := [][3]int{}

	for scanner.Scan() {
		line := scanner.Text()
		coordsStr := strings.Split(line, ",")
		if len(coordsStr) != 3 {
			panic(fmt.Errorf("invalid input"))
		}
		x, err := strconv.Atoi(coordsStr[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(coordsStr[1])
		if err != nil {
			panic(err)
		}
		z, err := strconv.Atoi(coordsStr[2])
		if err != nil {
			panic(err)
		}
		inp = append(inp, [3]int{x, y, z})
	}

	return inp
}

type jbDistance struct {
	distance float64
	a, b     int
}

func getSortedDistances(junctionBoxes [][3]int) []jbDistance {
	n := len(junctionBoxes)
	distances := []jbDistance{}

	for a := 0; a < n-1; a++ {
		for b := a + 1; b < n; b++ {
			dist := math.Sqrt(
				math.Pow(float64(junctionBoxes[a][0]-junctionBoxes[b][0]), 2) +
					math.Pow(float64(junctionBoxes[a][1]-junctionBoxes[b][1]), 2) +
					math.Pow(float64(junctionBoxes[a][2]-junctionBoxes[b][2]), 2),
			)
			distances = append(distances, jbDistance{distance: dist, a: a, b: b})
		}
	}

	sort.Slice(distances, func(i, j int) bool {
		return distances[i].distance < distances[j].distance
	})

	return distances
}

func main() {
	jboxes := parseInput()
	n := len(jboxes)

	connectionsTarget := 0
	if len(os.Args) != 2 {
		fmt.Println("usage: go run 08.go (connectionsTarget)")
		os.Exit(1)
	}
	connectionsTarget, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("expected 'connectionsTarget' argument to be a valid integer")
		os.Exit(2)
	}

	distances := getSortedDistances(jboxes)

	// initially all junction boxes are in their own circuit
	circuits := make(map[int]int)
	for i := range n {
		circuits[i] = i
	}
	circuitSize := make([]int, n)
	for i := range n {
		circuitSize[i] = 1
	}

	// find
	getCircuitRoot := func(i int) int {
		root := circuits[i]
		for root != circuits[root] {
			root = circuits[root]
		}
		return root
	}

	// union
	connectCircuits := func(rootA, rootB int) int {
		if rootA == rootB {
			return circuitSize[rootA]
		}
		newRoot := min(rootA, rootB)
		removeRoot := max(rootA, rootB)
		circuits[rootA] = newRoot
		circuits[rootB] = newRoot
		circuitSize[newRoot] += circuitSize[removeRoot]
		circuitSize[removeRoot] = 0
		return circuitSize[newRoot]
	}

	connectionsMade := 0

	for _, dist := range distances {
		rootA := getCircuitRoot(dist.a)
		rootB := getCircuitRoot(dist.b)
		resultCircuitSize := connectCircuits(rootA, rootB)
		connectionsMade += 1

		if connectionsMade == connectionsTarget {
			sortedCircuitSizes := append([]int{}, circuitSize...)
			sort.Slice(sortedCircuitSizes, func(i, j int) bool {
				return sortedCircuitSizes[i] > sortedCircuitSizes[j]
			})

			top3CircuitSizesProduct := 1
			for i := range 3 {
				top3CircuitSizesProduct *= sortedCircuitSizes[i]
			}

			fmt.Println("Part 1:", top3CircuitSizesProduct)
		}

		if resultCircuitSize == n {
			part2Ans := jboxes[dist.a][0] * jboxes[dist.b][0]
			fmt.Println("Part 2:", part2Ans)
			break
		}
	}
}

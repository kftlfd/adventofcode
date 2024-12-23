package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func parseInput() Graph {
	file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	graph := Graph{}

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "-")
		c1 := line[0]
		c2 := line[1]

		graph[c1] = append(graph[c1], c2)
		graph[c2] = append(graph[c2], c1)
	}

	return graph
}

type Graph map[string][]string

func (g *Graph) areConnected(name1, name2 string) bool {
	for _, v := range (*g)[name1] {
		if v == name2 {
			return true
		}
	}
	return false
}

func (g *Graph) areInterconnected(name1, name2, name3 string) bool {
	return g.areConnected(name1, name2) && g.areConnected(name1, name3) && g.areConnected(name2, name3)
}

func beginsWithT(name string) bool {
	return name[:1] == "t"
}

type AlphArr []string

func (a AlphArr) Len() int {
	return len(a)
}

func (a AlphArr) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a AlphArr) Less(i, j int) bool {
	return a[i] < a[j]
}

func sortAlphabetically(names []string) []string {
	res := make([]string, len(names))
	for i, v := range names {
		res[i] = v
	}
	sort.Sort(AlphArr(res))
	return res
}

func formatAns(network []string) string {
	return strings.Join(sortAlphabetically(network), ",")
}

func main() {
	graph := parseInput()

	names := []string{}
	for k := range graph {
		names = append(names, k)
	}
	computers_count := len(names)

	// part 1
	gorups_of_3 := 0
	for c1 := 0; c1 < computers_count-2; c1++ {
		for c2 := c1 + 1; c2 < computers_count-1; c2++ {
			for c3 := c2 + 1; c3 < computers_count; c3++ {
				comp1 := names[c1]
				comp2 := names[c2]
				comp3 := names[c3]

				if !beginsWithT(comp1) && !beginsWithT(comp2) && !beginsWithT(comp3) {
					continue
				}

				if graph.areInterconnected(comp1, comp2, comp3) {
					gorups_of_3 += 1
				}
			}
		}
	}

	fmt.Println("groups of 3:", gorups_of_3)

	// part 2
	networks := [][]string{}
	for _, name := range names {
		for i, network := range networks {
			belongs_to_network := true
			for _, c := range network {
				if !graph.areConnected(name, c) {
					belongs_to_network = false
					break
				}
			}
			if belongs_to_network {
				networks[i] = append(networks[i], name)
			}
		}
		networks = append(networks, []string{name})
	}

	max_network_len := 0
	max_network_i := -1
	for i, network := range networks {
		if len(network) > max_network_len {
			max_network_len = len(network)
			max_network_i = i
		}
	}

	fmt.Printf("max network: %v\n", formatAns(networks[max_network_i]))
}

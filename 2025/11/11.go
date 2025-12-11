package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var (
	YOU = "you"
	OUT = "out"
	SVR = "svr"
	DAC = "dac"
	FFT = "fft"
)

type Node struct {
	label    string
	children map[string]*Node
}

type Graph map[string]*Node

func parseInput() *Graph {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	graph := make(Graph)

	getNode := func(label string) *Node {
		node, ok := graph[label]
		if ok {
			return node
		}
		node = &Node{label: label, children: make(map[string]*Node)}
		graph[label] = node
		return node
	}

	for scanner.Scan() {
		line := strings.Split(scanner.Text(), ":")

		if len(line) != 2 {
			panic(fmt.Errorf("invalid input"))
		}

		cur := getNode(line[0])

		for _, nxtLabel := range strings.Fields(line[1]) {
			nxtNode := getNode(nxtLabel)
			cur.children[nxtLabel] = nxtNode
		}
	}

	return &graph
}

func dfs(node *Node, endLabel, omit string, memo *map[string]int) int {
	if node.label == omit {
		return 0
	}
	if count, visited := (*memo)[node.label]; visited {
		return count
	}
	if node.label == endLabel {
		return 1
	}
	paths := 0
	for nxtLabel := range node.children {
		paths += dfs(node.children[nxtLabel], endLabel, omit, memo)
	}
	(*memo)[node.label] = paths
	return paths
}

func (g *Graph) getPathsBetween(startLabel, endLabel, omit string) (int, error) {
	root, ok := (*g)[startLabel]
	if !ok {
		return 0, fmt.Errorf("no root node")
	}
	_, ok = (*g)[endLabel]
	if !ok {
		return 0, fmt.Errorf("no end node")
	}
	memo := make(map[string]int)
	paths := dfs(root, endLabel, omit, &memo)
	return paths, nil
}

func solvePart1(g *Graph) {
	paths, err := g.getPathsBetween(YOU, OUT, "")
	errMsg := ""
	if err != nil {
		errMsg = err.Error()
	}
	fmt.Println("Part 1:", paths, errMsg)
}

func solvePart2(g *Graph) {
	svrToDac, _ := g.getPathsBetween(SVR, DAC, FFT)
	svrToFft, _ := g.getPathsBetween(SVR, FFT, DAC)
	dacToFft, _ := g.getPathsBetween(DAC, FFT, "")
	fftToDac, _ := g.getPathsBetween(FFT, DAC, "")
	dacToOut, _ := g.getPathsBetween(DAC, OUT, FFT)
	fftToOut, _ := g.getPathsBetween(FFT, OUT, DAC)
	total := (svrToDac * dacToFft * fftToOut) + (svrToFft * fftToDac * dacToOut)
	fmt.Println("Part 2:", total)
}

func main() {
	graph := parseInput()
	solvePart1(graph)
	solvePart2(graph)
}

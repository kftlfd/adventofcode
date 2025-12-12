package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type box struct {
	m, n   int
	shape  [][]string
	filled int
}

type region struct {
	m, n  int
	boxes []int
}

func parseInput() ([]box, []region) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	sections := [][]string{}
	curSection := []string{}

	for scanner.Scan() {
		line := scanner.Text()

		if len(line) < 1 {
			sections = append(sections, curSection)
			curSection = []string{}
			continue
		}

		curSection = append(curSection, line)
	}
	sections = append(sections, curSection)

	boxes := []box{}
	for _, section := range sections[:len(sections)-1] {
		boxes = append(boxes, parseBox(section))
	}

	regions := parseRegions(sections[len(sections)-1])

	return boxes, regions
}

func parseBox(section []string) box {
	m := len(section) - 1
	n := len(section[1])
	filled := 0
	shape := [][]string{}
	for _, line := range section[1:] {
		vals := strings.Split(line, "")
		shape = append(shape, vals)
		for _, v := range vals {
			if v == "#" {
				filled += 1
			}
		}
	}
	return box{m: m, n: n, shape: shape, filled: filled}
}

func parseRegions(section []string) []region {
	regions := []region{}
	for _, line := range section {
		fields := strings.Fields(line)
		dimensionStr := strings.Split(fields[0][:len(fields[0])-1], "x")
		n, err := strconv.Atoi(dimensionStr[0])
		if err != nil {
			panic(err)
		}
		m, err := strconv.Atoi(dimensionStr[1])
		if err != nil {
			panic(err)
		}
		boxes := []int{}
		for _, boxStr := range fields[1:] {
			boxNum, err := strconv.Atoi(boxStr)
			if err != nil {
				panic(err)
			}
			boxes = append(boxes, boxNum)
		}
		regions = append(regions, region{m: m, n: n, boxes: boxes})
	}
	return regions
}

func (r region) isPossibleByTotalBoxesFill(boxes []box) bool {
	totalBoxesFill := 0
	for b, count := range r.boxes {
		totalBoxesFill += boxes[b].filled * count
	}
	return totalBoxesFill <= r.m*r.n
}

func (r region) canFitBoxesByDimensions(boxM, boxN int) bool {
	totalSpaces := (r.m / boxM) * (r.n / boxN)
	boxesNeeded := 0
	for _, count := range r.boxes {
		boxesNeeded += count
	}
	return totalSpaces >= boxesNeeded
}

func main() {
	boxes, regions := parseInput()

	boxM := boxes[0].m
	boxN := boxes[0].n
	for _, b := range boxes[1:] {
		if b.m != boxM || b.n != boxN {
			fmt.Println("boxes have different dimensions, can't proceed")
			os.Exit(0)
		}
	}

	maybePossibleByTotalFill := 0
	definetelyPossibleByDimensions := 0

	for _, r := range regions {
		if r.isPossibleByTotalBoxesFill(boxes) {
			maybePossibleByTotalFill += 1
			if r.canFitBoxesByDimensions(boxM, boxN) {
				definetelyPossibleByDimensions += 1
			}
		}
	}

	fmt.Println("total regions:                ", len(regions))
	fmt.Println("maybe possible by total fill: ", maybePossibleByTotalFill)
	fmt.Println("can fit just by dimensions:   ", definetelyPossibleByDimensions)
}

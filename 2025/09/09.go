package main

import (
	"aoc2025-09/part2"
	"bufio"
	"fmt"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/bits-and-blooms/bitset"
)

func parseInput() [][2]int {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanLines)

	inp := [][2]int{}

	for scanner.Scan() {
		line := scanner.Text()
		numStr := strings.Split(line, ",")
		if len(numStr) != 2 {
			panic(fmt.Errorf("invalid input"))
		}
		x, err := strconv.Atoi(numStr[0])
		if err != nil {
			panic(err)
		}
		y, err := strconv.Atoi(numStr[1])
		if err != nil {
			panic(err)
		}
		inp = append(inp, [2]int{x, y})
	}

	return inp
}

func main() {
	redSquares := parseInput()

	solvePart1(redSquares)

	p2 := part2.Solve(redSquares)
	if p2 > 0 {
		return
	}

	p2 = solvePart2Compact(redSquares)
	if p2 > 0 {
		return
	}

	solvePart2(redSquares)
}

func solvePart1(redSquares [][2]int) {
	n := len(redSquares)
	maxArea := 0
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			a := redSquares[i][0] - redSquares[j][0]
			if a < 0 {
				a *= -1
			}
			b := redSquares[i][1] - redSquares[j][1]
			if b < 0 {
				b *= -1
			}
			maxArea = max(maxArea, (a+1)*(b+1))
		}
	}
	fmt.Println("Part 1:", maxArea)
}

func solvePart2(redSquares [][2]int) {
	n := len(redSquares)

	// move points closer to (0,0), so that minX = 0, minY = 0
	minX := math.MaxInt
	maxX := -math.MaxInt
	minY := math.MaxInt
	maxY := -math.MaxInt
	for _, point := range redSquares {
		x, y := point[0], point[1]
		minX = min(minX, x)
		maxX = max(maxX, x)
		minY = min(minY, y)
		maxY = max(maxY, y)
	}
	// fmt.Printf("(%d,%d) (%d,%d)\n", minX, minY, maxX, maxY)
	// fmt.Printf("(0,0) (%d,%d)\n", maxX-minX, maxY-minY)
	// fmt.Println((maxX - minX) * (maxY - minY))
	cols := maxX - minX + 2
	rows := maxY - minY + 2
	points := [][2]int{}
	for _, ogPoint := range redSquares {
		points = append(points, [2]int{ogPoint[0] - minX, ogPoint[1] - minY})
	}

	// prepare a bit grid
	grid := make([]bitset.BitSet, rows)
	for r := range rows {
		grid[r] = *bitset.New(uint(cols))
	}

	// draw outline
	curX, curY := points[len(points)-1][0], points[len(points)-1][1]
	for _, point := range points {
		nxtX, nxtY := point[0], point[1]

		// startX, startY := curX, curY

		if nxtX == curX {
			move := 1
			if nxtY < curY {
				move = -1
			}
			for curY != nxtY {
				curY += move
				// fmt.Printf("(%d,%d)->(%d,%d) => +%d (%d,%d)\n", startX, startY, nxtX, nxtY, move, curX, curY)
				if grid[curY].Test(uint(curX)) {
					panic(fmt.Errorf("setting the same bit twice"))
				}
				grid[curY].Set(uint(curX))
			}
		} else if nxtY == curY {
			move := 1
			if nxtX < curX {
				move = -1
			}
			for curX != nxtX {
				curX += move
				if grid[curY].Test(uint(curX)) {
					panic(fmt.Errorf("setting the same bit twice"))
				}
				grid[curY].Set(uint(curX))
			}
		} else {
			panic(fmt.Errorf("next point not on the same line as prev"))
		}
	}

	// find the starting point to paint the inside
	startX, startY := 0, 0
	for !grid[startY].Test(uint(startX)) {
		startX += 1
	}
	for grid[startY+1].Test(uint(startX)) {
		startX += 1
	}
	startY += 1
	// print the start point for sanity check
	// for y := 0; y < 5; y++ {
	// 	for x := startX - 5; x < startX+5; x++ {
	// 		val := "."
	// 		if grid[y].Test(uint(x)) {
	// 			val = "#"
	// 		}
	// 		fmt.Printf("%v", val)
	// 	}
	// 	fmt.Printf("\n")
	// }

	// paint the inside of the polygon
	grid[startY].Set(uint(startX))
	q := [][2]int{{startX, startY}}
	moves := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	// painted := 0
	for len(q) > 0 {
		nxtq := [][2]int{}
		for _, cur := range q {

			curX, curY := cur[0], cur[1]
			for _, d := range moves {
				nxtX, nxtY := curX+d[0], curY+d[1]
				if !grid[nxtY].Test(uint(nxtX)) {
					grid[nxtY].Set(uint(nxtX))
					// painted += 1
					nxtq = append(nxtq, [2]int{nxtX, nxtY})
				}
			}

		}
		q = nxtq
	}
	// fmt.Println("painted", painted)

	// go through each possible rectangle, check if it's inside polygon
	// brute-force, but with go **concurrency**
	n = len(points)
	maxArea := 0

	getAreaIfValid := func(i, j int) int {
		iX, iY := points[i][0], points[i][1]
		jX, jY := points[j][0], points[j][1]

		// check if all points of the rectangle are inside the polygon
		startY := min(iY, jY)
		endY := max(iY, jY)
		startX := min(iX, jX)
		endX := max(iX, jX)

		ok := func() bool {
			// snail path: right, down, left, up, one in
			for i := 0; i < max(endY-startY, endX-startX); i++ {
				sY, sX := startY+i, startX+i
				eY, eX := endY-i, endX-i
				if sY > eY || sX > eX {
					return true
				}
				for x := sX; x < eX; x++ {
					if !grid[sY].Test(uint(x)) {
						return false
					}
				}
				for y := sY; y < eY; y++ {
					if !grid[y].Test(uint(eX)) {
						return false
					}
				}
				for x := eX; x > sX; x-- {
					if !grid[eY].Test(uint(x)) {
						return false
					}
				}
				for y := eY; y > sY; y-- {
					if !grid[y].Test(uint(sX)) {
						return false
					}
				}
			}

			// left to right, top to bottom
			// for y := startY; y <= endY; y++ {
			// 	for x := startX; x <= endX; x++ {
			// 		if !grid[y].Test(uint(x)) {
			// 			return false
			// 		}
			// 	}
			// }

			return true
		}()

		if !ok {
			return -1
		}

		return (endY - startY + 1) * (endX - startX + 1)
	}

	in := make(chan [2]int)
	out := make(chan int)

	var wgWorkers sync.WaitGroup
	var wgResults sync.WaitGroup

	for range 20 {
		wgWorkers.Go(func() {
			for pointsIn := range in {
				area := getAreaIfValid(pointsIn[0], pointsIn[1])
				out <- area
			}
		})
	}

	wgResults.Go(func() {
		cur := 0
		// totalRects := (n) * (n - 1) / 2
		for area := range out {
			cur += 1
			// fmt.Printf("%d / %d: ", cur, totalRects)
			if area == -1 {
				// fmt.Printf("X\n")
			} else {
				// fmt.Printf("%d\n", area)
				maxArea = max(maxArea, area)
			}
		}
	})

	go func() {
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				in <- [2]int{i, j}
			}
		}
		close(in)
	}()

	wgWorkers.Wait()
	close(out)
	wgResults.Wait()

	fmt.Println("Part 2:", maxArea)
}

func solvePart2Compact(input [][2]int) int {
	// make compact grid, compact points by removing unused x,y coordinates between points
	xVals := []int{}
	yVals := []int{}
	seenX := map[int]bool{}
	seenY := map[int]bool{}
	for _, point := range input {
		_, seenCurX := seenX[point[0]]
		if !seenCurX {
			xVals = append(xVals, point[0])
			seenX[point[0]] = true
		}
		_, seenCurY := seenY[point[1]]
		if !seenCurY {
			yVals = append(yVals, point[1])
			seenY[point[1]] = true
		}
	}
	slices.Sort(xVals)
	slices.Sort(yVals)
	newXCoord := map[int]int{}
	newX := 1
	for _, oldX := range xVals {
		newXCoord[oldX] = newX
		newX += 2
	}
	maxX := newX + 1
	newYCoord := map[int]int{}
	newY := 1
	for _, oldY := range yVals {
		newYCoord[oldY] = newY
		newY += 2
	}
	maxY := newY + 1
	compactPoints := [][2]int{}
	for _, oldPoint := range input {
		oldX, oldY := oldPoint[0], oldPoint[1]
		newX, newY := newXCoord[oldX], newYCoord[oldY]
		compactPoints = append(compactPoints, [2]int{newX, newY})
	}
	compactGrid := make([][]bool, maxY+1)
	for r := range maxY + 1 {
		compactGrid[r] = make([]bool, maxX+1)
	}

	// draw polygon outline
	curX, curY := compactPoints[len(compactPoints)-1][0], compactPoints[len(compactPoints)-1][1]
	for _, point := range compactPoints {
		nxtX, nxtY := point[0], point[1]
		// startX, startY := curX, curY
		if nxtX == curX {
			move := 1
			if nxtY < curY {
				move = -1
			}
			for curY != nxtY {
				curY += move
				// fmt.Printf("(%d,%d)->(%d,%d) => +%d (%d,%d)\n", startX, startY, nxtX, nxtY, move, curX, curY)
				if compactGrid[curY][curX] {
					panic(fmt.Errorf("setting the same bit twice"))
				}
				compactGrid[curY][curX] = true
			}
		} else if nxtY == curY {
			move := 1
			if nxtX < curX {
				move = -1
			}
			for curX != nxtX {
				curX += move
				if compactGrid[curY][curX] {
					panic(fmt.Errorf("setting the same bit twice"))
				}
				compactGrid[curY][curX] = true
			}
		} else {
			panic(fmt.Errorf("next point not on the same line as prev"))
		}
	}

	// find the starting point to paint the inside
	startX, startY := 0, 1
	for !compactGrid[startY][startX] {
		startX += 1
	}
	for compactGrid[startY+1][startX] {
		startX += 1
	}
	startY += 1

	// paint the inside of the polygon
	compactGrid[startY][startX] = true
	q := [][2]int{{startX, startY}}
	moves := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}
	for len(q) > 0 {
		nxtq := [][2]int{}
		for _, cur := range q {

			curX, curY := cur[0], cur[1]
			for _, d := range moves {
				nxtX, nxtY := curX+d[0], curY+d[1]
				if !compactGrid[nxtY][nxtX] {
					compactGrid[nxtY][nxtX] = true
					nxtq = append(nxtq, [2]int{nxtX, nxtY})
				}
			}

		}
		q = nxtq
	}

	// go through each possible rectangle, check if it's inside polygon
	n := len(compactPoints)
	maxArea := 0

	compactSquareIsInsideCompactPolygon := func(i, j int) bool {
		iX, iY := compactPoints[i][0], compactPoints[i][1]
		jX, jY := compactPoints[j][0], compactPoints[j][1]

		// check if all points of the rectangle are inside the polygon
		startY := min(iY, jY)
		endY := max(iY, jY)
		startX := min(iX, jX)
		endX := max(iX, jX)

		// snail path: right, down, left, up, one in
		for i := 0; i < max(endY-startY, endX-startX); i++ {
			sY, sX := startY+i, startX+i
			eY, eX := endY-i, endX-i
			if sY > eY || sX > eX {
				return true
			}
			for x := sX; x < eX; x++ {
				if !compactGrid[sY][x] {
					return false
				}
			}
			for y := sY; y < eY; y++ {
				if !compactGrid[y][eX] {
					return false
				}
			}
			for x := eX; x > sX; x-- {
				if !compactGrid[eY][x] {
					return false
				}
			}
			for y := eY; y > sY; y-- {
				if !compactGrid[y][sX] {
					return false
				}
			}
		}
		return true
	}

	getSquareArea := func(i, j int) int {
		x1, y1 := input[i][0], input[i][1]
		x2, y2 := input[j][0], input[j][1]
		yLen := max(y1, y2) - min(y1, y2)
		xLen := max(x1, x2) - min(x1, x2)
		return (yLen + 1) * (xLen + 1)
	}

	in := make(chan [2]int)
	out := make(chan int)

	var wgWorkers sync.WaitGroup
	var wgResults sync.WaitGroup

	for range 20 {
		wgWorkers.Go(func() {
			for pointsIn := range in {
				ok := compactSquareIsInsideCompactPolygon(pointsIn[0], pointsIn[1])
				area := -1
				if ok {
					area = getSquareArea(pointsIn[0], pointsIn[1])
				}
				out <- area
			}
		})
	}

	wgResults.Go(func() {
		cur := 0
		// totalRects := (n) * (n - 1) / 2
		for area := range out {
			cur += 1
			// fmt.Printf("%d / %d: ", cur, totalRects)
			if area == -1 {
				// fmt.Printf("X\n")
			} else {
				// fmt.Printf("%d\n", area)
				maxArea = max(maxArea, area)
			}
		}
	})

	go func() {
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				in <- [2]int{i, j}
			}
		}
		close(in)
	}()

	wgWorkers.Wait()
	close(out)
	wgResults.Wait()

	fmt.Println("Part 2:", maxArea)

	return maxArea
}

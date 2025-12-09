package part2

import "slices"

type point struct {
	x, y int
}

func toPoints(input [][2]int) []point {
	points := []point{}
	for _, p := range input {
		points = append(points, point{x: p[0], y: p[1]})
	}
	return points
}

type compactPoints struct {
	points []point
	m, n   int
}

func toCompactPoints(points []point) compactPoints {
	// get all used unique values of X/Y coordinates
	xVals := []int{}
	yVals := []int{}
	seenX := map[int]bool{}
	seenY := map[int]bool{}
	for _, point := range points {
		_, seenCurX := seenX[point.x]
		if !seenCurX {
			xVals = append(xVals, point.x)
			seenX[point.x] = true
		}
		_, seenCurY := seenY[point.y]
		if !seenCurY {
			yVals = append(yVals, point.y)
			seenY[point.y] = true
		}
	}

	// sort in increasing order
	slices.Sort(xVals)
	slices.Sort(yVals)

	// map original X values to compact values
	newXCoord := map[int]int{}
	newX := 1
	for _, oldX := range xVals {
		newXCoord[oldX] = newX
		newX += 2
	}

	// map original Y values to compact values
	newYCoord := map[int]int{}
	newY := 1
	for _, oldY := range yVals {
		newYCoord[oldY] = newY
		newY += 2
	}

	compPoints := []point{}
	for _, oldPoint := range points {
		oldX, oldY := oldPoint.x, oldPoint.y
		newX, newY := newXCoord[oldX], newYCoord[oldY]
		compPoints = append(compPoints, point{x: newX, y: newY})
	}

	return compactPoints{points: compPoints, m: newY, n: newX}
}

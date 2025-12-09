package part2

import (
	"fmt"
	"sync"
)

func getSquareArea(a, b point) int {
	yLen := max(a.y, b.y) - min(a.y, b.y)
	xLen := max(a.x, b.x) - min(a.x, b.x)
	return (yLen + 1) * (xLen + 1)
}

func Solve(input [][2]int) int {
	n := len(input)
	points := toPoints(input)
	compactPoints := toCompactPoints(points)
	compactPolygon := newCompactPolygon(compactPoints)
	maxArea := 0

	for i := range n - 1 {
		for j := i + 1; j < n; j++ {
			// check if compact square is inside polygon in the compact space
			// if yes, then original non-compacted square also is inside non-compacted polygon
			compA, compB := compactPoints.points[i], compactPoints.points[j]
			ok := compactPolygon.squareIsInside(compA, compB)
			if ok {
				// get area of original square, non-compacted
				pointA, pointB := points[i], points[j]
				area := getSquareArea(pointA, pointB)
				maxArea = max(maxArea, area)
			}
		}
	}

	fmt.Println("Part 2:", maxArea)
	return maxArea
}

func SolveConcurrently(input [][2]int) int {
	n := len(input)
	points := toPoints(input)
	compactPoints := toCompactPoints(points)
	compactPolygon := newCompactPolygon(compactPoints)
	maxArea := 0

	in := make(chan [2]int)
	out := make(chan int)

	var wgIn sync.WaitGroup
	var wgOut sync.WaitGroup

	for range 10 {
		wgIn.Go(func() {
			for p := range in {
				i, j := p[0], p[1]
				// check if compact square is inside polygon in the compact space
				// if yes, then original non-compacted square also is inside non-compacted polygon
				compA, compB := compactPoints.points[i], compactPoints.points[j]
				ok := compactPolygon.squareIsInside(compA, compB)
				if ok {
					// get area of original square, non-compacted
					pointA, pointB := points[i], points[j]
					out <- getSquareArea(pointA, pointB)
				}
			}
		})
	}

	wgOut.Go(func() {
		for area := range out {
			maxArea = max(maxArea, area)
		}
	})

	go func() {
		for i := range n - 1 {
			for j := i + 1; j < n; j++ {
				// send pair of points indexes in processing
				in <- [2]int{i, j}
			}
		}
		close(in)
	}()

	wgIn.Wait()
	close(out)
	wgOut.Wait()

	fmt.Println("Part 2:", maxArea)

	return maxArea
}

package part2

import "fmt"

type compactPolygon struct {
	grid [][]bool
	m, n int
}

func newCompactPolygon(compPoints compactPoints) compactPolygon {
	m, n := compPoints.m, compPoints.n
	grid := make([][]bool, m)
	for r := range m {
		grid[r] = make([]bool, n)
	}
	compPolygon := compactPolygon{grid: grid, m: m, n: n}

	compPolygon.drawPolygon(compPoints.points)

	return compPolygon
}

func (g *compactPolygon) drawPolygon(points []point) {
	err := g.drawPolygonOutline(points)
	if err != nil {
		panic(err)
	}
	fillStart, err := g.findFillStartingPoint()
	if err != nil {
		panic(err)
	}
	err = g.fillPolygon(fillStart)
	if err != nil {
		panic(err)
	}
}

func (g *compactPolygon) drawPolygonOutline(points []point) error {
	last := len(points) - 1

	curX, curY := points[last].x, points[last].y

	for _, point := range points {
		nxtX, nxtY := point.x, point.y
		if nxtX == curX {
			move := 1
			if nxtY < curY {
				move = -1
			}
			for curY != nxtY {
				curY += move
				if (*g).grid[curY][curX] {
					return fmt.Errorf("setting the same bit twice")
				}
				(*g).grid[curY][curX] = true
			}
		} else if nxtY == curY {
			move := 1
			if nxtX < curX {
				move = -1
			}
			for curX != nxtX {
				curX += move
				if (*g).grid[curY][curX] {
					return fmt.Errorf("setting the same bit twice")
				}
				(*g).grid[curY][curX] = true
			}
		} else {
			return fmt.Errorf("next point not on the same line as prev")
		}
	}

	return nil
}

func (g *compactPolygon) findFillStartingPoint() (point, error) {
	m, n := (*g).m, (*g).n

	for y := range m - 1 {
		for x := range n {
			if (*g).grid[y][x] && !(*g).grid[y+1][x] {
				return point{x: x, y: y + 1}, nil
			}
		}
	}

	return point{}, fmt.Errorf("starting point not found")
}

func (g *compactPolygon) fillPolygon(start point) error {
	moves := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	(*g).grid[start.y][start.x] = true

	q := []point{start}

	for len(q) > 0 {
		nxtq := []point{}
		for _, cur := range q {
			for _, d := range moves {
				nxtX, nxtY := cur.x+d[0], cur.y+d[1]
				if !(*g).grid[nxtY][nxtX] {
					(*g).grid[nxtY][nxtX] = true
					nxtq = append(nxtq, point{x: nxtX, y: nxtY})
				}
			}
		}
		q = nxtq
	}

	return nil
}

func (g *compactPolygon) squareIsInside(a, b point) bool {
	m, n := (*g).m, (*g).n

	if a.x < 0 || b.x < 0 || a.x >= n || b.x >= n {
		return false
	}
	if a.y < 0 || b.y < 0 || a.y >= m || b.y >= m {
		return false
	}

	startY := min(a.y, b.y)
	endY := max(a.y, b.y)
	startX := min(a.x, b.x)
	endX := max(a.x, b.x)

	// snail path: right, down, left, up, one in
	maxI := max(endY-startY, endX-startX)
	for i := range maxI {
		sY, sX := startY+i, startX+i
		eY, eX := endY-i, endX-i
		if sY > eY || sX > eX {
			return true
		}
		for x := sX; x < eX; x++ {
			if !(*g).grid[sY][x] {
				return false
			}
		}
		for y := sY; y < eY; y++ {
			if !(*g).grid[y][eX] {
				return false
			}
		}
		for x := eX; x > sX; x-- {
			if !(*g).grid[eY][x] {
				return false
			}
		}
		for y := eY; y > sY; y-- {
			if !(*g).grid[y][sX] {
				return false
			}
		}
	}

	return true
}

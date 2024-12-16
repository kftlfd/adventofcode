package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"math"
	"os"
	"strings"
)

const (
	R     = 0
	C     = 1
	NORTH = 0
	EAST  = 1
	SOUTH = 2
	WEST  = 3
)

var DIRS = [4][2]int{
	{-1, 0}, // North
	{0, 1},  // East
	{1, 0},  // South
	{0, -1}, // West
}

func isOppositeDir(cur, nxt int) bool {
	if cur < 0 || cur > 3 || nxt < 0 || nxt > 3 {
		panic(fmt.Errorf("isOppositeDir: bad input: want values in range [0,3], got: %v %v", cur, nxt))
	}
	return nxt == (cur+2)%4
}

type Maze struct {
	grid       [][]string
	m, n       int
	start, end [2]int
}

func parseInput(file *os.File) Maze {
	reader := bufio.NewScanner(file)
	reader.Split(bufio.ScanLines)

	grid := [][]string{}

	for reader.Scan() {
		grid = append(grid, strings.Split(reader.Text(), ""))
	}

	m := len(grid)
	n := len(grid[0])

	maze := Maze{grid: grid, m: m, n: n}

	for r := 0; r < m; r++ {
		for c := 0; c < n; c++ {
			if grid[r][c] == "S" {
				maze.start = [2]int{r, c}
			}
			if grid[r][c] == "E" {
				maze.end = [2]int{r, c}
			}
		}
	}

	return maze
}

func main() {
	inp_file, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}
	defer inp_file.Close()

	maze := parseInput(inp_file)

	min_score, err := maze.getMinScore()
	if err != nil {
		panic(err)
	}
	fmt.Println("min score:", min_score)

	// score, tiles := maze.getMinScoreAndPathsTileCount()
	score, tiles := maze.getMinScoreAndPathsTileCount()
	fmt.Printf("min score: %v, tiles: %v\n", score, tiles)

}

//
// part 1
//

type DeerPos struct {
	score, r, c, dir, dist int
	path                   [][2]int
}

/*
=======================================
Priority Queue
https://www.slingacademy.com/article/exploring-priority-queues-using-the-container-heap-package-in-go/
*/

type PriorityQueue []*DeerPos

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].score < pq[j].score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*DeerPos))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil // avoid memory leak
	*pq = old[0 : n-1]
	return item
}

/*
=======================================
*/

type DeerPosPQ struct {
	q *PriorityQueue
}

func newDeerPosPQ() DeerPosPQ {
	return DeerPosPQ{q: &PriorityQueue{}}
}

func (pq DeerPosPQ) Len() int {
	return len(*pq.q)
}

func (pq DeerPosPQ) push(d *DeerPos) {
	heap.Push(pq.q, d)
}

func (pq DeerPosPQ) pop() *DeerPos {
	return heap.Pop(pq.q).(*DeerPos)
}

func (m Maze) getMinScore() (int, error) {
	sr := m.start[R]
	sc := m.start[C]

	q := newDeerPosPQ()
	q.push(&DeerPos{score: 0, r: sr, c: sc, dir: EAST})

	scores := [][]int{}
	for r := 0; r < m.m; r++ {
		row := make([]int, m.n)
		for c := 0; c < m.n; c++ {
			row[c] = math.MaxInt
		}
		scores = append(scores, row)
	}
	scores[sr][sc] = 0

	for q.Len() > 0 {
		cur := q.pop()

		if cur.r == m.end[R] && cur.c == m.end[C] {
			return cur.score, nil
		}

		for i, dir := range DIRS {
			nr := cur.r + dir[R]
			nc := cur.c + dir[C]

			nscore := cur.score + 1
			if isOppositeDir(cur.dir, i) {
				nscore += 2000
			} else if cur.dir != i {
				nscore += 1000
			}

			if m.grid[nr][nc] == "#" || scores[nr][nc] <= nscore {
				continue
			}

			scores[nr][nc] = nscore

			q.push(&DeerPos{score: nscore, r: nr, c: nc, dir: i, dist: cur.dist + 1})
		}
	}

	return 0, fmt.Errorf("min path not found")
}

//
// part 2
//

type ScoreFromDir [4]int

func newScoreFromDir() ScoreFromDir {
	return ScoreFromDir{math.MaxInt, math.MaxInt, math.MaxInt, math.MaxInt}
}

type FromNodes [4][][3]int

func newFromNodes() FromNodes {
	return FromNodes{make([][3]int, 0), make([][3]int, 0), make([][3]int, 0), make([][3]int, 0)}
}

/* 3D Dijkstra + record from which node (RxCxDirection) coming to the next node, then walk back from the end */
func (m Maze) getMinScoreAndPathsTileCount() (int, int) {
	sr := m.start[R]
	sc := m.start[C]

	q := newDeerPosPQ()
	q.push(&DeerPos{score: 0, r: sr, c: sc, dir: EAST, dist: 0, path: [][2]int{{sr, sc}}})

	scores := [][]ScoreFromDir{}
	for r := 0; r < m.m; r++ {
		row := make([]ScoreFromDir, m.n)
		for c := 0; c < m.n; c++ {
			row[c] = newScoreFromDir()
		}
		scores = append(scores, row)
	}
	scores[sr][sc] = ScoreFromDir{0, 0, 0, 0}

	from := [][]FromNodes{}
	for r := 0; r < m.m; r++ {
		row := make([]FromNodes, m.n)
		for c := 0; c < m.n; c++ {
			row[c] = newFromNodes()
		}
		from = append(from, row)
	}

	min_score := math.MaxInt

	for q.Len() > 0 {
		cur := q.pop()

		if cur.score > min_score {
			break
		}

		if cur.r == m.end[R] && cur.c == m.end[C] {
			min_score = min(min_score, cur.score)
			continue
		}

		for ndir, ndir_d := range DIRS {
			nr := cur.r + ndir_d[R]
			nc := cur.c + ndir_d[C]

			if m.grid[nr][nc] == "#" {
				continue
			}

			nscore := cur.score + 1
			if isOppositeDir(cur.dir, ndir) {
				nscore += 2000
			} else if cur.dir != ndir {
				nscore += 1000
			}

			if scores[nr][nc][ndir] < nscore {
				continue
			}
			if scores[nr][nc][ndir] == nscore {
				from[nr][nc][ndir] = append(from[nr][nc][ndir], [3]int{cur.r, cur.c, cur.dir})
				continue
			}

			scores[nr][nc][ndir] = nscore
			from[nr][nc][ndir] = append(from[nr][nc][ndir], [3]int{cur.r, cur.c, cur.dir})

			q.push(&DeerPos{score: nscore, r: nr, c: nc, dir: ndir})
		}
	}

	// fmt.Println("min score:", min_score)

	tiles := [][]int{}
	for r := 0; r < m.m; r++ {
		tiles = append(tiles, make([]int, m.n))
	}

	tq := [][3]int{
		{m.end[R], m.end[C], NORTH},
		{m.end[R], m.end[C], SOUTH},
		{m.end[R], m.end[C], EAST},
		{m.end[R], m.end[C], WEST},
	}

	visited := [][][4]bool{}
	for r := 0; r < m.m; r++ {
		visited = append(visited, make([][4]bool, m.n))
	}
	visited[m.end[R]][m.end[C]] = [4]bool{true, true, true, true}

	for len(tq) > 0 {
		cur := tq[0]
		tq = tq[1:]

		r := cur[R]
		c := cur[C]
		d := cur[2]

		tiles[r][c] = 1

		froms := from[r][c][d]
		for _, nxt := range froms {
			nr := nxt[R]
			nc := nxt[C]
			nd := nxt[2]

			if visited[nr][nc][nd] {
				continue
			}

			visited[nr][nc][nd] = true
			tq = append(tq, [3]int{nr, nc, nd})
		}
	}

	tiles_count := 0
	for _, row := range tiles {
		for _, val := range row {
			tiles_count += val
		}
	}

	// fmt.Println("tiles count:", tiles_count)

	return min_score, tiles_count
}

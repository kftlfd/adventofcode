package flags

import "flag"

var (
	Time        = flag.Bool("t", false, "show running time")
	Debug       = flag.Bool("d", false, "show more info")
	P2SimpleDFS = flag.Bool("p2dfs", false, "use simple DFS for Part 2 (super slow)")
)

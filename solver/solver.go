package solver

import "projekt3/graph"

type ATSPSolver interface {
	SetGraph(graph graph.Graph)
	GetGraph() graph.Graph
	Solve() ([]int, int)
}

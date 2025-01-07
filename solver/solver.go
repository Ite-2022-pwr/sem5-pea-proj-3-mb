package solver

import "projekt2/graph"

type ATSPSolver interface {
	SetGraph(graph graph.Graph)
	GetGraph() graph.Graph
	SetStartVertex(startVertex int)
	Solve() ([]int, int)
}

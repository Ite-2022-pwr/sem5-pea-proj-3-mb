package gr

import (
	"projekt2/graph"
)

type GRATSPSolver struct {
	graph       graph.Graph
	startVertex int
}

func NewGreedyATSPSolver(sv int) GRATSPSolver {
	return GRATSPSolver{
		startVertex: sv,
	}
}

func (g *GRATSPSolver) SetGraph(graph graph.Graph) {
	g.graph = graph
}

func (g *GRATSPSolver) GetGraph() graph.Graph {
	return g.graph
}

func (g *GRATSPSolver) SetStartVertex(startVertex int) {
	g.startVertex = startVertex
}

func (g *GRATSPSolver) Solve() ([]int, int) {
	bestPath := g.graph.GetHamiltonianPathGreedy(g.startVertex)
	bestPathWeight := g.graph.CalculatePathWeight(bestPath)
	for i := 1; i < g.graph.GetVertexCount(); i++ {
		path := g.graph.GetHamiltonianPathGreedy(i)
		pathWeight := g.graph.CalculatePathWeight(path)
		if pathWeight < bestPathWeight {
			bestPath = path
			bestPathWeight = pathWeight
		}
	}
	return bestPath, bestPathWeight
}

package graph

type Graph interface {
	GetNoEdgeValue() int
	SetNoEdgeValue(int)
	GetVertexCount() int
	GetEdgeCount() int
	GetAllEdges() []Edge
	GetEdgesFromVertex(startVertex int) []Edge
	GetEdgesToVertex(endVertex int) []Edge
	GetEdge(startVertex, endVertex int) Edge
	GetMinEdgeFromWeight(vertex int) int
	AddEdge(startVertex, endVertex, weight int)
	RemoveEdge(startVertex, endVertex int)
	IsAdjacent(startVertex, endVertex int) bool
	CalculatePathWeight(path []int) int
	PathWithWeightsToString(path []int) string
	GetHamiltonianPathGreedy(startVertex int) []int
	GetHamiltonianPathRandom(startVertex int) []int
	ToString() string
}

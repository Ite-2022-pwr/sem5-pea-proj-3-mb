package graph

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
)

type AdjMatrixGraph struct {
	adjMatrix   [][]int
	vertexCount int
	edgeCount   int
	noEdgeValue int
}

func NewAdjMatrixGraph(vertexCount, noEdgeValue int) *AdjMatrixGraph {
	newGraph := new(AdjMatrixGraph)
	newGraph.vertexCount = vertexCount
	newGraph.noEdgeValue = noEdgeValue
	newGraph.edgeCount = -1
	newGraph.adjMatrix = make([][]int, vertexCount)
	for i := 0; i < vertexCount; i++ {
		newGraph.adjMatrix[i] = make([]int, vertexCount)
		for j := 0; j < vertexCount; j++ {
			newGraph.adjMatrix[i][j] = noEdgeValue
		}
	}
	return newGraph
}

func (a *AdjMatrixGraph) GetNoEdgeValue() int {
	return a.noEdgeValue
}

func (a *AdjMatrixGraph) SetNoEdgeValue(noEdgeValue int) {
	a.noEdgeValue = noEdgeValue
}

func (a *AdjMatrixGraph) GetVertexCount() int {
	return len(a.adjMatrix)
}

func (a *AdjMatrixGraph) GetEdgeCount() int {
	if a.edgeCount == -1 {
		count := 0
		for i := 0; i < a.GetVertexCount(); i++ {
			for j := 0; j < a.GetVertexCount(); j++ {
				if a.adjMatrix[i][j] != a.noEdgeValue {
					count++
				}
			}
		}
		a.edgeCount = count
	}
	return a.edgeCount
}

func (a *AdjMatrixGraph) GetAllEdges() []Edge {
	edges := make([]Edge, 0)
	for i := 0; i < a.GetVertexCount(); i++ {
		for j := 0; j < a.GetVertexCount(); j++ {
			if a.adjMatrix[i][j] != a.noEdgeValue {
				edges = append(edges, Edge{StartVertex: i, EndVertex: j, Weight: a.adjMatrix[i][j]})
			}
		}
	}
	return edges
}

func (a *AdjMatrixGraph) GetEdgesFromVertex(startVertex int) []Edge {
	edges := make([]Edge, 0)
	for i := 0; i < a.GetVertexCount(); i++ {
		if a.adjMatrix[startVertex][i] != a.noEdgeValue {
			edges = append(edges, Edge{StartVertex: startVertex, EndVertex: i, Weight: a.adjMatrix[startVertex][i]})
		}
	}
	return edges
}

func (a *AdjMatrixGraph) GetEdgesToVertex(endVertex int) []Edge {
	edges := make([]Edge, 0)
	for i := 0; i < a.GetVertexCount(); i++ {
		if a.adjMatrix[i][endVertex] != a.noEdgeValue {
			edges = append(edges, Edge{StartVertex: i, EndVertex: endVertex, Weight: a.adjMatrix[i][endVertex]})
		}
	}
	return edges
}

func (a *AdjMatrixGraph) GetEdge(startVertex, endVertex int) Edge {
	return Edge{StartVertex: startVertex, EndVertex: endVertex, Weight: a.adjMatrix[startVertex][endVertex]}
}

func (a *AdjMatrixGraph) GetMinEdgeFromWeight(vertex int) int {
	minEdge := math.MaxInt
	for i := 0; i < a.GetVertexCount(); i++ {
		if a.adjMatrix[vertex][i] < minEdge && a.adjMatrix[vertex][i] != a.noEdgeValue {
			minEdge = a.adjMatrix[vertex][i]
		}
	}
	return minEdge
}

func (a *AdjMatrixGraph) AddEdge(startVertex, endVertex, weight int) {
	a.adjMatrix[startVertex][endVertex] = weight
	a.edgeCount += 1
}

func (a *AdjMatrixGraph) RemoveEdge(startVertex, endVertex int) {
	a.adjMatrix[startVertex][endVertex] = a.noEdgeValue
	a.edgeCount -= 1
}

func (a *AdjMatrixGraph) IsAdjacent(startVertex, endVertex int) bool {
	return a.adjMatrix[startVertex][endVertex] != a.noEdgeValue
}

func (a *AdjMatrixGraph) CalculatePathWeight(path []int) int {
	weight := 0
	for i := 0; i < len(path)-1; i++ {
		weight += a.adjMatrix[path[i]][path[i+1]]
	}
	return weight
}

func (a *AdjMatrixGraph) PathWithWeightsToString(path []int) string {
	var out strings.Builder
	for i := 0; i < len(path)-1; i++ {
		out.WriteString("v" + strconv.Itoa(path[i]) + "--(" + strconv.Itoa(a.adjMatrix[path[i]][path[i+1]]) + ")-->")
	}
	out.WriteString("v" + strconv.Itoa(path[len(path)-1]))
	return out.String()
}

func (a *AdjMatrixGraph) GetHamiltonianPathGreedy(startVertex int) []int {
	visited := make([]bool, a.GetVertexCount())
	path := make([]int, 0)
	path = append(path, startVertex)
	visited[startVertex] = true
	currentVertex := startVertex
	for len(path) < a.GetVertexCount() {
		minEdgeWeight := math.MaxInt
		edgesFromCurrentVertex := a.GetEdgesFromVertex(currentVertex)
		nextVertex := -1
		for _, edge := range edgesFromCurrentVertex {
			if !visited[edge.EndVertex] && edge.Weight < minEdgeWeight {
				minEdgeWeight = edge.Weight
				nextVertex = edge.EndVertex
			}
		}
		path = append(path, nextVertex)
		currentVertex = nextVertex
		visited[currentVertex] = true
	}
	path = append(path, startVertex)
	return path
}

func (a *AdjMatrixGraph) GetHamiltonianPathRandom(startVertex int) []int {
	visited := make([]bool, a.GetVertexCount())
	path := make([]int, 0)
	path = append(path, startVertex)
	visited[startVertex] = true
	currentVertex := startVertex
	for len(path) < a.GetVertexCount() {
		nextVertex := -1
		for {
			nextVertex = rand.Int() % a.GetVertexCount()
			if !visited[nextVertex] {
				break
			}
		}
		path = append(path, nextVertex)
		currentVertex = nextVertex
		visited[currentVertex] = true
	}
	path = append(path, startVertex)
	return path
}

func (a *AdjMatrixGraph) ToString() string {
	var out strings.Builder

	// Nagłówki kolumn z przesunięciem dla wiersza nagłówka
	out.WriteString("\t|") // Pusty tabulator na początku dla wyrównania z wierszami
	for i := 0; i < a.GetVertexCount(); i++ {
		out.WriteString("v" + strconv.Itoa(i) + "\t|")
	}
	out.WriteString("\n")

	// Linia przerywana
	out.WriteString(strings.Repeat("-", (a.GetVertexCount()+1)*8) + "\n")

	// Wiersze z nagłówkami i wartościami macierzy
	for i := 0; i < a.GetVertexCount(); i++ {
		out.WriteString("v" + strconv.Itoa(i) + "\t|") // Nagłówek wiersza

		for j := 0; j < a.GetVertexCount(); j++ {
			out.WriteString(strconv.Itoa(a.adjMatrix[i][j]) + "\t|") // Wartość z macierzy z tabulatorem i kreską
		}
		out.WriteString("\n")

		// Linia przerywana, oprócz ostatniego wiersza
		if i < a.GetVertexCount()-1 {
			out.WriteString(strings.Repeat("-", (a.GetVertexCount()+1)*8) + "\n")
		}
	}

	return out.String()
}

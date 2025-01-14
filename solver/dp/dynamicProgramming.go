package dp

import (
	"log"
	"math"
	"projekt2/graph"
)

type DPATSPSolver struct {
	graph       graph.Graph
	startVertex int
}

func NewDynamicProgrammingATSPSolver(sv int) DPATSPSolver {
	return DPATSPSolver{
		startVertex: sv,
	}
}

func (d *DPATSPSolver) SetGraph(graph graph.Graph) {
	d.graph = graph
}

func (d *DPATSPSolver) GetGraph() graph.Graph {
	return d.graph
}

func (d *DPATSPSolver) SetStartVertex(startVertex int) {
	d.startVertex = startVertex
}

func (d *DPATSPSolver) Solve() ([]int, int) {
	log.Println("Rozpoczęcie programowania dynamicznego dla wierzchołka początkowego:", d.startVertex, "z liczbą wierzchołków:", d.graph.GetVertexCount())

	vertexCount := d.graph.GetVertexCount()
	allVisited := (1 << vertexCount) - 1 // Maskowanie dla wszystkich wierzchołków odwiedzonych
	// np. dla 4 wierzchołków: 1111 (czyli 15)
	// bo 1 << 4 = 10000, a 10000 - 1 = 1111
	// gdzie każdy bit odpowiada jednemu wierzchołkowi

	// Tworzenie mapy przechowującej koszty częściowych rozwiązań
	memo := make([][]int, vertexCount)
	parent := make([][]int, vertexCount) // Dodatkowa tablica, aby zapamiętać, skąd przychodzimy
	for i := range memo {
		memo[i] = make([]int, 1<<vertexCount)   // np dla 4 wierzchołków 16 (2^4) bo od 0000 do 1111
		parent[i] = make([]int, 1<<vertexCount) // Inicjalizacja ścieżki (skąd przychodzimy)
		for j := range memo[i] {
			memo[i][j] = math.MaxInt // Inicjalizacja maksymalnym kosztem
			parent[i][j] = -1        // Brak poprzedniego wierzchołka (na początku)
		}
	}

	// Inicjalizacja wierzchołka początkowego
	memo[d.startVertex][1<<d.startVertex] = 0

	// Dynamiczne przeliczanie wartości dla podproblemów
	for subset := 1; subset <= allVisited; subset++ {
		if (subset & (1 << d.startVertex)) == 0 {
			continue // Pomijamy zbiory, które nie zawierają startVertex
		}

		for currentVertex := 0; currentVertex < vertexCount; currentVertex++ {
			if (subset&(1<<currentVertex)) == 0 || currentVertex == d.startVertex {
				continue // Pomijamy wierzchołki, które nie są w bieżącym podzbiorze, lub są wierzchołkiem startowym
			}

			previousSubset := subset ^ (1 << currentVertex) // Podzbiór bez bieżącego wierzchołka

			// Szukamy minimalnego kosztu przejścia do bieżącego wierzchołka
			for prevVertex := 0; prevVertex < vertexCount; prevVertex++ {
				if (previousSubset & (1 << prevVertex)) == 0 {
					continue // Pomijamy, jeśli prevVertex nie jest w zbiorze
				}

				edge := d.graph.GetEdge(prevVertex, currentVertex)
				if edge.Weight == d.graph.GetNoEdgeValue() {
					continue // Pomijamy, jeśli nie ma krawędzi
				}

				if memo[prevVertex][previousSubset] == math.MaxInt {
					continue // Pomijamy, jeśli nie ma wartości dla tego podzbioru
				}

				newCost := memo[prevVertex][previousSubset] + edge.Weight
				if newCost < memo[currentVertex][subset] {
					memo[currentVertex][subset] = newCost
					parent[currentVertex][subset] = prevVertex // Zapamiętujemy poprzednika
				}
			}
		}
	}

	// Znalezienie minimalnej ścieżki powrotnej do wierzchołka startowego
	minCost := math.MaxInt
	lastVertex := -1
	for vertex := 0; vertex < vertexCount; vertex++ {
		if vertex == d.startVertex {
			continue
		}
		edge := d.graph.GetEdge(vertex, d.startVertex)
		if edge.Weight == d.graph.GetNoEdgeValue() {
			continue
		}

		if memo[vertex][allVisited] == math.MaxInt {
			continue // Pomijamy, jeśli nie ma rozwiązania dla tego podzbioru
		}

		totalCost := memo[vertex][allVisited] + edge.Weight
		if totalCost < minCost {
			minCost = totalCost
			lastVertex = vertex
		}
	}

	// Odtwarzanie najlepszej ścieżki
	if lastVertex == -1 {
		return nil, -1 // Jeśli nie znaleziono żadnej ścieżki
	}

	bestPath := []int{}
	currentVertex := lastVertex
	currentSubset := allVisited

	// Odtwarzanie trasy na podstawie zapisanych poprzedników
	for currentVertex != -1 {
		bestPath = append([]int{currentVertex}, bestPath...)
		nextVertex := parent[currentVertex][currentSubset]
		currentSubset ^= (1 << currentVertex)
		currentVertex = nextVertex
	}

	// Dodajemy wierzchołek startowy na końcu trasy, aby utworzyć cykl
	bestPath = append(bestPath, d.startVertex)

	log.Println("Zakończenie programowania dynamicznego dla wierzchołka początkowego:", d.startVertex, "z liczbą wierzchołków:", vertexCount)

	return bestPath, minCost
}

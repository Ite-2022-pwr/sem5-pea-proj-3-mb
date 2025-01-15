package ts

import (
	"fmt"
	"log"
	"math"
	"projekt3/graph"
	"time"
)

const (
	NeighborhoodSwap   = "swap"
	NeighborhoodInsert = "insert"
)

type TsATSPSolver struct {
	graph              graph.Graph
	startVertex        int
	iterations         int   // Maksymalna liczba iteracji
	timeout            int64 // Limit czasu w nanosekundach (-1 oznacza brak limitu)
	startTime          time.Time
	tabuTenure         int    // Ile iteracji ruch pozostaje tabu
	neighborhoodMethod string // Metoda sąsiedztwa: "swap" lub "insert"
}

func (t *TsATSPSolver) SetGraph(g graph.Graph) {
	t.graph = g
}

func (t *TsATSPSolver) GetGraph() graph.Graph {
	return t.graph
}

func (t *TsATSPSolver) SetStartVertex(startVertex int) {
	t.startVertex = startVertex
}

func (t *TsATSPSolver) SetTimeout(timeout int64) {
	t.timeout = timeout
}

func (t *TsATSPSolver) GetTimeout() int64 {
	return t.timeout
}

func (t *TsATSPSolver) SetNeighborhoodMethod(method string) error {
	if method != NeighborhoodSwap && method != NeighborhoodInsert {
		return fmt.Errorf("nieprawidłowa metoda sąsiedztwa: %s", method)
	}
	t.neighborhoodMethod = method
	return nil
}

func (t *TsATSPSolver) GetNeighborhoodMethod() string {
	return t.neighborhoodMethod
}

func NewTabuSearchATSPSolver(iterations int, timeout int64, tabuTenure int, neighborhoodMethod string) TsATSPSolver {
	solver := TsATSPSolver{
		iterations: iterations,
		timeout:    timeout,
		tabuTenure: tabuTenure,
	}
	if err := solver.SetNeighborhoodMethod(neighborhoodMethod); err != nil {
		log.Printf("Nieprawidłowa metoda sąsiedztwa: %s. Użyto domyślnej metody '%s'.\n", neighborhoodMethod, NeighborhoodSwap)
		solver.neighborhoodMethod = NeighborhoodSwap
	}
	return solver
}

// calculateCost oblicza koszt danej ścieżki
func (t *TsATSPSolver) calculateCost(path []int) int {
	return t.graph.CalculatePathWeight(path)
}

// findBestNeighbor znajduje najlepszego sąsiada (zgodnie z metodą sąsiedztwa)
func (t *TsATSPSolver) findBestNeighbor(currentSolution []int, tabuList [][]int, bestCost int) (bestPath []int, bestI, bestJ, bestNeighborCost int) {
	vertexCount := len(currentSolution)
	bestNeighborCost = math.MaxInt32

	switch t.neighborhoodMethod {
	case NeighborhoodSwap:
		// Swap: zamiana par wierzchołków
		for i := 1; i < vertexCount-1; i++ {
			for j := i + 1; j < vertexCount-1; j++ {
				// Zamiana
				currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]

				cost := t.calculateCost(currentSolution)
				isTabu := (tabuList[i][j] > 0 || tabuList[j][i] > 0)
				// Warunek aspiracji lub nie-tabu
				if cost < bestNeighborCost && (!isTabu || cost < bestCost) {
					bestNeighborCost = cost
					// Kopiujemy ścieżkę tylko kiedy jest to potrzebne (najlepszy dotąd)
					if bestPath == nil {
						bestPath = make([]int, vertexCount)
					}
					copy(bestPath, currentSolution)
					bestI, bestJ = i, j
				}

				// Cofnięcie zamiany
				currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]
			}
		}

	case NeighborhoodInsert:
		// Insert: przeniesienie wierzchołka z pozycji i na pozycję j
		for i := 1; i < vertexCount-1; i++ {
			for j := 1; j < vertexCount-1; j++ {
				if i == j {
					continue
				}
				// Tymczasowo zmodyfikujemy currentSolution, aby uzyskać sąsiada
				// Zapisz wartość elementu do przeniesienia
				element := currentSolution[i]

				if i < j {
					// Przesuwamy wierzchołki w lewo
					copy(currentSolution[i:], currentSolution[i+1:j+1])
					currentSolution[j] = element
				} else {
					// i > j
					// Przesuwamy wierzchołki w prawo
					copy(currentSolution[j+1:i+1], currentSolution[j:i])
					currentSolution[j] = element
				}

				cost := t.calculateCost(currentSolution)
				isTabu := (tabuList[i][j] > 0 || tabuList[j][i] > 0)
				if cost < bestNeighborCost && (!isTabu || cost < bestCost) {
					bestNeighborCost = cost
					if bestPath == nil {
						bestPath = make([]int, vertexCount)
					}
					copy(bestPath, currentSolution)
					bestI, bestJ = i, j
				}

				// Przywracamy oryginalną kolejność
				if i < j {
					// element był w i, przeniesiony do j
					// cofamy zmianę
					copy(currentSolution[i+1:j+1], currentSolution[i:j])
					currentSolution[i] = element
				} else {
					// i > j
					copy(currentSolution[j:i], currentSolution[j+1:i+1])
					currentSolution[i] = element
				}
			}
		}

	default:
		// Domyślne podejście swap, jeśli coś jest nie tak
		for i := 1; i < vertexCount-1; i++ {
			for j := i + 1; j < vertexCount-1; j++ {
				currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]

				cost := t.calculateCost(currentSolution)
				isTabu := (tabuList[i][j] > 0 || tabuList[j][i] > 0)
				if cost < bestNeighborCost && (!isTabu || cost < bestCost) {
					bestNeighborCost = cost
					if bestPath == nil {
						bestPath = make([]int, vertexCount)
					}
					copy(bestPath, currentSolution)
					bestI, bestJ = i, j
				}

				currentSolution[i], currentSolution[j] = currentSolution[j], currentSolution[i]
			}
		}
	}

	return
}

func (t *TsATSPSolver) Solve() ([]int, int) {
	vertexCount := t.graph.GetVertexCount()
	if vertexCount == 0 {
		return nil, -1
	}

	t.startTime = time.Now()

	// Generujemy początkowe losowe rozwiązanie
	currentSolution := t.graph.GetHamiltonianPathRandom(t.startVertex)
	currentCost := t.calculateCost(currentSolution)

	bestSolution := make([]int, len(currentSolution))
	copy(bestSolution, currentSolution)
	bestCost := currentCost

	tabuList := make([][]int, vertexCount)
	for i := 0; i < vertexCount; i++ {
		tabuList[i] = make([]int, vertexCount)
	}

	for iteration := 0; iteration < t.iterations; iteration++ {
		// Sprawdzenie limitu czasu na początku każdej iteracji
		if t.timeout != -1 {
			elapsed := time.Since(t.startTime).Nanoseconds()
			if elapsed >= t.timeout {
				log.Println("Zatrzymano przy iteracji:", iteration, "z powodu przekroczenia limitu czasu.")
				break
			}
		}

		// Znajdujemy najlepszego sąsiada
		newSolution, bestI, bestJ, neighborCost := t.findBestNeighbor(currentSolution, tabuList, bestCost)

		if newSolution == nil {
			// Brak sąsiadów lub nie udało się poprawić, kończymy
			break
		}

		// Przejście do wybranego sąsiada
		copy(currentSolution, newSolution)
		currentCost = neighborCost

		// Aktualizacja najlepszego rozwiązania
		if currentCost < bestCost {
			bestCost = currentCost
			copy(bestSolution, currentSolution)
		}

		// Aktualizacja listy tabu:
		tabuList[bestI][bestJ] = t.tabuTenure
		tabuList[bestJ][bestI] = t.tabuTenure

		// Zmniejszanie karencji tabu
		for i := 0; i < vertexCount; i++ {
			for j := 0; j < vertexCount; j++ {
				if tabuList[i][j] > 0 {
					tabuList[i][j]--
				}
			}
		}
	}

	log.Println("Zakończono Tabu Search. Najlepszy znaleziony koszt:", bestCost)
	return bestSolution, bestCost
}

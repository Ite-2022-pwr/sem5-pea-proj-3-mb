package aco

import (
	"fmt"
	"log"
	"math"
	"math/rand/v2"
	"projekt2/graph"
	"time"
)

// Małe epsilon, aby uniknąć dzielenia przez 0, jeśli krawędź ma wagę 0.
const epsilon = 1e-9

// ACOASSolver - przykład algorytmu Mrówkowego spełniającego podane warunki.
type ACOASSolver struct {
	graph                           graph.Graph
	antsCount                       int         // Liczba mrówek
	pheromonesPerAnt                float64     // Ilość feromonów na jedną mrówkę
	iterations                      int         // Liczba iteracji
	alpha                           float64     // Wpływ (waga) feromonów
	beta                            float64     // Wpływ (waga) heurystyki (1 / (waga + epsilon))
	evaporationRate                 float64     // Współczynnik parowania feromonów
	startPheromones                 float64     // Początkowa ilość feromonów
	pheromones                      [][]float64 // Macierz feromonów
	decisionMatrix                  [][]float64 // Macierz decyzji
	invDistancesToBetaPow           [][]float64 // Macierz odwrotności wag krawędzi
	bestSolution                    []int       // Najlepsze dotąd znalezione rozwiązanie
	bestCost                        int         // Koszt najlepszej trasy
	timeout                         int64       // Czas wykonania w nanosekundach
	maxIterationsWithoutImprovement int         // Liczba iteracji bez poprawy
	startTime                       time.Time   // Czas rozpoczęcia
}

// NewACOZeroEdgeSolver - konstruktor algorytmu
func NewACOZeroEdgeSolver(antsCount, iterations, maxIterationsWithoutImprovement int, alpha, beta, evap, pherPA, startPher float64, timeout int64) *ACOASSolver {
	log.Println(timeout)
	return &ACOASSolver{
		antsCount:                       antsCount, // recommended: graph.GetVertexCount()
		pheromonesPerAnt:                pherPA,    // default: 5.0 recommended(?): graph.CalculatePathWeight(graph.GetHamiltonianPathGreedy(0))
		iterations:                      iterations,
		alpha:                           alpha,       // default: 1.0
		beta:                            beta,        // default: 2.0-5.0
		evaporationRate:                 evap,        // default: 0.1-0.5
		startPheromones:                 startPher,   // default: 1.0 recommended: antsCount / graph.CalculatePathWeight(graph.GetHamiltonianPathGreedy(0))
		bestCost:                        math.MaxInt, // Na start przyjmujemy bardzo dużą wartość
		maxIterationsWithoutImprovement: maxIterationsWithoutImprovement,
		timeout:                         timeout,
	}
}

// SetGraph przypina solverowi dany graf.
func (s *ACOASSolver) SetGraph(g graph.Graph) {
	s.graph = g
	if s.startPheromones == 0 {
		s.startPheromones = 1.0
	}
	s.initializePheromones()
	s.initializeInvDistancesToBetaPow()
	s.initializeDecisionMatrix()
}

// GetGraph zwraca aktualnie ustawiony graf.
func (s *ACOASSolver) GetGraph() graph.Graph {
	return s.graph
}

// GetPheromones zwraca macierz feromonów.
func (s *ACOASSolver) GetPheromones() [][]float64 {
	return s.pheromones
}

// GetDecisionMatrix zwraca macierz decyzji.
func (s *ACOASSolver) GetDecisionMatrix() [][]float64 {
	return s.decisionMatrix
}

// PheromonesToString zwraca macierz feromonów w postaci stringa.
func (s *ACOASSolver) PheromonesToString() string {
	result := ""
	for i := 0; i < len(s.pheromones); i++ {
		for j := 0; j < len(s.pheromones[i]); j++ {
			result += fmt.Sprintf("|%8.2f ", s.pheromones[i][j])
		}
		result += "|\n"
	}
	return result
}

// DecisionMatrixToString zwraca macierz decyzji w postaci stringa.
func (s *ACOASSolver) DecisionMatrixToString() string {
	result := ""
	for i := 0; i < len(s.decisionMatrix); i++ {
		for j := 0; j < len(s.decisionMatrix[i]); j++ {
			result += fmt.Sprintf("|%8.2f ", s.decisionMatrix[i][j])
		}
		result += "|\n"
	}
	return result
}

// Solve uruchamia algorytm i zwraca najlepszą znalezioną ścieżkę Hamiltona wraz z kosztem.
// Jeżeli żadna mrówka nie zbuduje pełnej ścieżki, solver zwróci bestSolution == nil i bestCost == math.MaxInt.
func (s *ACOASSolver) Solve() ([]int, int) {
	vertexCount := s.graph.GetVertexCount()
	if vertexCount == 0 {
		return nil, -1
	}

	// Rejestracja czasu rozpoczęcia
	s.startTime = time.Now()

	iterationsWithoutImprovement := 0

	for i := 0; i < s.iterations; i++ {
		// Sprawdzenie limitu czasu
		if s.timeout != -1 {
			elapsed := time.Since(s.startTime).Nanoseconds()
			if elapsed >= s.timeout {
				log.Println("Przekroczono limit czasu. Kończenie algorytmu.")
				break
			}
		}
		if iterationsWithoutImprovement >= s.maxIterationsWithoutImprovement {
			log.Println("Przekroczono limit iteracji bez poprawy. Kończenie algorytmu.")
			break
		}

		iterationsWithoutImprovement++

		antPaths := make([][]int, s.antsCount)
		antCosts := make([]int, s.antsCount)
		for j := 0; j < s.antsCount; j++ {
			antPath := make([]int, 0)
			antCost := 0
			antVisited := make([]bool, vertexCount)
			antStartVertex := rand.Int() % vertexCount
			antPath = append(antPath, antStartVertex)
			antVisited[antStartVertex] = true
			for k := 0; k < vertexCount-1; k++ {
				currentVertex := antPath[len(antPath)-1]
				edgeDecisionValues := s.decisionMatrix[currentVertex]
				nextVertex := -1
				availableProbabilitiesSum := s.calculateAvailableProbabilitiesSum(antVisited, currentVertex)
				sumBefore := 0.0
				currSum := 0.0
				randomValue := rand.Float64()
				for l := 0; l < vertexCount; l++ {
					if !antVisited[l] && s.graph.GetEdge(currentVertex, l).Weight != s.graph.GetNoEdgeValue() {
						currSum += edgeDecisionValues[l] / availableProbabilitiesSum
						if randomValue >= sumBefore && randomValue < currSum {
							nextVertex = l
							break
						}
						sumBefore = currSum
					}
				}
				if nextVertex == -1 {
					for l := 0; l < vertexCount; l++ {
						if !antVisited[l] {
							nextVertex = l
							break
						}
					}
					if nextVertex == -1 {
						panic("Nie znaleziono kolejnego wierzchołka")
					}
				}
				antVisited[nextVertex] = true
				currentVertex = nextVertex
				antPath = append(antPath, nextVertex)
			}
			antPath = append(antPath, antStartVertex)
			antCost = s.graph.CalculatePathWeight(antPath)
			antPaths[j] = make([]int, len(antPath))
			copy(antPaths[j], antPath)
			antCosts[j] = antCost
			if antCost < s.bestCost {
				log.Println("Znaleziono lepsze rozwiązanie!")
				log.Println("Mrówka", j, "iteracja", i, "koszt:", antCost)
				log.Println("Mrówka", j, "iteracja", i, "ścieżka:", antPath)
				s.bestCost = antCost
				s.bestSolution = make([]int, len(antPath))
				copy(s.bestSolution, antPath)
				iterationsWithoutImprovement = 0
			}
		}
		//fmt.Println(s.PheromonesToString())
		s.updatePheromones(antPaths, antCosts)
		s.updateDecisionMatrix()
	}
	return s.bestSolution, s.bestCost
}

// initializePheromones inicjalizuje macierz feromonów.
func (s *ACOASSolver) initializePheromones() {
	vertexCount := s.graph.GetVertexCount()
	s.pheromones = make([][]float64, vertexCount)
	for i := 0; i < vertexCount; i++ {
		s.pheromones[i] = make([]float64, vertexCount)
		for j := 0; j < vertexCount; j++ {
			s.pheromones[i][j] = s.startPheromones
		}
	}
}

// evaporationPheromones paruje feromony.
func (s *ACOASSolver) evaporationPheromones() {
	vertexCount := s.graph.GetVertexCount()
	for i := 0; i < vertexCount; i++ {
		for j := 0; j < vertexCount; j++ {
			s.pheromones[i][j] *= 1 - s.evaporationRate
		}
	}
}

// updatePheromones aktualizuje macierz feromonów po zakończeniu iteracji.
func (s *ACOASSolver) updatePheromones(antPaths [][]int, antCosts []int) {
	s.evaporationPheromones()
	for i := 0; i < len(antPaths); i++ {
		//pheromonePerEdge := s.pheromonesPerAnt / float64(antCosts[i])
		for j := 0; j < len(antPaths[i])-1; j++ {
			s.pheromones[antPaths[i][j]][antPaths[i][j+1]] += s.pheromonesPerAnt / float64(antCosts[i])
		}
	}
}

// initializeInvDistances inicjalizuje macierz odwrotności wag krawędzi.
func (s *ACOASSolver) initializeInvDistancesToBetaPow() {
	vertexCount := s.graph.GetVertexCount()
	s.invDistancesToBetaPow = make([][]float64, vertexCount)
	for i := 0; i < vertexCount; i++ {
		s.invDistancesToBetaPow[i] = make([]float64, vertexCount)
	}
	for i := 0; i < vertexCount; i++ {
		for j := 0; j < vertexCount; j++ {
			if i == j || s.graph.GetEdge(i, j).Weight == s.graph.GetNoEdgeValue() {
				s.invDistancesToBetaPow[i][j] = 0
			} else {
				if s.graph.GetEdge(i, j).Weight == 0 {
					s.invDistancesToBetaPow[i][j] = 1.0 / epsilon
				} else {
					s.invDistancesToBetaPow[i][j] = 1.0 / float64(s.graph.GetEdge(i, j).Weight)
				}
				s.invDistancesToBetaPow[i][j] = math.Pow(s.invDistancesToBetaPow[i][j], s.beta)
			}
		}
	}
}

// initializeDecisionMatrix inicjalizuje macierz decyzji.
func (s *ACOASSolver) initializeDecisionMatrix() {
	vertexCount := s.graph.GetVertexCount()
	s.decisionMatrix = make([][]float64, vertexCount)
	for i := 0; i < vertexCount; i++ {
		s.decisionMatrix[i] = make([]float64, vertexCount)
	}
	s.updateDecisionMatrix()
}

// updateDecisionMatrix aktualizuje macierz decyzji.
func (s *ACOASSolver) updateDecisionMatrix() {
	vertexCount := s.graph.GetVertexCount()
	for i := 0; i < vertexCount; i++ {
		edgeAttractivenessList := make([]float64, vertexCount)
		edgeAttractivenessSum := 0.0
		for j := 0; j < vertexCount; j++ {
			if i == j || s.graph.GetEdge(i, j).Weight == s.graph.GetNoEdgeValue() {
				s.decisionMatrix[i][j] = -1.0
				continue
			}
			//edgeAttractivenessList[j] = math.Pow(s.pheromones[i][j], s.alpha) * s.invDistancesToBetaPow[i][j]
			edgeAttractivenessList[j] = math.Pow(s.pheromones[i][j], s.alpha) * s.invDistancesToBetaPow[i][j]
			edgeAttractivenessSum += edgeAttractivenessList[j]
		}
		for j := 0; j < vertexCount; j++ {
			if i != j && s.graph.GetEdge(i, j).Weight != s.graph.GetNoEdgeValue() {
				s.decisionMatrix[i][j] = edgeAttractivenessList[j] / edgeAttractivenessSum
			}
		}
	}
}

// calculateAvailableProbabilitiesSum oblicza sumę dostępnych prawdopodobieństw.
func (s *ACOASSolver) calculateAvailableProbabilitiesSum(antVisited []bool, currentVertex int) float64 {
	vertexCount := s.graph.GetVertexCount()
	availableProbabilitiesSum := 0.0
	for i := 0; i < vertexCount; i++ {
		if !antVisited[i] && s.graph.GetEdge(currentVertex, i).Weight != s.graph.GetNoEdgeValue() {
			availableProbabilitiesSum += s.decisionMatrix[currentVertex][i]
		}
	}
	return availableProbabilitiesSum
}

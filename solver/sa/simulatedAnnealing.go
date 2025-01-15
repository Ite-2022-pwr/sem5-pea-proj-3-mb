package sa

import (
	"log"
	"math"
	"math/rand"
	"projekt3/graph"
	"time"
)

// SaATSPSolver implementuje interfejs ATSPSolver
type SaATSPSolver struct {
	graph              graph.Graph
	startVertex        int
	initialTemperature float64   // Początkowa temperatura
	minimalTemperature float64   // Minimalna temperatura
	alpha              float64   // Współczynnik chłodzenia, np. 0.99
	iterations         int       // Liczba iteracji
	timeout            int64     // Czas wykonania w nanosekundach
	startTime          time.Time // Czas rozpoczęcia
}

// SetGraph ustawia graf dla solvera
func (s *SaATSPSolver) SetGraph(g graph.Graph) {
	s.graph = g
}

// GetGraph zwraca przypisany graf
func (s *SaATSPSolver) GetGraph() graph.Graph {
	return s.graph
}

// SetStartVertex ustawia wierzchołek startowy
func (s *SaATSPSolver) SetStartVertex(startVertex int) {
	s.startVertex = startVertex
}

// SetTimeout ustawia czas wykonania w sekundach
func (s *SaATSPSolver) SetTimeout(timeout int64) {
	s.timeout = timeout
}

// GetTimeout zwraca czas wykonania w sekundach
func (s *SaATSPSolver) GetTimeout() int64 {
	return s.timeout
}

// NewSimulatedAnnealingATSPSolver tworzy nowy solver
func NewSimulatedAnnealingATSPSolver(initialTemperature float64, minimalTemperature float64, alpha float64, iterations int, timeout int64) SaATSPSolver {
	return SaATSPSolver{
		initialTemperature: initialTemperature,
		minimalTemperature: minimalTemperature,
		alpha:              alpha,
		iterations:         iterations,
		timeout:            timeout,
	}
}

// calculateCost oblicza koszt danej ścieżki w grafie, wliczając powrót do startu
// Zakłada, że path już kończy się na startVertex, więc nie dodaje go ponownie.
func (s *SaATSPSolver) calculateCost(path []int) int {
	return s.graph.CalculatePathWeight(path)
}

// getNeighbor generuje sąsiednie rozwiązanie poprzez zamianę pozycji dwóch wierzchołków (oprócz startVertex na początku i końca)
func (s *SaATSPSolver) getNeighbor(currentPath []int) []int {
	vertexCount := len(currentPath)
	if vertexCount <= 3 {
		// Jeżeli jest tylko start i jeden inny wierzchołek, lub tylko startVertex na początku i końcu
		return currentPath
	}

	newPath := make([]int, vertexCount)
	copy(newPath, currentPath)

	// Losowanie dwóch pozycji do zamiany, pomijamy indeks 0 (startVertex) i ostatni indeks (również startVertex)
	i := rand.Intn(vertexCount-2) + 1 // [1, vertexCount-2]
	j := rand.Intn(vertexCount-2) + 1
	for j == i {
		j = rand.Intn(vertexCount-2) + 1
	}

	// Zamiana
	newPath[i], newPath[j] = newPath[j], newPath[i]

	return newPath
}

// acceptanceProbability oblicza prawdopodobieństwo przyjęcia gorszego rozwiązania
func (s *SaATSPSolver) acceptanceProbability(delta int, temperature float64) float64 {
	return math.Exp(-float64(delta) / temperature)
}

// Solve rozwiązuje ATSP metodą Symulowanego Wyżarzania
func (s *SaATSPSolver) Solve() ([]int, int) {
	vertexCount := s.graph.GetVertexCount()
	if vertexCount == 0 {
		return nil, -1
	}

	// Rejestracja czasu rozpoczęcia
	s.startTime = time.Now()

	// Wygeneruj początkowe rozwiązanie metodą zachłanną
	currentSolution := s.graph.GetHamiltonianPathRandom(s.startVertex)
	currentCost := s.calculateCost(currentSolution)

	// Ustawiamy najlepsze znane rozwiązanie
	bestSolution := make([]int, len(currentSolution))
	copy(bestSolution, currentSolution)
	bestCost := currentCost

	// Inicjalizacja parametrów SA
	T := s.initialTemperature

	// Główna pętla symulowanego wyżarzania
	for T > s.minimalTemperature {
		// Sprawdzenie limitu czasu
		if s.timeout != -1 {
			elapsed := time.Since(s.startTime).Nanoseconds()
			if elapsed >= s.timeout {
				log.Println("Przekroczono limit czasu. Kończenie algorytmu.")
				break
			}
		}

		for iteration := 0; iteration < s.iterations; iteration++ {
			// Generujemy sąsiada
			newSolution := s.getNeighbor(currentSolution)
			newCost := s.calculateCost(newSolution)
			delta := newCost - currentCost

			if delta < 0 {
				// Lepsze rozwiązanie
				currentSolution = newSolution
				currentCost = newCost
			} else {
				// Gorsze rozwiązanie - sprawdzamy prawdopodobieństwo przyjęcia
				ap := s.acceptanceProbability(delta, T)
				chance := rand.Float64()
				if chance < ap {
					currentSolution = newSolution
					currentCost = newCost
				}
			}

			// Aktualizacja najlepszego znalezionego rozwiązania
			if currentCost < bestCost {
				bestCost = currentCost
				copy(bestSolution, currentSolution)
			}
		}

		// Schładzanie temperatury
		T *= s.alpha
	}

	log.Println("Zakończono Symulowane Wyżarzanie. Najlepszy znaleziony koszt:", bestCost)
	log.Println("Temperatura końcowa:", T)
	log.Println("wartoś exp(-1/Tk) =", s.acceptanceProbability(1, T))

	// bestSolution już kończy się na startVertex, więc nie musimy go doklejać
	return bestSolution, bestCost
}

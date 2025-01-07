package menu

import (
	"bufio"
	"fmt"
	"os"
	"projekt2/utils"
	"strconv"
	"strings"
	"time"

	"projekt2/graph"
	"projekt2/solver/bf"
	"projekt2/solver/bnb"
	"projekt2/solver/dp"
	"projekt2/solver/gr"
	"projekt2/solver/sa"
	"projekt2/solver/ts"
)

// Menu struktura obsługująca dostępne funkcjonalności
type Menu struct {
	bfATSPSolver  bf.BFATSPSolver
	bnbATSPSolver bnb.BNBATSPSolver
	dpATSPSolver  dp.DPATSPSolver
	grATSPSolver  gr.GRATSPSolver
	saATSPSolver  sa.SaATSPSolver
	tsATSPSolver  ts.TsATSPSolver
	graph         graph.Graph
	startVertex   int
}

// NewMenu tworzy nową instancję menu bez grafu
func NewMenu() *Menu {
	return &Menu{
		startVertex:   0,
		bfATSPSolver:  bf.BFATSPSolver{},
		bnbATSPSolver: bnb.BNBATSPSolver{},
		dpATSPSolver:  dp.DPATSPSolver{},
		grATSPSolver:  gr.GRATSPSolver{},
		saATSPSolver:  sa.SaATSPSolver{},
		tsATSPSolver:  ts.TsATSPSolver{},
	}
}

// NewDefaultMenu tworzy nową instancję menu z podanym grafem
func NewDefaultMenu(g graph.Graph) *Menu {
	return &Menu{
		graph: g,
	}
}

// SetGraph ustawia graf
func (m *Menu) SetGraph(g graph.Graph) {
	m.graph = g
	m.bnbATSPSolver.SetGraph(g)
	m.bfATSPSolver.SetGraph(g)
	m.dpATSPSolver.SetGraph(g)
	m.grATSPSolver.SetGraph(g)
	m.saATSPSolver.SetGraph(g)
	m.tsATSPSolver.SetGraph(g)
}

// Graph zwraca aktualny graf
func (m *Menu) Graph() graph.Graph {
	return m.graph
}

// SetStartVertex ustawia wierzchołek startowy dla problemu
func (m *Menu) SetStartVertex(startVertex int) {
	m.startVertex = startVertex
	m.bfATSPSolver.SetStartVertex(startVertex)
	m.bnbATSPSolver.SetStartVertex(startVertex)
	m.dpATSPSolver.SetStartVertex(startVertex)
	m.grATSPSolver.SetStartVertex(startVertex)
	m.saATSPSolver.SetStartVertex(startVertex)
	m.tsATSPSolver.SetStartVertex(startVertex)
}

// LoadGraphFromFile wczytuje graf z pliku
func (m *Menu) LoadGraphFromFile(filePath string, isATSP bool) error {
	adjGraph := &graph.AdjMatrixGraph{}
	err := graph.LoadGraphFromFile(filePath, adjGraph, isATSP)
	if err != nil {
		return err
	}
	m.SetGraph(adjGraph)
	return nil
}

// GenerateRandomGraph generuje losowy graf
func (m *Menu) GenerateRandomGraph(vertexCount, noEdgeValue, maxWeight int) {
	adjGraph := &graph.AdjMatrixGraph{}
	graph.GenerateRandomGraph(adjGraph, vertexCount, noEdgeValue, maxWeight)
	m.SetGraph(adjGraph)
}

// SetNoEdgeValue ustawia wartość braku krawędzi
func (m *Menu) SetNoEdgeValue(value int) {
	if m.graph != nil {
		m.graph.SetNoEdgeValue(value)
	} else {
		fmt.Println("Brak grafu, nie można ustawić wartości noEdgeValue.")
	}
}

// DisplayGraph wyświetla aktualny graf
func (m *Menu) DisplayGraph() {
	if m.graph != nil {
		fmt.Println(m.graph.ToString())
	} else {
		fmt.Println("Graf nie został zainicjalizowany.")
	}
}

// Konfiguracja solverów bez dodatkowych parametrów - tylko start vertex
func (m *Menu) ConfigureBfSolver() {
	m.bfATSPSolver = bf.NewBruteForceATSPSolver(m.startVertex)
	m.bfATSPSolver.SetGraph(m.graph)
	fmt.Println("BF skonfigurowane.")
}

func (m *Menu) ConfigureBnbSolver() {
	m.bnbATSPSolver = bnb.NewBranchAndBoundATSPSolver(m.startVertex)
	m.bnbATSPSolver.SetGraph(m.graph)
	fmt.Println("BnB skonfigurowane.")
}

func (m *Menu) ConfigureDpSolver() {
	m.dpATSPSolver = dp.NewDynamicProgrammingATSPSolver(m.startVertex)
	m.dpATSPSolver.SetGraph(m.graph)
	fmt.Println("DP skonfigurowane.")
}

func (m *Menu) ConfigureGrSolver() {
	m.grATSPSolver = gr.NewGreedyATSPSolver(m.startVertex)
	m.grATSPSolver.SetGraph(m.graph)
	fmt.Println("GR skonfigurowane.")
}

// Konfiguracja solverów z parametrami
func (m *Menu) ConfigureSaSolver(initialTemperature float64, minimalTemperature float64, alpha float64, iterations int, timeout int64) {
	solver := sa.NewSimulatedAnnealingATSPSolver(initialTemperature, minimalTemperature, alpha, iterations, timeout)
	solver.SetGraph(m.graph)
	solver.SetStartVertex(m.startVertex)
	m.saATSPSolver = solver
	fmt.Println("SA skonfigurowane.")
}

func (m *Menu) ConfigureTsSolver(iterations int, timeout int64, tabuTenure int, neighborhoodMethod string) {
	// Inicjalizacja solvera Tabu Search
	solver := ts.NewTabuSearchATSPSolver(iterations, timeout, tabuTenure, neighborhoodMethod)
	solver.SetGraph(m.graph)
	solver.SetStartVertex(m.startVertex)

	// Ustawienie metody sąsiedztwa
	err := solver.SetNeighborhoodMethod(neighborhoodMethod)
	if err != nil {
		fmt.Printf("Błąd ustawiania metody sąsiedztwa: %v. Użyto domyślnej metody 'swap'.\n", err)
	}
	m.tsATSPSolver = solver
	fmt.Println("TS skonfigurowane.")
}

// printSolution wypisuje ścieżkę i koszt
func (m *Menu) printSolution(path []int, cost int) {
	if path == nil {
		fmt.Println("Nie znaleziono rozwiązania.")
		return
	}
	fmt.Println("Koszt:", cost)
	if m.graph != nil {
		fmt.Println("Ścieżka ze szczegółami wag:", m.graph.PathWithWeightsToString(path))
	}
}

// RunBf uruchamia brute force
func (m *Menu) RunBf() {
	if m.bfATSPSolver.GetGraph() == nil {
		fmt.Println("Solver BF nie ma przypisanego grafu.")
		return
	}
	start := time.Now()
	path, cost := m.bfATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania Brute Force:", elapsed)
}

// RunBnb uruchamia branch and bound
func (m *Menu) RunBnb() {
	if m.bnbATSPSolver.GetGraph() == nil {
		fmt.Println("Solver BnB nie ma przypisanego grafu.")
		return
	}
	start := time.Now()
	path, cost := m.bnbATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania BnB:", elapsed)
}

// RunDp uruchamia dynamic programming
func (m *Menu) RunDp() {
	if m.dpATSPSolver.GetGraph() == nil {
		fmt.Println("Solver DP nie ma przypisanego grafu.")
		return
	}
	start := time.Now()
	path, cost := m.dpATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania DP:", elapsed)
}

// RunGr uruchamia greedy solver
func (m *Menu) RunGr() {
	if m.grATSPSolver.GetGraph() == nil {
		fmt.Println("Solver GR nie ma przypisanego grafu.")
		return
	}
	start := time.Now()
	path, cost := m.grATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania Greedy:", elapsed)
}

// RunSa uruchamia simulated annealing
func (m *Menu) RunSa() {
	if m.saATSPSolver.GetGraph() == nil {
		fmt.Println("Solver SA nie ma przypisanego grafu.")
		return
	}
	start := time.Now()
	path, cost := m.saATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania SA:", elapsed)
}

// RunTs uruchamia tabu search
func (m *Menu) RunTs() {
	if m.tsATSPSolver.GetGraph() == nil {
		fmt.Println("Solver TS nie jest skonfigurowany.")
		return
	}
	start := time.Now()
	path, cost := m.tsATSPSolver.Solve()
	elapsed := time.Since(start)
	m.printSolution(path, cost)
	fmt.Println("Czas wykonania Tabu Search:", elapsed)
}

// SaveGraphToFile zapis grafu do pliku
func (m *Menu) SaveGraphToFile(filePath string, useTabs bool) error {
	if m.graph == nil {
		return fmt.Errorf("brak grafu do zapisania")
	}
	err := graph.SaveGraphToFile(m.graph, filePath, useTabs)
	if err != nil {
		return err
	}
	fmt.Println("Graf zapisany do pliku:", filePath)
	return nil
}

// Submenu konfiguracji solverów
func (m *Menu) solverConfigurationSubmenu() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\n=== Konfiguracja Solverów ===")
		fmt.Println("1. BF")
		fmt.Println("2. BnB")
		fmt.Println("3. DP")
		fmt.Println("4. GR")
		fmt.Println("5. SA")
		fmt.Println("6. TS")
		fmt.Println("b. Powrót")
		fmt.Print("Wybierz solver do konfiguracji: ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "1":
			// BF - tylko ustawienie start vertex
			fmt.Print("Podaj wierzchołek startowy (lub enter aby nie zmieniać): ")
			sv, _ := reader.ReadString('\n')
			sv = strings.TrimSpace(sv)
			if sv != "" {
				newSV, err := strconv.Atoi(sv)
				if err == nil {
					m.SetStartVertex(newSV)
				} else {
					fmt.Println("Nieprawidłowa wartość. Start vertex pozostaje bez zmian.")
				}
			}
			m.ConfigureBfSolver()
		case "2":
			// BnB - tylko ustawienie start vertex
			fmt.Print("Podaj wierzchołek startowy (lub enter aby nie zmieniać): ")
			sv, _ := reader.ReadString('\n')
			sv = strings.TrimSpace(sv)
			if sv != "" {
				newSV, err := strconv.Atoi(sv)
				if err == nil {
					m.SetStartVertex(newSV)
				} else {
					fmt.Println("Nieprawidłowa wartość. Start vertex pozostaje bez zmian.")
				}
			}
			m.ConfigureBnbSolver()
		case "3":
			// DP - tylko ustawienie start vertex
			fmt.Print("Podaj wierzchołek startowy (lub enter aby nie zmieniać): ")
			sv, _ := reader.ReadString('\n')
			sv = strings.TrimSpace(sv)
			if sv != "" {
				newSV, err := strconv.Atoi(sv)
				if err == nil {
					m.SetStartVertex(newSV)
				} else {
					fmt.Println("Nieprawidłowa wartość. Start vertex pozostaje bez zmian.")
				}
			}
			m.ConfigureDpSolver()
		case "4":
			// GR - tylko ustawienie start vertex
			fmt.Print("Podaj wierzchołek startowy (lub enter aby nie zmieniać): ")
			sv, _ := reader.ReadString('\n')
			sv = strings.TrimSpace(sv)
			if sv != "" {
				newSV, err := strconv.Atoi(sv)
				if err == nil {
					m.SetStartVertex(newSV)
				} else {
					fmt.Println("Nieprawidłowa wartość. Start vertex pozostaje bez zmian.")
				}
			}
			m.ConfigureGrSolver()
		case "5":
			// SA - konfiguracja parametrów
			fmt.Println("Konfiguracja SA:")
			initialTemp, err := readFloat("Podaj initialTemperature: ")
			if err != nil {
				fmt.Println("Błąd odczytu initialTemperature:", err)
				break
			}

			minTemp, err := readFloat("Podaj minimalTemperature: ")
			if err != nil {
				fmt.Println("Błąd odczytu minimalTemperature:", err)
				break
			}

			alpha, err := readFloat("Podaj alpha (współczynnik chłodzenia): ")
			if err != nil {
				fmt.Println("Błąd odczytu alpha:", err)
				break
			}

			iterations, err := readInt("Podaj liczbę iteracji na epokę: ")
			if err != nil {
				fmt.Println("Błąd odczytu iteracji:", err)
				break
			}

			timeout, err := readInt64("Podaj timeout w sekundach (-1 brak limitu): ")
			if timeout != -1 {
				timeout = utils.SecondsToNanoSeconds(timeout)
			}
			if err != nil {
				fmt.Println("Błąd odczytu timeout:", err)
				break
			}

			m.ConfigureSaSolver(initialTemp, minTemp, alpha, iterations, timeout)
		case "6":
			// TS - konfiguracja parametrów
			fmt.Println("Konfiguracja Tabu Search:")
			iterations, err := readInt("Podaj liczbę iteracji: ")
			if err != nil {
				fmt.Println("Błąd odczytu iteracji:", err)
				break
			}

			timeout, err := readInt64("Podaj timeout w sekundach (-1 brak limitu): ")
			if timeout != -1 {
				timeout = utils.SecondsToNanoSeconds(timeout)
			}
			if err != nil {
				fmt.Println("Błąd odczytu timeout:", err)
				break
			}

			tabuTenure, err := readInt("Podaj tabuTenure (ilość iteracji ruchu w tabu): ")
			if err != nil {
				fmt.Println("Błąd odczytu tabuTenure:", err)
				break
			}

			// Wybór metody sąsiedztwa
			fmt.Println("Wybierz metodę sąsiedztwa dla Tabu Search:")
			fmt.Println("1. Swap")
			fmt.Println("2. Insert")
			fmt.Print("Wybierz opcję: ")

			methodOpt, _ := reader.ReadString('\n')
			methodOpt = strings.TrimSpace(methodOpt)

			var neighborhoodMethod string
			switch methodOpt {
			case "1":
				neighborhoodMethod = ts.NeighborhoodSwap
			case "2":
				neighborhoodMethod = ts.NeighborhoodInsert
			default:
				fmt.Println("Nieznana opcja. Użyto domyślnej metody 'swap'.")
				neighborhoodMethod = ts.NeighborhoodSwap
			}

			m.ConfigureTsSolver(iterations, timeout, tabuTenure, neighborhoodMethod)
		case "b", "B":
			// Powrót do głównego menu
			return
		default:
			fmt.Println("Nieznana opcja.")
		}
	}
}

// RunInteractiveMenu uruchamia niekończącą się pętlę menu
func (m *Menu) RunInteractiveMenu() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n===== MENU =====")
		fmt.Println("1. Wczytaj graf z pliku")
		fmt.Println("2. Wygeneruj losowy graf")
		fmt.Println("3. Wyświetl aktualny graf")
		fmt.Println("4. Ustaw wierzchołek startowy")
		fmt.Println("5. Konfiguruj solvery")
		fmt.Println("6. Uruchom solvery")
		fmt.Println("7. Ustaw noEdgeValue w grafie")
		fmt.Println("q. Wyjście")
		fmt.Print("Wybierz opcję: ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		switch line {
		case "1":
			// Wczytaj graf z pliku
			fmt.Print("Podaj ścieżkę do pliku: ")
			filePath, _ := reader.ReadString('\n')
			filePath = strings.TrimSpace(filePath)

			fmt.Print("Czy to format ATSP? (y/n): ")
			atspLine, _ := reader.ReadString('\n')
			atspLine = strings.TrimSpace(atspLine)
			isATSP := (atspLine == "y" || atspLine == "Y")

			err := m.LoadGraphFromFile(filePath, isATSP)
			if err != nil {
				fmt.Println("Błąd wczytywania grafu:", err)
			} else {
				fmt.Println("Graf wczytany pomyślnie.")
			}
		case "2":
			// Wygeneruj losowy graf
			fmt.Print("Podaj liczbę wierzchołków: ")
			vcStr, _ := reader.ReadString('\n')
			vcStr = strings.TrimSpace(vcStr)
			vc, err := strconv.Atoi(vcStr)
			if err != nil || vc <= 0 {
				fmt.Println("Nieprawidłowa liczba wierzchołków.")
				break
			}

			fmt.Print("Podaj wartość dla braku krawędzi: ")
			neStr, _ := reader.ReadString('\n')
			neStr = strings.TrimSpace(neStr)
			neVal, err := strconv.Atoi(neStr)
			if err != nil {
				fmt.Println("Nieprawidłowa wartość dla braku krawędzi.")
				break
			}

			fmt.Print("Podaj maksymalną wagę krawędzi: ")
			mwStr, _ := reader.ReadString('\n')
			mwStr = strings.TrimSpace(mwStr)
			mw, err := strconv.Atoi(mwStr)
			if err != nil || mw <= 0 {
				fmt.Println("Nieprawidłowa maksymalna waga krawędzi.")
				break
			}

			m.GenerateRandomGraph(vc, neVal, mw)
			fmt.Println("Wygenerowano losowy graf.")
		case "3":
			// Wyświetl aktualny graf
			m.DisplayGraph()
		case "4":
			// Ustaw wierzchołek startowy
			fmt.Print("Podaj wierzchołek startowy: ")
			svStr, _ := reader.ReadString('\n')
			svStr = strings.TrimSpace(svStr)
			sv, err := strconv.Atoi(svStr)
			if err != nil {
				fmt.Println("Nieprawidłowa wartość wierzchołka startowego.")
				break
			}
			if m.graph != nil && (sv < 0 || sv >= m.graph.GetVertexCount()) {
				fmt.Println("Wierzchołek startowy poza zakresem.")
				break
			}
			m.SetStartVertex(sv)
			fmt.Println("Ustawiono wierzchołek startowy na:", sv)
		case "5":
			// Konfiguruj solvery
			m.solverConfigurationSubmenu()
		case "6":
			// Uruchom solvery
			if m.graph == nil {
				fmt.Println("Najpierw wczytaj lub wygeneruj graf.")
				break
			}

			fmt.Println("\nWybierz solver do uruchomienia:")
			fmt.Println("1. Brute Force")
			fmt.Println("2. Branch and Bound")
			fmt.Println("3. Dynamic Programming")
			fmt.Println("4. Greedy")
			fmt.Println("5. Simulated Annealing")
			fmt.Println("6. Tabu Search")
			fmt.Print("Wybierz solver: ")

			solverOpt, _ := reader.ReadString('\n')
			solverOpt = strings.TrimSpace(solverOpt)

			switch solverOpt {
			case "1":
				m.RunBf()
			case "2":
				m.RunBnb()
			case "3":
				m.RunDp()
			case "4":
				m.RunGr()
			case "5":
				m.RunSa()
			case "6":
				m.RunTs()
			default:
				fmt.Println("Nieznana opcja.")
			}
		case "7":
			// Ustaw noEdgeValue w grafie
			if m.graph == nil {
				fmt.Println("Najpierw wczytaj lub wygeneruj graf.")
				break
			}
			fmt.Print("Podaj nową wartość dla braku krawędzi: ")
			nevStr, _ := reader.ReadString('\n')
			nevStr = strings.TrimSpace(nevStr)
			nev, err := strconv.Atoi(nevStr)
			if err != nil {
				fmt.Println("Nieprawidłowa wartość.")
				break
			}
			m.SetNoEdgeValue(nev)
			fmt.Println("Ustawiono noEdgeValue na:", nev)
		case "q", "Q":
			// Wyjście z menu
			fmt.Println("Zakończono działanie programu.")
			return
		default:
			fmt.Println("Nieznana opcja, spróbuj ponownie.")
		}
	}
}

// Helper functions to read input
func readFloat(prompt string) (float64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("brak wartości")
	}
	return strconv.ParseFloat(line, 64)
}

func readInt(prompt string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("brak wartości")
	}
	return strconv.Atoi(line)
}

func readInt64(prompt string) (int64, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)
	if line == "" {
		return 0, fmt.Errorf("brak wartości")
	}
	return strconv.ParseInt(line, 10, 64)
}

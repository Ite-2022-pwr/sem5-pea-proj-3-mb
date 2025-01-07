package graph

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"
)

// LoadGraphFromFile wczytuje graf z pliku tekstowego o podanej ścieżce.
// Jeśli isATSP = true, próbuje wczytać plik w formacie ATSP (EXPLICIT, FULL_MATRIX).
// W przeciwnym wypadku oczekuje formatu:
// - w pierwszej linii: liczba wierzchołków
// - w kolejnych liniach: macierz n x n
// Funkcja poprawnie interpretuje wielokrotne spacje jako separator.
func LoadGraphFromFile(filePath string, graph *AdjMatrixGraph, isATSP ...bool) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if len(isATSP) > 0 && isATSP[0] {
		// Wczytywanie formatu ATSP
		var dimension int
		inMatrixSection := false

		// Tablica do której wczytamy wszystkie liczby z sekcji EDGE_WEIGHT_SECTION
		var allNumbers []int

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())

			// Szukamy wymiaru grafu w linii zawierającej "DIMENSION:"
			if strings.HasPrefix(line, "DIMENSION:") {
				parts := strings.Split(line, ":")
				if len(parts) == 2 {
					dimStr := strings.TrimSpace(parts[1])
					dimension, err = strconv.Atoi(dimStr)
					if err != nil {
						return errors.New("błąd przy konwersji wymiaru z pliku ATSP")
					}
				}
			}

			// Szukamy sekcji z wagami
			if strings.HasPrefix(line, "EDGE_WEIGHT_SECTION") {
				inMatrixSection = true
				continue
			}

			// Jeśli jesteśmy w sekcji z macierzą
			if inMatrixSection {
				// Wczytujemy wszystkie liczby do jednej tablicy
				if line != "" {
					values := strings.Fields(line)
					for _, val := range values {
						num, err := strconv.Atoi(val)
						if err != nil {
							return errors.New("błąd konwersji wartości w macierzy ATSP")
						}
						allNumbers = append(allNumbers, num)
					}
				}

				// Jeżeli mamy już wystarczającą liczbę elementów do zapełnienia macierzy dimension x dimension, możemy przerwać
				if dimension > 0 && len(allNumbers) >= dimension*dimension {
					break
				}
			}
		}

		if dimension == 0 {
			return errors.New("nie znaleziono wymiaru w pliku ATSP")
		}
		if len(allNumbers) < dimension*dimension {
			return errors.New("zbyt mało danych w sekcji EDGE_WEIGHT_SECTION aby uzupełnić macierz")
		}

		// Inicjalizacja grafu
		graph.vertexCount = dimension
		graph.adjMatrix = make([][]int, dimension)
		for i := 0; i < dimension; i++ {
			graph.adjMatrix[i] = make([]int, dimension)
		}

		// Uzupełnienie macierzy z tablicy allNumbers
		idx := 0
		for i := 0; i < dimension; i++ {
			for j := 0; j < dimension; j++ {
				graph.adjMatrix[i][j] = allNumbers[idx]
				idx++
			}
		}

		return nil
	} else {
		// Wczytywanie standardowego formatu
		if !scanner.Scan() {
			return errors.New("plik jest pusty lub nieprawidłowy")
		}
		firstLine := strings.TrimSpace(scanner.Text())
		vertexCount, err := strconv.Atoi(firstLine)
		if err != nil {
			return errors.New("błąd podczas odczytu liczby wierzchołków")
		}

		graph.vertexCount = vertexCount
		graph.adjMatrix = make([][]int, vertexCount)
		for i := 0; i < vertexCount; i++ {
			graph.adjMatrix[i] = make([]int, vertexCount)
		}

		row := 0
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				continue
			}
			values := strings.Fields(line)
			if len(values) != vertexCount {
				return errors.New("niewłaściwa liczba elementów w wierszu macierzy")
			}
			for j, val := range values {
				num, err := strconv.Atoi(val)
				if err != nil {
					return errors.New("błąd przy konwersji wartości macierzy do liczby całkowitej")
				}
				graph.adjMatrix[row][j] = num
			}
			row++
			if row == vertexCount {
				break
			}
		}

		if row != vertexCount {
			return errors.New("niewłaściwa liczba wierszy w macierzy sąsiedztwa")
		}
		return nil
	}
}

// SaveGraphToFile zapisuje graf do pliku w formacie:
// Pierwsza linia: liczba wierzchołków
// Kolejne linie: macierz n x n
// Jeśli useTabsAsSeparator = true, wartości oddzielone tabulatorami, w przeciwnym wypadku spacjami.
func SaveGraphToFile(g Graph, filePath string, useTabsAsSeparator ...bool) error {
	separator := " "
	if len(useTabsAsSeparator) > 0 && useTabsAsSeparator[0] {
		separator = "\t"
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var out strings.Builder
	out.WriteString(strconv.Itoa(g.GetVertexCount()) + "\n")

	for i := 0; i < g.GetVertexCount(); i++ {
		for j := 0; j < g.GetVertexCount(); j++ {
			edge := g.GetEdge(i, j)
			out.WriteString(strconv.Itoa(edge.Weight))
			if j < g.GetVertexCount()-1 {
				out.WriteString(separator)
			}
		}
		out.WriteString("\n")
	}

	_, err = file.WriteString(out.String())
	if err != nil {
		return err
	}

	return nil
}

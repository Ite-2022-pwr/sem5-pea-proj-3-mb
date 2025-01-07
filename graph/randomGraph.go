package graph

import (
	"math/rand"
	"time"
)

// GenerateRandomGraph generuje losowy graf z daną liczbą wierzchołków i wypełnia krawędzie losowymi wagami większymi niż zero.
// Wartość `noEdgeValue` jest przypisywana tam, gdzie krawędź nie istnieje (między wierzchołkiem a samym sobą).
// `maxWeight` - maksymalna wartość wag krawędzi (losowane wartości będą z zakresu od 1 do maxWeight).
func GenerateRandomGraph(g *AdjMatrixGraph, vertexCount int, noEdgeValue int, maxWeight int) {
	// Inicjalizuj licznik losowy za pomocą obecnego czasu
	rand.Seed(time.Now().UnixNano())

	// Ustaw liczbę wierzchołków i inicjalizuj macierz sąsiedztwa
	g.vertexCount = vertexCount
	g.noEdgeValue = noEdgeValue
	g.adjMatrix = make([][]int, vertexCount)
	for i := 0; i < vertexCount; i++ {
		g.adjMatrix[i] = make([]int, vertexCount)
	}

	// Przejdź przez wszystkie pary wierzchołków i generuj losowe wagi
	for i := 0; i < vertexCount; i++ {
		for j := 0; j < vertexCount; j++ {
			if i == j {
				// Brak krawędzi do samego siebie
				g.adjMatrix[i][j] = noEdgeValue
			} else {
				// Generowanie losowej wagi krawędzi większej niż 0
				g.adjMatrix[i][j] = rand.Intn(maxWeight) + 1 // Losowe liczby od 1 do maxWeight
			}
		}
	}
}

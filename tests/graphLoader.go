package tests

import (
	"log"
	"projekt2/graph"
)

func LoadTestGraphs() (graph.Graph, graph.Graph, graph.Graph) {
	smallGraph := graph.NewAdjMatrixGraph(56, 100000000)
	mediumGraph := graph.NewAdjMatrixGraph(171, 100000000)
	largeGraph := graph.NewAdjMatrixGraph(358, 0)
	errS := graph.LoadGraphFromFile("ftv55.atsp", smallGraph, true)
	errM := graph.LoadGraphFromFile("ftv170.atsp", mediumGraph, true)
	errG := graph.LoadGraphFromFile("rbg358.atsp", largeGraph, true)
	if errS != nil || errM != nil || errG != nil {
		log.Println(errS)
		log.Println(errM)
		log.Println(errG)
		log.Fatal("Błąd wczytywania grafów")
		return nil, nil, nil
	}
	return smallGraph, mediumGraph, largeGraph
}

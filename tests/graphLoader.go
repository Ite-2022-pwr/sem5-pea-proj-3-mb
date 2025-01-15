package tests

import (
	"log"
	"projekt3/graph"
)

func LoadTestGraphs() (graph.Graph, graph.Graph, graph.Graph, graph.Graph) {
	tinyGraph := graph.NewAdjMatrixGraph(17, 9999)
	smallGraph := graph.NewAdjMatrixGraph(56, 100000000)
	mediumGraph := graph.NewAdjMatrixGraph(171, 100000000)
	largeGraph := graph.NewAdjMatrixGraph(358, 0)
	errT := graph.LoadGraphFromFile("br17.atsp", tinyGraph, true)
	errS := graph.LoadGraphFromFile("ftv55.atsp", smallGraph, true)
	errM := graph.LoadGraphFromFile("ftv170.atsp", mediumGraph, true)
	errG := graph.LoadGraphFromFile("rbg358.atsp", largeGraph, true)
	if errT != nil || errS != nil || errM != nil || errG != nil {
		log.Println(errT)
		log.Println(errS)
		log.Println(errM)
		log.Println(errG)
		log.Fatal("Błąd wczytywania grafów")
		return nil, nil, nil, nil
	}
	return tinyGraph, smallGraph, mediumGraph, largeGraph
}

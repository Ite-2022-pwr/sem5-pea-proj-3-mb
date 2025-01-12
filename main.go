package main

import (
	"fmt"
	"log"
	"math"
	"projekt2/graph"
	"projekt2/solver/aco"
	"time"
)

func main() {

	//runInteractiveMenuPTR := flag.Bool("interactive", true, "Run interactive menu(default true)")
	//flag.Parse()
	//
	//if *runInteractiveMenuPTR {
	//	mainMenu := menu.NewMenu()
	//	mainMenu.RunInteractiveMenu()
	//} else {
	//	tuningTests.RunTuning()
	//	amountTests.RunAmountTests()
	//	tests.RunOptimalSA()
	//	tests.RunOptimalTS()
	//}
	g := graph.NewAdjMatrixGraph(443, 0)
	err := graph.LoadGraphFromFile("rbg443.atsp", g, true)
	if err != nil {
		return
	}
	fmt.Println(g.ToString())
	ants := g.GetVertexCount()
	bestGreedyWeight := math.MaxInt
	for i := 0; i < ants; i++ {
		greedyWeight := g.CalculatePathWeight(g.GetHamiltonianPathGreedy(i))
		if greedyWeight < bestGreedyWeight {
			bestGreedyWeight = greedyWeight
		}
	}
	startPheromone := float64(bestGreedyWeight) / float64(ants)
	pheromonesPerAnt := float64(bestGreedyWeight)
	fmt.Println(pheromonesPerAnt)
	acoSolv := aco.NewACOZeroEdgeSolver(100, 100000, 10000, 1.0, 5.0, 0.6, pheromonesPerAnt, startPheromone, 60)
	acoSolv.SetGraph(g)
	startTime := time.Now()
	log.Println("Start time: ", startTime)
	path, cost := acoSolv.Solve()
	elapsedTime := time.Since(startTime)
	log.Println("Elapsed time: ", elapsedTime)
	log.Println(path, "\n", cost)
}

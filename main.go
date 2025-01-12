package main

import (
	"fmt"
	"log"
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
	g := graph.NewAdjMatrixGraph(17, 9999)
	err := graph.LoadGraphFromFile("br17.atsp", g, true)
	if err != nil {
		return
	}
	fmt.Println(g.ToString())
	ants := g.GetVertexCount()
	greedyWeight := g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0))
	startPheromone := float64(greedyWeight) / float64(ants)
	pheromonesPerAnt := float64(greedyWeight)
	fmt.Println(pheromonesPerAnt)
	acoSolv := aco.NewACOZeroEdgeSolver(ants, 100, 300, 1.0, 4.0, 0.5, pheromonesPerAnt, startPheromone, 30)
	acoSolv.SetGraph(g)
	startTime := time.Now()
	log.Println("Start time: ", startTime)
	path, cost := acoSolv.Solve()
	elapsedTime := time.Since(startTime)
	log.Println("Elapsed time: ", elapsedTime)
	log.Println(path, "\n", cost)
	fmt.Println(acoSolv.PheromonesToString())
}

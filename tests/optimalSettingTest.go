package tests

import (
	"log"
	"projekt2/graph"
	"projekt2/solver/sa"
	"projekt2/solver/ts"
	"projekt2/utils"
	"time"
)

func RunOptimalTS() {
	smallGraph, mediumGraph, largeGraph := LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphTS(smallGraph, timeoutInNs, "ts_optimal_small_")
	runSingleGraphTS(mediumGraph, timeoutInNs, "ts_optimal_medium_")
	runSingleGraphTS(largeGraph, timeoutInNs, "ts_optimal_large_")
}

func RunOptimalSA() {
	smallGraph, mediumGraph, largeGraph := LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphSA(smallGraph, timeoutInNs, "sa_optimal_small_")
	runSingleGraphSA(mediumGraph, timeoutInNs, "sa_optimal_medium_")
	runSingleGraphSA(largeGraph, timeoutInNs, "sa_optimal_large_")
}

func runSingleGraphTS(g graph.Graph, timeoutInNs int64, fileOutName string) {
	results := make([][]int64, 2)
	for i := 0; i < 2; i++ {
		results[i] = make([]int64, 10)
	}

	tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, 10, "insert")
	tsSolver.SetGraph(g)
	tsSolver.SetStartVertex(0)
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, weight := tsSolver.Solve()
		elapsed := time.Since(start)
		log.Println(" Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
		results[0][i] = elapsed.Nanoseconds()
		results[1][i] = int64(weight)
	}
	utils.SaveTimesToCSVFile(results, fileOutName+utils.GetDateForFilename()+".csv")

}

func runSingleGraphSA(g graph.Graph, timeoutInNs int64, fileOutName string) {
	results := make([][]int64, 2)
	for i := 0; i < 2; i++ {
		results[i] = make([]int64, 10)
	}

	saSolver := sa.NewSimulatedAnnealingATSPSolver(1000000, 1e-9, 0.995, 5000, timeoutInNs)
	saSolver.SetGraph(g)
	saSolver.SetStartVertex(0)
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, weight := saSolver.Solve()
		elapsed := time.Since(start)
		log.Println(" Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
		results[0][i] = elapsed.Nanoseconds()
		results[1][i] = int64(weight)
	}
	utils.SaveTimesToCSVFile(results, fileOutName+utils.GetDateForFilename()+".csv")
}

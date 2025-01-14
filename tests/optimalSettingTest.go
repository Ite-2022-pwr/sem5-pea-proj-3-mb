package tests

import (
	"log"
	"math"
	"projekt2/graph"
	"projekt2/solver/aco"
	"projekt2/utils"
	"time"
)

func RunOptimalACO() {
	tinyGraph, smallGraph, mediumGraph, largeGraph := LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	runSingleGraphACO(tinyGraph, timeoutInNs, "aco_optimal_tiny_")
	runSingleGraphACO(smallGraph, timeoutInNs, "aco_optimal_small_")
	runSingleGraphACO(mediumGraph, timeoutInNs, "aco_optimal_medium_")
	runSingleGraphACO(largeGraph, timeoutInNs, "aco_optimal_large_")
}

func runSingleGraphACO(g graph.Graph, timeoutInNs int64, fileOutName string) {
	results := make([][]int64, 2)
	for i := 0; i < 2; i++ {
		results[i] = make([]int64, 10)
	}

	acoSolver := aco.NewACOZeroEdgeSolver(100, 100, math.MaxInt, 2.0, 2.0, 0.5, 30.0, float64(g.GetVertexCount())/float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0))), timeoutInNs)
	acoSolver.SetGraph(g)
	for i := 0; i < 10; i++ {
		start := time.Now()
		_, weight := acoSolver.Solve()
		elapsed := time.Since(start)
		log.Println(" Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
		results[0][i] = elapsed.Nanoseconds()
		results[1][i] = int64(weight)
	}
	utils.SaveTimesToCSVFile(results, fileOutName+utils.GetDateForFilename()+".csv")

}

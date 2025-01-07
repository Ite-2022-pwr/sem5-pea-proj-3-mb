package saTuning

import (
	"log"
	"projekt2/graph"
	"projekt2/solver/sa"

	"projekt2/tests"
	"projekt2/utils"
	"strconv"
	"time"
)

func RunIterationsTuningSA() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	iterations := []int{1000, 3000, 5000}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphIterTuning(smallGraph, iterations, timeoutInNs, "sa_iterations_small_")
	runSingleGraphIterTuning(mediumGraph, iterations, timeoutInNs, "sa_iterations_medium_")
	runSingleGraphIterTuning(largeGraph, iterations, timeoutInNs, "sa_iterations_large_")
}

func runSingleGraphIterTuning(g graph.Graph, iterations []int, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(iterations))
	for i, _ := range iterations {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, it := range iterations {
		saSolver := sa.NewSimulatedAnnealingATSPSolver(10000, 1e-9, 0.995, it, timeoutInNs)
		saSolver.SetGraph(g)
		saSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := saSolver.Solve()
			elapsed := time.Since(start)
			log.Println("Iteration: ", it, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.Itoa(it)+"iterations_"+utils.GetDateForFilename()+".csv")
	}
}

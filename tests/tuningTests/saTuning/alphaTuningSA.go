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

func RunAlphaTuningSA() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	alphas := []float64{0.9, 0.95, 0.99, 0.995}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphAlphaTuning(smallGraph, alphas, timeoutInNs, "sa_alpha_small_")
	runSingleGraphAlphaTuning(mediumGraph, alphas, timeoutInNs, "sa_alpha_medium_")
	runSingleGraphAlphaTuning(largeGraph, alphas, timeoutInNs, "sa_alpha_large_")
}

func runSingleGraphAlphaTuning(g graph.Graph, alphas []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(alphas))
	for i, _ := range alphas {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, alpha := range alphas {
		saSolver := sa.NewSimulatedAnnealingATSPSolver(10000, 1e-9, alpha, 5000, timeoutInNs)
		saSolver.SetGraph(g)
		saSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := saSolver.Solve()
			elapsed := time.Since(start)
			log.Println("Alpha: ", alpha, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(alpha, 'E', -1, 64)+"min_temp_"+utils.GetDateForFilename()+".csv")
	}
}

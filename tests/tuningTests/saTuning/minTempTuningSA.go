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

func RunMinTempTuningSA() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	minTemps := []float64{1, 0.001, 0.000001, 0.000000001}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphMinTempTuning(smallGraph, minTemps, timeoutInNs, "sa_min_temp_small_")
	runSingleGraphMinTempTuning(mediumGraph, minTemps, timeoutInNs, "sa_min_temp_medium_")
	runSingleGraphMinTempTuning(largeGraph, minTemps, timeoutInNs, "sa_min_temp_large_")
}

func runSingleGraphMinTempTuning(g graph.Graph, minTemps []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(minTemps))
	for i, _ := range minTemps {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, minTemp := range minTemps {
		saSolver := sa.NewSimulatedAnnealingATSPSolver(10000, minTemp, 0.995, 5000, timeoutInNs)
		saSolver.SetGraph(g)
		saSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := saSolver.Solve()
			elapsed := time.Since(start)
			log.Println("MinTemp: ", minTemp, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(minTemp, 'E', -1, 64)+"min_temp_"+utils.GetDateForFilename()+".csv")
	}
}

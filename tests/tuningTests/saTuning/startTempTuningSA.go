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

func RunInitialTempTuningSA() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	initTemps := []float64{1000, 1000000, 1000000000}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphInitialTempTuning(smallGraph, initTemps, timeoutInNs, "sa_init_temp_small_")
	runSingleGraphInitialTempTuning(mediumGraph, initTemps, timeoutInNs, "sa_init_temp_medium_")
	runSingleGraphInitialTempTuning(largeGraph, initTemps, timeoutInNs, "sa_init_temp_large_")
}

func runSingleGraphInitialTempTuning(g graph.Graph, initTemps []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(initTemps))
	for i, _ := range initTemps {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, initTemp := range initTemps {
		saSolver := sa.NewSimulatedAnnealingATSPSolver(initTemp, 0.000000001, 0.995, 5000, timeoutInNs)
		saSolver.SetGraph(g)
		saSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := saSolver.Solve()
			elapsed := time.Since(start)
			log.Println("InitTemp: ", initTemp, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(initTemp, 'E', -1, 64)+"min_temp_"+utils.GetDateForFilename()+".csv")
	}
}

package tsTuning

import (
	"log"
	"projekt2/graph"
	"projekt2/solver/ts"
	"projekt2/tests"
	"projekt2/utils"
	"strconv"
	"time"
)

func RunTenureTuningTS() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	tenures := []int{5, 10, 15}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphTenureTuning(smallGraph, tenures, timeoutInNs, "ts_tenure_small_")
	runSingleGraphTenureTuning(mediumGraph, tenures, timeoutInNs, "ts_tenure_medium_")
	runSingleGraphTenureTuning(largeGraph, tenures, timeoutInNs, "ts_tenure_large_")
}

func runSingleGraphTenureTuning(g graph.Graph, tenures []int, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(tenures))
	for i, _ := range tenures {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, ten := range tenures {
		tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, ten, "insert")
		tsSolver.SetGraph(g)
		tsSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := tsSolver.Solve()
			elapsed := time.Since(start)
			log.Println("Tenure: ", ten, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.Itoa(ten)+"tenure_"+utils.GetDateForFilename()+".csv")
	}
}

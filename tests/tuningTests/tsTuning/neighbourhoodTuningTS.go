package tsTuning

import (
	"log"
	"projekt2/graph"
	"projekt2/solver/ts"
	"projekt2/tests"
	"projekt2/utils"
	"time"
)

func RunNeighbourTuningTS() {
	smallGraph, mediumGraph, largeGraph := tests.LoadTestGraphs()
	neighbours := []string{"swap", "insert"}
	timeoutInNs := utils.MinutesToNanoSeconds(5)
	runSingleGraphNeighTuning(smallGraph, neighbours, timeoutInNs, "ts_neighbours_small_")
	runSingleGraphNeighTuning(mediumGraph, neighbours, timeoutInNs, "ts_neighbours_medium_")
	runSingleGraphNeighTuning(largeGraph, neighbours, timeoutInNs, "ts_neighbours_large_")
}

func runSingleGraphNeighTuning(g graph.Graph, neighbours []string, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(neighbours))
	for i, _ := range neighbours {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, neigh := range neighbours {
		tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, 10, neigh)
		tsSolver.SetGraph(g)
		tsSolver.SetStartVertex(0)
		for j := 0; j < 10; j++ {
			start := time.Now()
			_, weight := tsSolver.Solve()
			elapsed := time.Since(start)
			log.Println("Neighbourhood: ", neigh, " Time: ", elapsed, " Weight: ", weight, " Graph size: ", g.GetVertexCount())
			results[i][0][j] = elapsed.Nanoseconds()
			results[i][1][j] = int64(weight)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+neigh+"neighbour_"+utils.GetDateForFilename()+".csv")
	}
}

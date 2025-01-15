package acoTuning

import (
	"log"
	"math"
	"projekt3/graph"
	"projekt3/solver/aco"
	"projekt3/tests"
	"projekt3/utils"
	"strconv"
	"time"
)

func RunAntCountTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	antCountsT := []int{10, 50, 100, 17}
	antCountsS := []int{10, 50, 100, 56}
	antCountsM := []int{10, 50, 100, 171}
	antCountsL := []int{10, 50, 100, 358}
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	runSingleGraphAntCountTuning(tinyG, antCountsT, timeoutInNs, "aco_ant_count_tiny_")
	runSingleGraphAntCountTuning(smallG, antCountsS, timeoutInNs, "aco_ant_count_small_")
	runSingleGraphAntCountTuning(mediumG, antCountsM, timeoutInNs, "aco_ant_count_medium_")
	runSingleGraphAntCountTuning(largeG, antCountsL, timeoutInNs, "aco_ant_count_large_")
}

func runSingleGraphAntCountTuning(g graph.Graph, antCounts []int, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(antCounts))
	iterations := 100
	alpha := 1.0
	beta := 5.0
	rho := 0.5
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0

	for i, _ := range antCounts {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, antCount := range antCounts {
		acoSolver := aco.NewACOZeroEdgeSolver(antCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Ant count: ", antCount, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.Itoa(antCount)+"_"+utils.GetDateForFilename()+".csv")
	}

}

package acoTuning

import (
	"log"
	"math"
	"projekt2/graph"
	"projekt2/solver/aco"
	"projekt2/tests"
	"projekt2/utils"
	"strconv"
	"time"
)

func RunBetaTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	betaValues := []float64{2.0, 5.0, 8.0}
	runSingleGraphBetaTuning(tinyG, betaValues, timeoutInNs, "aco_beta_tiny_")
	runSingleGraphBetaTuning(smallG, betaValues, timeoutInNs, "aco_beta_small_")
	runSingleGraphBetaTuning(mediumG, betaValues, timeoutInNs, "aco_beta_medium_")
	runSingleGraphBetaTuning(largeG, betaValues, timeoutInNs, "aco_beta_large_")
}

func runSingleGraphBetaTuning(g graph.Graph, betaValues []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(betaValues))
	iterations := 100
	antsCount := 30
	alpha := 1.0
	rho := 0.5
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0

	for i, _ := range betaValues {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, beta := range betaValues {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Beta: ", beta, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(beta, 'E', -1, 64)+"_"+utils.GetDateForFilename()+".csv")
	}
}

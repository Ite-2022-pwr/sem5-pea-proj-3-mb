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

func RunEvaporationRateTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	evaporationRates := []float64{0.1, 0.3, 0.5}
	timeoutInNs := utils.MinutesToNanoSeconds(2)
	runSingleGraphEvaporationRateTuning(tinyG, evaporationRates, timeoutInNs, "aco_evaporation_tiny_")
	runSingleGraphEvaporationRateTuning(smallG, evaporationRates, timeoutInNs, "aco_evaporation_small_")
	runSingleGraphEvaporationRateTuning(mediumG, evaporationRates, timeoutInNs, "aco_evaporation_medium_")
	runSingleGraphEvaporationRateTuning(largeG, evaporationRates, timeoutInNs, "aco_evaporation_large_")
}

func runSingleGraphEvaporationRateTuning(g graph.Graph, evaporationRates []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(evaporationRates))
	iterations := 100
	antsCount := g.GetVertexCount()
	alpha := 1.0
	beta := 5.0
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0

	for i, _ := range evaporationRates {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, evaporationRate := range evaporationRates {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, evaporationRate, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Evaporation rate: ", evaporationRate, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(evaporationRate, 'E', -1, 64)+"_"+utils.GetDateForFilename()+".csv")
	}
}

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

func RunStartPheromonesTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	startPheromones := []float64{1, 10}
	tinyGPheromones := float64(tinyG.GetVertexCount()) / float64(tinyG.CalculatePathWeight(tinyG.GetHamiltonianPathGreedy(0)))
	smallGPheromones := float64(smallG.GetVertexCount()) / float64(smallG.CalculatePathWeight(smallG.GetHamiltonianPathGreedy(0)))
	mediumGPheromones := float64(mediumG.GetVertexCount()) / float64(mediumG.CalculatePathWeight(mediumG.GetHamiltonianPathGreedy(0)))
	largeGPheromones := float64(largeG.GetVertexCount()) / float64(largeG.CalculatePathWeight(largeG.GetHamiltonianPathGreedy(0)))
	startPheromonesT := append(startPheromones, tinyGPheromones)
	startPheromonesS := append(startPheromones, smallGPheromones)
	startPheromonesM := append(startPheromones, mediumGPheromones)
	startPheromonesL := append(startPheromones, largeGPheromones)
	runSingleGraphStartPheromonesTuning(tinyG, startPheromonesT, timeoutInNs, "aco_start_pheromones_tiny_")
	runSingleGraphStartPheromonesTuning(smallG, startPheromonesS, timeoutInNs, "aco_start_pheromones_small_")
	runSingleGraphStartPheromonesTuning(mediumG, startPheromonesM, timeoutInNs, "aco_start_pheromones_medium_")
	runSingleGraphStartPheromonesTuning(largeG, startPheromonesL, timeoutInNs, "aco_start_pheromones_large_")
}

func runSingleGraphStartPheromonesTuning(g graph.Graph, startPheromones []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(startPheromones))
	iterations := 100
	antsCount := 30
	alpha := 1.0
	beta := 5.0
	rho := 0.5
	pheromonesPerAnt := 1.0

	for i, _ := range startPheromones {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, startPheromone := range startPheromones {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Start pheromone: ", startPheromone, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(startPheromone, 'E', -1, 64)+"_"+utils.GetDateForFilename()+".csv")
	}
}

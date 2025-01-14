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

func RunIterTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	iterationsToTest := []int{10, 50, 100}
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	runSingleGraphIterTuning(tinyG, iterationsToTest, timeoutInNs, "aco_iter_tiny_")
	runSingleGraphIterTuning(smallG, iterationsToTest, timeoutInNs, "aco_iter_small_")
	runSingleGraphIterTuning(mediumG, iterationsToTest, timeoutInNs, "aco_iter_medium_")
	runSingleGraphIterTuning(largeG, iterationsToTest, timeoutInNs, "aco_iter_large_")
}

func runSingleGraphIterTuning(g graph.Graph, iterations []int, timeoutInNs int64, fileOutName string) {
	antsCount := 30
	results := make([][][]int64, len(iterations))
	alpha := 1.0
	beta := 1.0
	rho := 0.5 // evaporation rate
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0

	for i, _ := range iterations {
		//0 - time, 1 - cost
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, iteration := range iterations {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iteration, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Iteration: ", iteration, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.Itoa(iteration)+"_"+utils.GetDateForFilename()+".csv")
	}
}

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

func RunAlphaTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	alphaValues := []float64{1.0, 2.0, 3.0, 5.0}
	runSingleGraphAlphaTuning(tinyG, alphaValues, timeoutInNs, "aco_alpha_tiny_")
	runSingleGraphAlphaTuning(smallG, alphaValues, timeoutInNs, "aco_alpha_small_")
	runSingleGraphAlphaTuning(mediumG, alphaValues, timeoutInNs, "aco_alpha_medium_")
	runSingleGraphAlphaTuning(largeG, alphaValues, timeoutInNs, "aco_alpha_large_")
}

func runSingleGraphAlphaTuning(g graph.Graph, alphaValues []float64, timeoutInNs int64, fileOutName string) {
	results := make([][][]int64, len(alphaValues))
	iterations := 100
	antsCount := 30
	beta := 5.0
	rho := 0.5
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0

	for i, _ := range alphaValues {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, alpha := range alphaValues {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Alpha: ", alpha, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(alpha, 'E', -1, 64)+"_"+utils.GetDateForFilename()+".csv")
	}
}

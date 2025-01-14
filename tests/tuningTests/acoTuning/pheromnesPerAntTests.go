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

func RunPheromonesPerAntTests() {
	tinyG, smallG, mediumG, bigG := tests.LoadTestGraphs()
	pheromonesPerAnts := []float64{1.0, 30.0, 50.0, 100.0}
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	runSingleGraphPheromonesPerAntTuning(tinyG, pheromonesPerAnts, timeoutInNs, "aco_pheromones_per_ant_tiny_")
	runSingleGraphPheromonesPerAntTuning(smallG, pheromonesPerAnts, timeoutInNs, "aco_pheromones_per_ant_small_")
	runSingleGraphPheromonesPerAntTuning(mediumG, pheromonesPerAnts, timeoutInNs, "aco_pheromones_per_ant_medium_")
	runSingleGraphPheromonesPerAntTuning(bigG, pheromonesPerAnts, timeoutInNs, "aco_pheromones_per_ant_big_")

}

func runSingleGraphPheromonesPerAntTuning(g graph.Graph, pheromonesPerAnts []float64, timeOutNs int64, fileOutName string) {
	results := make([][][]int64, len(pheromonesPerAnts))
	iterations := 100
	antsCount := 30
	alpha := 1.0
	beta := 5.0
	rho := 0.5
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))

	for i, _ := range pheromonesPerAnts {
		results[i] = make([][]int64, 2)
		for j := 0; j < 2; j++ {
			results[i][j] = make([]int64, 10)
		}
	}

	for i, pheromonesPerAnt := range pheromonesPerAnts {
		acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeOutNs)
		acoSolver.SetGraph(g)
		for j := 0; j < 10; j++ {
			startTime := time.Now()
			path, cost := acoSolver.Solve()
			elapsedTime := time.Since(startTime)
			log.Println("Pheromones per ant: ", pheromonesPerAnt, " Time: ", elapsedTime, " Weight: ", cost, " Graph size: ", g.GetVertexCount(), "\n", g.PathWithWeightsToString(path))
			results[i][0][j] = elapsedTime.Nanoseconds()
			results[i][1][j] = int64(cost)
		}
		utils.SaveTimesToCSVFile(results[i], fileOutName+strconv.FormatFloat(pheromonesPerAnt, 'E', -1, 64)+"_"+utils.GetDateForFilename()+".csv")
	}
}

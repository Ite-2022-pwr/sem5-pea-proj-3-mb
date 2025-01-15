package comparisonTests

import (
	"math"
	"projekt3/graph"
	"projekt3/solver/aco"
	"projekt3/solver/bnb"
	"projekt3/solver/dp"
	"projekt3/solver/sa"
	"projekt3/solver/ts"
	"projekt3/tests"
	"projekt3/utils"
	"time"
)

func RunComparisonTests() {
	tinyG, smallG, mediumG, largeG := tests.LoadTestGraphs()
	timeoutInNs := utils.MinutesToNanoSeconds(2)
	runSingleGraphComparisonTest(smallG, timeoutInNs, "comparison_small_")
	runSingleGraphComparisonTest(mediumG, timeoutInNs, "comparison_medium_")
	runSingleGraphComparisonTest(largeG, timeoutInNs, "comparison_large_")
	runTinyComparisonTests(tinyG, "comparison_br17_")
}

func runTinyComparisonTests(g graph.Graph, filename string) {
	timeoutInNs := utils.MinutesToNanoSeconds(1)
	results := make([][]int64, 8) // aco, sa, ts, dp, bnb
	for i := 0; i < 8; i++ {
		results[i] = make([]int64, 10)
	}

	for i := 0; i < 10; i++ {
		results[0][i], results[1][i] = runACO(g, timeoutInNs)
		results[2][i], results[3][i] = runSA(g, timeoutInNs)
		results[4][i], results[5][i] = runTS(g, timeoutInNs)
		results[6][i], _ = runDp(g)
		results[7][i], _ = runBnB(g)
	}
	utils.SaveTimesToCSVFile(results, filename+utils.GetDateForFilename()+".csv")
}

func runSingleGraphComparisonTest(g graph.Graph, timeoutInNs int64, fileOutName string) {
	results := make([][]int64, 6) // col 0,1 aco; col 2,3 sa; col 4,5 ts
	for i := 0; i < 6; i++ {
		results[i] = make([]int64, 10)
	}
	for i := 0; i < 10; i++ {
		results[0][i], results[1][i] = runACO(g, timeoutInNs)
		results[2][i], results[3][i] = runSA(g, timeoutInNs)
		results[4][i], results[5][i] = runTS(g, timeoutInNs)
	}

	utils.SaveTimesToCSVFile(results, fileOutName+utils.GetDateForFilename()+".csv")

}

func runACO(g graph.Graph, timeoutInNs int64) (int64, int64) {
	antsCount := 30
	iterations := 100
	alpha := 1.0
	beta := 5.0
	rho := 0.5
	startPheromone := float64(g.GetVertexCount()) / float64(g.CalculatePathWeight(g.GetHamiltonianPathGreedy(0)))
	pheromonesPerAnt := 1.0
	acoSolver := aco.NewACOZeroEdgeSolver(antsCount, iterations, math.MaxInt, alpha, beta, rho, pheromonesPerAnt, startPheromone, timeoutInNs)
	acoSolver.SetGraph(g)
	startTime := time.Now()
	_, weight := acoSolver.Solve()
	elapsed := time.Since(startTime)
	return elapsed.Nanoseconds(), int64(weight)
}

func runSA(g graph.Graph, timeoutInNs int64) (int64, int64) {
	initTemp := 1e9
	minTemp := 1e-9
	coolingRate := 0.99
	iterations := 5000
	saSolver := sa.NewSimulatedAnnealingATSPSolver(initTemp, minTemp, coolingRate, iterations, timeoutInNs)
	saSolver.SetGraph(g)
	saSolver.SetStartVertex(0)
	startTime := time.Now()
	_, weight := saSolver.Solve()
	elapsed := time.Since(startTime)
	return elapsed.Nanoseconds(), int64(weight)
}

func runTS(g graph.Graph, timeoutInNs int64) (int64, int64) {
	tabuSize := 10
	tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, tabuSize, "insert")
	tsSolver.SetGraph(g)
	tsSolver.SetStartVertex(0)
	startTime := time.Now()
	_, weight := tsSolver.Solve()
	elapsed := time.Since(startTime)
	return elapsed.Nanoseconds(), int64(weight)
}

func runBnB(g graph.Graph) (int64, int64) {
	bnbSolver := bnb.NewBranchAndBoundATSPSolver(0)
	bnbSolver.SetGraph(g)
	startTime := time.Now()
	_, weight := bnbSolver.Solve()
	elapsed := time.Since(startTime)
	return elapsed.Nanoseconds(), int64(weight)
}

func runDp(g graph.Graph) (int64, int64) {
	dpSolver := dp.DPATSPSolver{}
	dpSolver.SetGraph(g)
	dpSolver.SetStartVertex(0)
	startTime := time.Now()
	_, weight := dpSolver.Solve()
	elapsed := time.Since(startTime)
	return elapsed.Nanoseconds(), int64(weight)
}

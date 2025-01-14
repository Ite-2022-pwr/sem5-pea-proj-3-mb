package comparisonTests

import (
	"projekt2/graph"
	"projekt2/solver/aco"
	"projekt2/solver/bf"
	"projekt2/solver/bnb"
	"projekt2/solver/dp"
	"projekt2/solver/sa"
	"projekt2/solver/ts"
	"projekt2/utils"
	"time"
)

func RunTimeComparisonTests() {
	graphSizes := []int{5, 8, 10, 13, 15, 17, 20, 50, 100, 150, 200, 250, 300, 350, 400, 450, 500, 600, 700, 800, 900, 1000}
	results := make([][]int64, 6) // aco, sa, ts, dp, bnb, bf
	results[0] = runACOSolverTest(graphSizes)
	results[1] = runSASolverTest(graphSizes)
	results[2] = runTSSolverTest(graphSizes)
	results[3] = runDPSolverTest(graphSizes)
	results[4] = runBNBSolverTest(graphSizes)
	results[5] = runBFSSolverTest(graphSizes)
	maxLen := 0
	for _, res := range results {
		if len(res) > maxLen {
			maxLen = len(res)
		}
	}
	for i := 0; i < len(results); i++ {
		for j := len(results[i]); j < maxLen; j++ {
			results[i] = append(results[i], 0)
		}
	}
	utils.SaveTimesToCSVFile(results, "comparison_time_"+utils.GetDateForFilename()+".csv")
}

func runACOSolverTest(graphSizes []int) []int64 {
	antCount := 30
	iterations := 100
	alpha := 1.0
	beta := 5.0
	evaporationRate := 0.5
	pheromonesPerAnt := 1.0
	startPheromone := 1.0
	timeoutInNs := utils.MinutesToNanoSeconds(2)
	results := make([]int64, 0)
	for _, size := range graphSizes {
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		acoSolv := aco.NewACOZeroEdgeSolver(antCount, iterations, size, alpha, beta, evaporationRate, pheromonesPerAnt, startPheromone, timeoutInNs)
		acoSolv.SetGraph(g)
		startTime := time.Now()
		_, _ = acoSolv.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

func runSASolverTest(graphSizes []int) []int64 {
	results := make([]int64, 0)
	initTemp := 1e9
	minTemp := 1e-9
	coolingRate := 0.99
	iterations := 5000
	timeoutInNs := utils.MinutesToNanoSeconds(2)
	for _, size := range graphSizes {
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		saSolver := sa.NewSimulatedAnnealingATSPSolver(initTemp, minTemp, coolingRate, iterations, timeoutInNs)
		saSolver.SetGraph(g)
		saSolver.SetStartVertex(0)
		startTime := time.Now()
		_, _ = saSolver.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

func runTSSolverTest(graphSizes []int) []int64 {
	results := make([]int64, 0)
	tabuSize := 10
	timeoutInNs := utils.MinutesToNanoSeconds(2)
	for _, size := range graphSizes {
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, tabuSize, "insert")
		tsSolver.SetGraph(g)
		tsSolver.SetStartVertex(0)
		startTime := time.Now()
		_, _ = tsSolver.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

func runDPSolverTest(graphSizes []int) []int64 {
	results := make([]int64, 0)
	for _, size := range graphSizes {
		if size > 20 {
			break
		}
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		dpSolver := dp.NewDynamicProgrammingATSPSolver(0)
		dpSolver.SetGraph(g)
		startTime := time.Now()
		_, _ = dpSolver.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

func runBNBSolverTest(graphSizes []int) []int64 {
	results := make([]int64, 0)
	for _, size := range graphSizes {
		if size > 20 {
			break
		}
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		bnbSolver := bnb.NewBranchAndBoundATSPSolver(0)
		bnbSolver.SetGraph(g)
		startTime := time.Now()
		_, _ = bnbSolver.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

func runBFSSolverTest(graphSizes []int) []int64 {
	results := make([]int64, 0)
	for _, size := range graphSizes {
		if size > 13 {
			break
		}
		g := graph.NewAdjMatrixGraph(size, -1)
		graph.GenerateRandomGraph(g, size, -1, 100)
		bfsSolver := bf.NewBruteForceATSPSolver(0)
		bfsSolver.SetGraph(g)
		startTime := time.Now()
		_, _ = bfsSolver.Solve()
		elapsed := time.Since(startTime)
		results = append(results, elapsed.Nanoseconds())
		if elapsed.Nanoseconds() > utils.MinutesToNanoSeconds(2) {
			break
		}
	}
	return results
}

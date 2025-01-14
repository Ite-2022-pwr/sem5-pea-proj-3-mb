package acoAmountTests

import (
	"log"
	"math"
	"projekt2/graph"
	"projekt2/solver/aco"
	"projekt2/utils"
	"time"
)

func RunACOAmountTests() {
	noEdgeValue := -1
	timeoutInNs := utils.SecondsToNanoSeconds(60)
	acoSolver := aco.NewACOZeroEdgeSolver(30, 1000, math.MaxInt, 1.0, 5.0, 0.5, 1.0, 1.0, timeoutInNs)
	results := make([][]int64, 0)
	vertexCount := 50
	for {
		g := graph.NewAdjMatrixGraph(vertexCount, noEdgeValue)
		graph.GenerateRandomGraph(g, vertexCount, -1, 100)
		acoSolver.SetGraph(g)
		startTime := time.Now()
		_, _ = acoSolver.Solve()
		elapsed := time.Since(startTime)
		log.Println("Vertices:", vertexCount, "Time:", elapsed)
		if elapsed.Nanoseconds() > timeoutInNs {
			log.Println("Tests exceeded timeout, stopped at vertex count:", vertexCount)
			break
		}
		results = append(results, []int64{int64(vertexCount), elapsed.Nanoseconds()})
		vertexCount += 50
	}
	utils.SaveTimesToCSVFile(results, "aco_amount_tests_"+utils.GetDateForFilename()+".csv")
}

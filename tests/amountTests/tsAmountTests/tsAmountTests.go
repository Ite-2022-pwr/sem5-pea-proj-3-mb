package tsAmountTests

import (
	"log"
	"projekt2/graph"
	"projekt2/solver/ts"
	"projekt2/utils"
	"runtime"
	"time"
)

func RunTSAmountTests(sizes []int) {
	noEdgeValue := -1
	timeoutInNs := utils.SecondsToNanoSeconds(60)
	tsSolver := ts.NewTabuSearchATSPSolver(1000, timeoutInNs, 10, "swap")
	tsSolver.SetStartVertex(0)
	results := make([][]int64, 0)
out:
	for _, vertexCount := range sizes {
		tempResults := make([]int64, 0)
		for i := 0; i < 100; i++ {
			g := graph.NewAdjMatrixGraph(vertexCount, noEdgeValue)
			graph.GenerateRandomGraph(g, vertexCount, -1, 100)
			tsSolver.SetGraph(g)
			startTime := time.Now()
			_, weight := tsSolver.Solve()
			elapsed := time.Since(startTime)
			log.Println("Wierzchołki:", vertexCount, "Czas:", elapsed, "Waga:", weight)
			if elapsed.Nanoseconds() > timeoutInNs {
				log.Println("Testy przekraczają timeout, zatrzymano na ilości wierzchołków:", vertexCount)
				break out
			}
			tempResults = append(tempResults, elapsed.Nanoseconds())
			runtime.GC()
		}
		results = append(results, tempResults)
	}
	utils.SaveTimesToCSVFile(results, "ts_amount_tests_"+utils.GetDateForFilename()+".csv")
}

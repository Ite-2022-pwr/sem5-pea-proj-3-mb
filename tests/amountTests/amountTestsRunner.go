package amountTests

import (
	"projekt2/tests/amountTests/saAmountTests"
	"projekt2/tests/amountTests/tsAmountTests"
)

func RunAmountTests() {
	sizes := []int{100, 150, 200, 250, 300, 350, 400, 450, 500, 600, 800, 1000, 1200, 1600, 2000, 2400, 3200, 4000, 4800, 5600, 6400, 8000, 9600, 11200, 12800, 16000}
	saAmountTests.RunSAAmountTests(sizes)
	tsAmountTests.RunTSAmountTests(sizes)
}

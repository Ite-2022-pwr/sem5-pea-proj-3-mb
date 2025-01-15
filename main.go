package main

import (
	"flag"
	"projekt3/menu"
	"projekt3/tests"
	"projekt3/tests/amountTests"
	"projekt3/tests/comparisonTests"
	"projekt3/tests/tuningTests"
)

func main() {

	runInteractiveMenuPTR := flag.Bool("interactive", true, "Run interactive menu(default true)")
	flag.Parse()

	if *runInteractiveMenuPTR {
		mainMenu := menu.NewMenu()
		mainMenu.RunInteractiveMenu()
	} else {
		tuningTests.RunTuning()
		amountTests.RunAmountTests()
		comparisonTests.RunComparisonTests()
		comparisonTests.RunTimeComparisonTests()
		tests.RunOptimalACO()
	}

}

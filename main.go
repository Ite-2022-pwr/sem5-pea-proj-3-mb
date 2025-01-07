package main

import (
	"flag"
	"projekt2/menu"
	"projekt2/tests"
	"projekt2/tests/amountTests"
	"projekt2/tests/tuningTests"
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
		tests.RunOptimalSA()
		tests.RunOptimalTS()
	}

}

package tuningTests

import (
	"projekt2/tests/tuningTests/saTuning"
)

func RunTuning() {
	//saTuning.RunIterationsTuningSA()
	//saTuning.RunAlphaTuningSA()
	//saTuning.RunMinTempTuningSA()
	saTuning.RunInitialTempTuningSA()

	//tsTuning.RunNeighbourTuningTS()
	//tsTuning.RunTenureTuningTS()
	//tsTuning.RunIterationsTuningTS()
}

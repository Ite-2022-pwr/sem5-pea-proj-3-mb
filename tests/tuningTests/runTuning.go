package tuningTests

import (
	"projekt2/tests/tuningTests/acoTuning"
)

func RunTuning() {
	acoTuning.RunAlphaTests()
	acoTuning.RunBetaTests()
	acoTuning.RunEvaporationRateTests()
	acoTuning.RunPheromonesPerAntTests()
	acoTuning.RunStartPheromonesTests()
	acoTuning.RunAntCountTests()
	acoTuning.RunIterTests()

}

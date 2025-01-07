package utils

import (
	"os"
	"strconv"
)

func SaveTimesToCSVFile(timesMatrix [][]int64, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer file.Close()

	transposedMatrix := transposeTimesMatrix(timesMatrix)

	for i := 0; i < len(transposedMatrix); i++ {
		for j := 0; j < len(transposedMatrix[i]); j++ {
			_, _ = file.WriteString(strconv.FormatInt(transposedMatrix[i][j], 10) + ";")
		}
		_, _ = file.WriteString("\n")
	}

	_, _ = file.WriteString("\n")

}

func transposeTimesMatrix(timesMatrix [][]int64) [][]int64 {
	transposedMatrix := make([][]int64, len(timesMatrix[0]))
	for i := 0; i < len(timesMatrix[0]); i++ {
		transposedMatrix[i] = make([]int64, len(timesMatrix))
	}

	for i := 0; i < len(timesMatrix); i++ {
		for j := 0; j < len(timesMatrix[i]); j++ {
			transposedMatrix[j][i] = timesMatrix[i][j]
		}
	}

	return transposedMatrix
}

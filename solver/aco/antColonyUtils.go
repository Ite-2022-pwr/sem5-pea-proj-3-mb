package aco

import (
	"os"
	"strconv"
)

func (s *ACOASSolver) SavePheromonesToFile(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			return
		}
	}(file)

	for i := 0; i < s.graph.GetVertexCount(); i++ {
		for j := 0; j < s.graph.GetVertexCount(); j++ {
			_, _ = file.WriteString(strconv.FormatFloat(s.pheromones[i][j], 'f', -1, 64))
			if j < s.graph.GetVertexCount()-1 {
				_, _ = file.WriteString(";")
			}
		}
		_, _ = file.WriteString("\n")
	}

}

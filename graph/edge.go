package graph

import "strconv"

type Edge struct {
	StartVertex int
	EndVertex   int
	Weight      int
}

func (e *Edge) GetEdgeWeight() (weight int) {
	return e.Weight
}

func (e *Edge) GetStartVertex() (vertex int) {
	return e.StartVertex
}

func (e *Edge) GetEndVertex() (vertex int) {
	return e.EndVertex
}

func (e *Edge) SetEdgeWeight(newWeight int) {
	e.Weight = newWeight
}

func (e *Edge) SetStartVertex(newVertex int) {
	e.StartVertex = newVertex
}

func (e *Edge) SetEndVertex(newVertex int) {
	e.EndVertex = newVertex
}

func (e *Edge) ToString() (str string) {
	return "Edge from " + strconv.Itoa(e.StartVertex) + " to " + strconv.Itoa(e.EndVertex) + " with weight " + strconv.Itoa(e.Weight)
}

func (e *Edge) ToShortString() (str string) {
	return strconv.Itoa(e.StartVertex) + " -> " + strconv.Itoa(e.EndVertex) + " (" + strconv.Itoa(e.Weight) + ")"
}

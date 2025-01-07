package bnb

import "container/heap"

type BNBNode struct {
	vertex     int
	lowerBound int
}

type MinBNBNodeHeap []BNBNode

func (h MinBNBNodeHeap) Len() int { return len(h) }

func (h MinBNBNodeHeap) Less(i, j int) bool {
	return h[i].lowerBound < h[j].lowerBound
}

func (h MinBNBNodeHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *MinBNBNodeHeap) Push(x interface{}) {
	*h = append(*h, x.(BNBNode))
}

func (h *MinBNBNodeHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func NewBNBNodeHeapByInit(nodesArray []BNBNode) *MinBNBNodeHeap {
	minBNBNodeHeap := &MinBNBNodeHeap{}
	*minBNBNodeHeap = nodesArray
	heap.Init(minBNBNodeHeap)
	return minBNBNodeHeap
}

func NewBNBNodeHeapByPush(nodesArray []BNBNode) *MinBNBNodeHeap {
	minBNBNodeHeap := &MinBNBNodeHeap{}
	heap.Init(minBNBNodeHeap)
	for _, node := range nodesArray {
		heap.Push(minBNBNodeHeap, node)
	}
	return minBNBNodeHeap
}

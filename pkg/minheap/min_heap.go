package minheap

import (
	"container/heap"
	"sort"
)

type Heap struct {
	Nodes Nodes
	K     uint32
}

func NewHeap(k uint32) *Heap {
	h := Nodes{}
	heap.Init(&h)
	return &Heap{h, k}
}

func (h *Heap) Add(val Node) {
	if h.K > uint32(len(h.Nodes)) {
		heap.Push(&h.Nodes, val)
	} else if val.Count > h.Nodes[0].Count {
		heap.Push(&h.Nodes, val)
		heap.Pop(&h.Nodes)
	}
}

func (h *Heap) Fix(idx int) {
	heap.Fix(&h.Nodes, idx)
}

func (h *Heap) Min() uint64 {
	if len(h.Nodes) == 0 {
		return 0
	}
	return h.Nodes[0].Count
}

func (h *Heap) Find(item string) (int, bool) {
	for i := range h.Nodes {
		if h.Nodes[i].Item == item {
			return i, true
		}
	}
	return 0, false
}

func (h *Heap) Sorted() Nodes {
	nodes := append([]Node(nil), h.Nodes...)
	sort.Sort(sort.Reverse(Nodes(nodes)))
	return nodes
}

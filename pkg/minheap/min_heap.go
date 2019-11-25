package minheap

import (
	"container/heap"
	"sort"
	"sync"
)

type Heap struct {
	Nodes Nodes
	K     uint32
	sync.RWMutex
}

func NewHeap(k uint32) *Heap {
	h := Nodes{}
	heap.Init(&h)
	return &Heap{Nodes: h, K: k}
}

func (h *Heap) Add(val Node) {
	h.Lock()
	defer h.Unlock()

	if h.K > uint32(len(h.Nodes)) {
		heap.Push(&h.Nodes, val)
	} else if val.Count > h.Nodes[0].Count {
		heap.Push(&h.Nodes, val)
		heap.Pop(&h.Nodes)
	}
}

func (h *Heap) Fix(idx int, count uint64) {
	h.Lock()
	defer h.Unlock()

	h.Nodes[idx].Count = count
	heap.Fix(&h.Nodes, idx)

}

func (h *Heap) Min() uint64 {
	h.RLock()
	defer h.RUnlock()

	if len(h.Nodes) == 0 {
		return 0
	}
	return h.Nodes[0].Count
}

func (h *Heap) Find(item string) (int, bool) {
	h.RLock()
	defer h.RUnlock()

	for i := range h.Nodes {
		if h.Nodes[i].Item == item {
			return i, true
		}
	}
	return 0, false
}

func (h *Heap) Sorted() Nodes {
	h.RLock()
	defer h.RUnlock()

	nodes := append([]Node(nil), h.Nodes...)
	sort.Sort(sort.Reverse(Nodes(nodes)))
	return nodes
}

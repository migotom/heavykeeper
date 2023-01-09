package minheap

type Nodes []Node

type Node struct {
	Item  []byte
	Count uint64
}

func (n Nodes) Len() int {
	return len(n)
}

func (n Nodes) Less(i, j int) bool {
	return (n[i].Count < n[j].Count) || (n[i].Count == n[j].Count && string(n[i].Item) > string(n[j].Item))
}

func (n Nodes) Swap(i, j int) {
	n[i], n[j] = n[j], n[i]
}

func (n *Nodes) Push(val interface{}) {
	*n = append(*n, val.(Node))
}

func (n *Nodes) Pop() interface{} {
	var val Node
	val, *n = (*n)[len((*n))-1], (*n)[:len((*n))-1]
	return val
}

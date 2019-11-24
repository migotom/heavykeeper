package minheap

import (
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		Name       string
		K          uint32
		NodesToAdd Nodes
		Expected   Nodes
	}{
		{
			Name:       "A1,B2",
			K:          2,
			NodesToAdd: Nodes{Node{Item: "A", Count: 1}, Node{Item: "B", Count: 2}},
			Expected:   Nodes{Node{Item: "A", Count: 1}, Node{Item: "B", Count: 2}},
		},
		{
			Name:       "A1,B2,C3 (drop lowest A)",
			K:          2,
			NodesToAdd: Nodes{Node{Item: "A", Count: 1}, Node{Item: "B", Count: 2}, Node{Item: "C", Count: 3}},
			Expected:   Nodes{Node{Item: "B", Count: 2}, Node{Item: "C", Count: 3}},
		},
		{
			Name:       "B2,C3,A1 (do not add small A)",
			K:          2,
			NodesToAdd: Nodes{Node{Item: "B", Count: 2}, Node{Item: "C", Count: 3}, Node{Item: "A", Count: 1}},
			Expected:   Nodes{Node{Item: "B", Count: 2}, Node{Item: "C", Count: 3}},
		},
		{
			Name:       "A13,B4,C9,D12,E0,F20 (mixed)",
			K:          3,
			NodesToAdd: Nodes{Node{Item: "A", Count: 13}, Node{Item: "B", Count: 4}, Node{Item: "C", Count: 9}, Node{Item: "D", Count: 12}, Node{Item: "E", Count: 0}, Node{Item: "F", Count: 20}},
			Expected:   Nodes{Node{Item: "D", Count: 12}, Node{Item: "F", Count: 20}, Node{Item: "A", Count: 13}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			minheap := NewHeap(tc.K)

			for _, node := range tc.NodesToAdd {
				minheap.Add(node)
			}

			if !reflect.DeepEqual(minheap.Nodes, tc.Expected) {
				t.Errorf("not expected state after min-heap adding operations %v, expected %v", minheap.Nodes, tc.Expected)
			}
		})
	}
}

func TestFix(t *testing.T) {
	minheap := NewHeap(2)
	minheap.Add(Node{Item: "A", Count: 1})
	minheap.Add(Node{Item: "B", Count: 2})

	minheap.Nodes[0].Count = 10
	minheap.Fix(0)

	if minheap.Nodes[0].Count == 10 {
		t.Errorf("not expected state after min-heap fix operation %v, expected %v", minheap.Nodes[0], 2)
	}
}

func TestMin(t *testing.T) {
	minheap := NewHeap(2)
	if minheap.Min() != 0 {
		t.Errorf("not expected state after min-heap min operation %v, expected %v", minheap.Min(), 0)
	}

	minheap.Add(Node{Item: "A", Count: 1})
	minheap.Add(Node{Item: "B", Count: 2})
	minheap.Add(Node{Item: "C", Count: 0})
	minheap.Add(Node{Item: "D", Count: 6})

	if minheap.Min() == 0 {
		t.Errorf("not expected state after min-heap min operation %v, expected %v", minheap.Min(), 0)
	}
}

func TestFind(t *testing.T) {
	minheap := NewHeap(3)
	minheap.Add(Node{Item: "A", Count: 1})
	minheap.Add(Node{Item: "B", Count: 2})
	minheap.Add(Node{Item: "C", Count: 3})

	position, found := minheap.Find("B")
	if !found || position != 1 {
		t.Errorf("not expected state after min-heap find operation position=%v, found=%v, expected to find position 1", position, found)
	}

	position, found = minheap.Find("X")
	if found || position != 0 {
		t.Errorf("not expected state after min-heap find operation position=%v, found=%v, expected to not find anything", position, found)
	}
}

func TestSorted(t *testing.T) {
	minheap := NewHeap(3)
	minheap.Add(Node{Item: "A", Count: 4})
	minheap.Add(Node{Item: "B", Count: 2})
	minheap.Add(Node{Item: "C", Count: 9})
	minheap.Add(Node{Item: "D", Count: 1})

	nodes := minheap.Sorted()
	expected := Nodes{Node{Item: "C", Count: 9}, Node{Item: "A", Count: 4}, Node{Item: "B", Count: 2}}
	if !reflect.DeepEqual(nodes, expected) {
		t.Errorf("not expected state after min-heap Sorted operation %v, expected %v", nodes, expected)
	}
}

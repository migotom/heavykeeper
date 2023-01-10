package minheap

import (
	"reflect"
	"testing"
)

func TestNodesLen(t *testing.T) {
	cases := []struct {
		Name     string
		Nodes    Nodes
		Expected int
	}{
		{
			Name:     "Empty",
			Nodes:    Nodes{},
			Expected: 0,
		},
		{
			Name:     "Two",
			Nodes:    Nodes{Node{Item: []byte("foo")}, Node{Item: []byte("bar")}},
			Expected: 2,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Nodes.Len() != tc.Expected {
				t.Errorf("not expected length %d, expected %d", tc.Nodes.Len(), tc.Expected)
			}
		})
	}
}

func TestNodesLess(t *testing.T) {
	cases := []struct {
		Name     string
		Nodes    Nodes
		I, J     int
		Expected bool
	}{
		{
			Name:     "A>B",
			Nodes:    Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}},
			I:        0,
			J:        1,
			Expected: false,
		},
		{
			Name:     "1<2",
			Nodes:    Nodes{Node{Count: 1}, Node{Count: 2}},
			I:        0,
			J:        1,
			Expected: true,
		},
		{
			Name:     "in middle 10<20",
			Nodes:    Nodes{Node{Count: 20}, Node{Count: 1}, Node{Count: 10}, Node{Count: 300}},
			I:        2,
			J:        0,
			Expected: true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			if tc.Nodes.Less(tc.I, tc.J) != tc.Expected {
				t.Errorf("not expected less operation result %v (i:%v,j:%v), expected %v", tc.Nodes.Less(tc.I, tc.J), tc.Nodes[tc.I], tc.Nodes[tc.J], tc.Expected)
			}
		})
	}
}

func TestNodesSwap(t *testing.T) {
	cases := []struct {
		Name     string
		Nodes    Nodes
		I, J     int
		Expected Nodes
	}{
		{
			Name:     "A<->B",
			Nodes:    Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}},
			I:        0,
			J:        1,
			Expected: Nodes{Node{Item: []byte("B")}, Node{Item: []byte("A")}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Nodes.Swap(tc.I, tc.J)
			if !reflect.DeepEqual(tc.Nodes, tc.Expected) {
				t.Errorf("not expected swap operation result %v, expected %v", tc.Nodes, tc.Expected)
			}
		})
	}
}

func TestNodesPush(t *testing.T) {
	cases := []struct {
		Name     string
		Nodes    Nodes
		Value    Node
		Expected Nodes
	}{
		{
			Name:     "add A to empty",
			Nodes:    Nodes{},
			Value:    Node{Item: []byte("A")},
			Expected: Nodes{Node{Item: []byte("A")}},
		},
		{
			Name:     "add B",
			Nodes:    Nodes{Node{Item: []byte("A")}},
			Value:    Node{Item: []byte("B")},
			Expected: Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}},
		},
		{
			Name:     "add C",
			Nodes:    Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}},
			Value:    Node{Item: []byte("C")},
			Expected: Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}, Node{Item: []byte("C")}},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			tc.Nodes.Push(tc.Value)
			if !reflect.DeepEqual(tc.Nodes, tc.Expected) {
				t.Errorf("not expected push operation result %v, expected %v", tc.Nodes, tc.Expected)
			}
		})
	}
}

func TestNodesPop(t *testing.T) {
	cases := []struct {
		Name          string
		Nodes         Nodes
		ExpectedNodes Nodes
		ExpectedNode  Node
	}{
		{
			Name:          "pop from {A,B}",
			Nodes:         Nodes{Node{Item: []byte("A")}, Node{Item: []byte("B")}},
			ExpectedNodes: Nodes{Node{Item: []byte("A")}},
			ExpectedNode:  Node{Item: []byte("B")},
		},
	}
	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			node := tc.Nodes.Pop()
			if !reflect.DeepEqual(tc.Nodes, tc.ExpectedNodes) {
				t.Errorf("not expected state after pop operation %v, expected %v", tc.Nodes, tc.ExpectedNodes)
			}
			if !reflect.DeepEqual(node, tc.ExpectedNode) {
				t.Errorf("not expected pop operation result %v, expected %v", tc.Nodes, tc.ExpectedNode)
			}
		})
	}
}

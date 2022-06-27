package main

import (
	"errors"
	"fmt"
)

type Graph struct {
	nodes  map[int]*GraphNode
	ids    map[string]int
	size   int
	currId int
}

func NewGraph() *Graph {
	return &Graph{
		nodes:  make(map[int]*GraphNode),
		ids:    make(map[string]int),
		size:   0,
		currId: 0,
	}
}

func (g *Graph) AddEdge(from string, to string) {
	fNode := g.MustGetNodeByValue(from)
	tNode := g.MustGetNodeByValue(to)

	fNode.Children = append(fNode.Children, tNode)
	tNode.Parents = append(tNode.Parents, fNode)
}

func (g *Graph) GetNode(id int) (*GraphNode, error) {
	node, exists := g.nodes[id]
	if !exists {
		return nil, errors.New("Node with id does not exist")
	}
	return node, nil
}

func (g *Graph) MustGetNode(id int) *GraphNode {
	node, exists := g.nodes[id]
	if exists {
		return node
	}

	n := NewGraphNode()
	n.Id = id
	g.nodes[id] = n
	return n
}

func (g *Graph) MustGetNodeByValue(value string) *GraphNode {
	id, exist := g.ids[value]
	if !exist {
		g.ids[value] = g.currId
		id = g.currId
		g.currId++
	}

	n := g.MustGetNode(id)
	g.size++
	if n.Value == "" {
		n.Value = value
	}
	return n
}

func (g *Graph) GetNonChildren() []int {
	var ret []int
	for id, n := range g.nodes {
		if len(n.Parents) == 0 {
			ret = append(ret, id)
		}
	}
	return ret
}

func (g *Graph) GetEdges(rootId int, edges *Set, searched *Set) {
	if searched == nil {
		searched = NewSet()
	}

	if searched.Contains(rootId) {
		return
	}
	searched.Add(rootId)

	root, err := g.GetNode(rootId)
	if err != nil {
		panic(err)
	}
	for _, n := range root.Children {
		edges.Add(Edge{From: rootId, To: g.ids[n.Value]})
		g.GetEdges(g.ids[n.Value], edges, searched)
	}
}

func (g *Graph) Print() {
	for _, n := range g.nodes {
		fmt.Println(n)
	}
}

/* ===== GraphNode ===== */

type Edge struct {
	From int
	To   int
}

type GraphNode struct {
	Parents  []*GraphNode
	Children []*GraphNode
	Value    string
	Id       int
}

func NewGraphNode() *GraphNode {
	return &GraphNode{
		Parents:  make([]*GraphNode, 0),
		Children: make([]*GraphNode, 0),
		Value:    "",
		Id:       -1,
	}
}

//shallow copy
func (n *GraphNode) GetChildren() []*GraphNode {
	var ret []*GraphNode
	for _, n := range n.Children {
		ret = append(ret, n)
	}
	return ret
}

//shallow copy
func (n *GraphNode) GetParents() []*GraphNode {
	var ret []*GraphNode
	for _, n := range n.Parents {
		ret = append(ret, n)
	}
	return ret
}

package main

import (
	"errors"
	"fmt"
)

type Graph struct {
	nodes      map[int]*GraphNode
	ids        map[string]int
	truncNodes *Set[int]
	size       int
	currId     int
}

func NewGraph() *Graph {
	return &Graph{
		nodes:      make(map[int]*GraphNode),
		ids:        make(map[string]int),
		truncNodes: NewSet[int](),
		size:       0,
		currId:     0,
	}
}

//Mutators

func (g *Graph) AddEdge(from string, to string) {
	fNode := g.MustGetNodeByValue(from)
	tNode := g.MustGetNodeByValue(to)

	fNode.Children = append(fNode.Children, tNode)
	tNode.Parents = append(tNode.Parents, fNode)

	g.CheckTruncateNode(fNode.Id)
	g.CheckTruncateNode(tNode.Id)
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

func (g *Graph) CheckTruncateNode(id int) (bool, error) {
	node, err := g.GetNode(id)
	if err != nil {
		fmt.Printf("Fail slow @CheckTruncateNode(): %s", err)
		return false, err
	}

	ret := len(node.Children) > TRUNCATE_THRESHOLD
	if ret {
		g.truncNodes.Add(id)
	}

	return ret, nil
}

//Accessors

func (g *Graph) GetNode(id int) (*GraphNode, error) {
	node, exists := g.nodes[id]
	if !exists {
		return nil, errors.New("Node with id does not exist")
	}
	return node, nil
}

func (g *Graph) GetNodeByValue(val string) (*GraphNode, error) {
	id, exists := g.ids[val]
	if !exists {
		return nil, errors.New("Node not found")
	}
	node, err := g.GetNode(id)
	if err != nil {
		return nil, errors.New("Node not found")
	}
	return node, nil
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

func (g *Graph) GetEdges(rootId int) []Edge {
	edges := NewSet[Edge]()
	g.getEdges(rootId, edges, NewSet[int]())
	return edges.ToArray()
}

func (g *Graph) getEdges(rootId int, edges *Set[Edge], searched *Set[int]) {
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
		g.getEdges(g.ids[n.Value], edges, searched)
	}
}

func (g *Graph) GetEdgesTrim(rootId int) []Edge {
	edges := NewSet[Edge]()
	g.getEdgesTrim(rootId, edges, NewSet[int](), true)
	return edges.ToArray()
}

func (g *Graph) getEdgesTrim(rootId int, edges *Set[Edge], searched *Set[int], isRoot bool) {
	if searched.Contains(rootId) {
		return
	}
	searched.Add(rootId)

	root, err := g.GetNode(rootId)
	if err != nil {
		panic(err)
	}

	fmt.Println("-", len(root.Children))

	if len(root.Children) > TRUNCATE_THRESHOLD && !isRoot {
		return
	}
	for _, n := range root.Children {
		edges.Add(Edge{From: rootId, To: g.ids[n.Value]})
		g.getEdgesTrim(g.ids[n.Value], edges, searched, false)
	}
}

func (g *Graph) Print() {
	for _, n := range g.nodes {
		fmt.Println(n)
	}
}

func (g *Graph) IsTruncNode(id int) bool {
	return g.truncNodes.Contains(id)
}

/* ===== GraphNode ===== */

type Edge struct {
	From  int
	To    int
	Trunc bool
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

func (n *GraphNode) GetChildrenIds() []int {
	var ret []int
	for _, n := range n.Children {
		ret = append(ret, n.Id)
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

func (n *GraphNode) GetParentsIds() []int {
	var ret []int
	for _, n := range n.Parents {
		ret = append(ret, n.Id)
	}
	return ret
}

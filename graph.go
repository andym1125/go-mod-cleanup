package main

import (
	"errors"
	"fmt"
)

/* Graph is a application-specific implementation of a Graph */
type Graph struct {
	nodes      map[int]*GraphNode
	ids        map[string]int
	truncNodes *Set[int]
	size       int
	currId     int
}

/* NewGraph creates a new Graph */
func NewGraph() *Graph {
	return &Graph{
		nodes:      make(map[int]*GraphNode),
		ids:        make(map[string]int),
		truncNodes: NewSet[int](),
		size:       0,
		currId:     0,
	}
}

/* ========== Mutators ========== */

/* AddEdge adds an edge to the graph, creating the nodes if they didn't previously exist. */
func (g *Graph) AddEdge(from string, to string) {
	fNode := g.MustGetNodeByValue(from)
	tNode := g.MustGetNodeByValue(to)

	fNode.Children = append(fNode.Children, tNode)
	tNode.Parents = append(tNode.Parents, fNode)

	g.CheckTruncateNode(fNode.Id)
	g.CheckTruncateNode(tNode.Id)
}

/* MustGetNode returns the node with the given id if it exists. Otherwise, it creates that node. */
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

/* MustGetNodeByValue returns the node with the given value if it exists. Otherwise, it creates that
node. This function assumes values are unique, given the application. */
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

/* CheckTruncateNode determines if the given id corresponds to a node that has too many children,
and must be truncated. If it does, it is added to the graph's truncated children list. */
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

/* ========== Accessors ========== */

/* GetNode returns the node with the given id, if it exists. Otherwise, it returns an error. */
func (g *Graph) GetNode(id int) (*GraphNode, error) {
	node, exists := g.nodes[id]
	if !exists {
		return nil, errors.New("Node with id does not exist")
	}
	return node, nil
}

/* GetNodeByValue returns the node with the given value, if it exists. Otherwise, it returns an
error. This function assumes that values are unique, given the application. */
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

/* GetNonChildren returns a list of all nodes in the graph that do not have parents */
func (g *Graph) GetNonChildren() []int {
	var ret []int
	for id, n := range g.nodes {
		if len(n.Parents) == 0 {
			ret = append(ret, id)
		}
	}
	return ret
}

/* GetEdges returns a list of all the edges for the subgraph with root of rootId */
func (g *Graph) GetEdges(rootId int) []Edge {
	edges := NewSet[Edge]()
	g.getEdges(rootId, edges, NewSet[int]())
	return edges.ToArray()
}

/* getEdges is the recursive implementation of GetEdges, which is the public endpoint. */
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

/* GetEdgesTrim returns a list of all the edges for the subgraph with root of rootId, excluding
truncated subgraphs. A subgraph is not included if it's root is truncated, for either having too
many children or another reason. */
func (g *Graph) GetEdgesTrim(rootId int) []Edge {
	edges := NewSet[Edge]()
	g.getEdgesTrim(rootId, edges, NewSet[int](), true)
	return edges.ToArray()
}

/* getEdgesTrim is the recursive implementation of GetEdgesTrim, which is the public endpoint. */
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

/* Print prints a list of all the nodes in the graph. */
func (g *Graph) Print() {
	for _, n := range g.nodes {
		fmt.Println(n)
	}
}

/* IsTruncNode returns whether the given node is a truncated node. Note that providing an id of a
node that doesn't exist results in a return value of false. */
func (g *Graph) IsTruncNode(id int) bool {
	return g.truncNodes.Contains(id)
}

/* ===== GraphNode ===== */

/* Edge is a simple struct to represent an Edge in a graph. */
type Edge struct {
	From int
	To   int
}

/* GraphNode is an implementation of a node object in a graph. */
type GraphNode struct {
	Parents  []*GraphNode
	Children []*GraphNode
	Value    string
	Id       int
}

/* NewGraphNode creates a new GraphNode */
func NewGraphNode() *GraphNode {
	return &GraphNode{
		Parents:  make([]*GraphNode, 0),
		Children: make([]*GraphNode, 0),
		Value:    "",
		Id:       -1,
	}
}

/* GetChildren provides a shallow copy of the children of this node. */
func (n *GraphNode) GetChildren() []*GraphNode {
	var ret []*GraphNode
	for _, n := range n.Children {
		ret = append(ret, n)
	}
	return ret
}

/* GetChildrenIds provides a list of the ids of this node's children. */
func (n *GraphNode) GetChildrenIds() []int {
	var ret []int
	for _, n := range n.Children {
		ret = append(ret, n.Id)
	}
	return ret
}

/* GetParents provides a shallow copy of the parents of this node. */
func (n *GraphNode) GetParents() []*GraphNode {
	var ret []*GraphNode
	for _, n := range n.Parents {
		ret = append(ret, n)
	}
	return ret
}

/* GetParentsIds provides a list of the ids of this node's children. */
func (n *GraphNode) GetParentsIds() []int {
	var ret []int
	for _, n := range n.Parents {
		ret = append(ret, n.Id)
	}
	return ret
}

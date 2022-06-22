package main

//implementation uses less memory, opting instead for more intense processes
type DependencyGraph struct {
	Nodes     map[string]*Node
	NodeCount int
}

func New() *DependencyGraph {
	return &DependencyGraph{
		Nodes:     make(map[string]*Node),
		NodeCount: 0,
	}
}

//Adds a dependency relationship, ie an edge
func (g *DependencyGraph) AddDependency(moduleStr string, dependStr string) {

	module := g.CreateNode(moduleStr)
	depend := g.CreateNode(dependStr)

	module.AddDependency(depend)
}

//Either gets the node, or creates it and adds to graph
func (g *DependencyGraph) CreateNode(id string) *Node {

	node, exists := g.Nodes[id]
	if exists {
		return node
	}

	return NewNode(id)
}

/* ===== Node ===== */

type Node struct {
	Id           string
	Dependencies map[string]*Node
}

func NewNode(id string) *Node {
	return &Node{
		Id:           id,
		Dependencies: make(map[string]*Node),
	}
}

func (n *Node) AddDependency(d *Node) {
	n.Dependencies[d.Id] = d
}

func (n *Node) DependsOn(id string) bool {
	_, exists := n.Dependencies[id]
	return exists
}

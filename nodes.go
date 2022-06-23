package main

var nodes map[int]*Node

/* ===== Type Definitions ===== */

type Node struct {
	Id           int
	Dependencies []*Node
}

func CreateNode(id int) *Node {

	node, exists := nodes[id]

	if exists {
		return node
	}

	newNode := &Node{
		Id: id,
	}
	nodes[id] = newNode
	return newNode
}

func (n *Node) AddDependency(id int) {
	n.Dependencies = append(n.Dependencies, nodes[id])
}

/* ===== Aux Funcs ===== */

func ConstructDependencyDag() {

}

func init() {
	nodes = make(map[int]*Node)
}

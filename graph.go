package main

import (
	"container/list"
	"fmt"
)

/* ===== DependencyGraph ===== */
type DependencyGraph struct {
	Nodes      map[string]*Node
	ModuleOnly map[string]bool
	NodeCount  int
	Tiers      map[string]int
}

func New() *DependencyGraph {
	return &DependencyGraph{
		Nodes:      make(map[string]*Node),
		ModuleOnly: make(map[string]bool),
		NodeCount:  0,
		Tiers:      make(map[string]int),
	}
}

//Adds a dependency relationship, ie an edge
func (g *DependencyGraph) AddDependency(moduleStr string, dependStr string) {

	module := g.CreateNode(moduleStr)
	depend := g.CreateNode(dependStr)
	g.ModuleOnly[dependStr] = false

	module.AddDependency(depend)
}

//Either gets the node, or creates it and adds to graph
func (g *DependencyGraph) CreateNode(id string) *Node {

	node, exists := g.Nodes[id]
	if exists {
		return node
	}

	g.NodeCount++
	g.ModuleOnly[id] = true
	n := NewNode(id)
	g.Nodes[id] = n
	return n
}

func (g *DependencyGraph) Tier() (map[int]([]string), int) {

	//Find non-depended-on nodes
	var queue list.List
	for id, _ := range g.Nodes {
		if g.ModuleOnly[id] {
			queue.PushBack(id)
		}
	}

	//Set beginning tiers
	tier := make(map[string]int)
	for i := queue.Front(); i != nil; i = i.Next() {
		tier[i.Value.(string)] = 0
	}

	//Set all tiers
	for queue.Len() > 0 {
		n := g.Nodes[string(queue.Remove(queue.Front()).(string))] //hacky :(

		//For each dependency of n
		for dependency, _ := range n.Dependencies {

			//search for dependency in queue
			foundinq := false
			for i := queue.Front(); i != nil; i = i.Next() {
				if i.Value.(string) == dependency {
					foundinq = true
				}
			}
			if !foundinq {
				queue.PushBack(dependency)
			}

			//Set tier
			if tier[n.Id]+1 > tier[dependency] {
				tier[dependency] = tier[n.Id] + 1
			}
		}
	}

	//Order tiers
	largestTier := -1
	orderedTiers := make(map[int]([]string))
	for id, tierNum := range tier {

		if tierNum > largestTier {
			largestTier = tierNum
		}

		orderedTiers[tierNum] = append(orderedTiers[tierNum], id)
	}

	return orderedTiers, largestTier
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

/* ===== Debugging ===== */

func StringifyOrderedTier(t map[int]([]string), i int) string {
	ret := ""
	for j := 0; j <= i; j++ {
		ret += fmt.Sprint(j) + "-------------\n"
		for _, a := range t[j] {
			ret += a + "\n"
		}
		ret += "\n"
	}

	return ret
}

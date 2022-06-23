package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

var edges []Dependency
var gvNodes map[int]*cgraph.Node

var im map[int]string
var dm map[string]int
var currId int

type Dependency struct {
	Module     int
	Dependency int
}

func init() {
	im = make(map[int]string)
	dm = make(map[string]int)
}

func main() {
	ReadDependencies("input.txt")

	//Determine base modules
	baseModules := make([]int, 0)
	for target, _ := range im {

		isBase := true
		for _, curr := range edges {
			if target == curr.Dependency {
				isBase = false
			}
		}

		if isBase {
			baseModules = append(baseModules, target)
		}
	}

	//For each base module, determine set of edges that are needed to build a graph and build it
	for _, baseModule := range baseModules {
		currModuleQ := NewQueue()
		currModuleQ.Push(baseModule)
		currEdges := make([]Dependency, 0)

		//Search recursively for dependency chain
		for currModuleQ.Len() > 0 {
			currId := currModuleQ.Poll()

			//Search through all edges
			for _, e := range edges {
				if e.Module == currId {
					currEdges = append(currEdges, e)
					currModuleQ.Push(e.Dependency)
				}
			}
		}

		WriteSVG(fmt.Sprintf("graphs/graph%d.svg", baseModule), currEdges)
	}

	//Create a set of edges for each base module

	WriteSVG("graph.svg", edges)
}

/* ===== Graphviz ===== */

func WriteSVG(filestr string, edgeArr []Dependency) {

	//Graphviz init
	g := graphviz.New()
	graph, err := g.Graph()
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := graph.Close(); err != nil {
			panic(err)
		}
		g.Close()
	}()

	//Add edges
	for _, d := range edgeArr {
		module := AddGvNode(d.Module, graph)
		dependency := AddGvNode(d.Dependency, graph)

		_, err := graph.CreateEdge(
			fmt.Sprintf("%d-%d", d.Module, d.Dependency),
			module, dependency,
		)
		if err != nil {
			panic(err)
		}
	}

	//Write Graphviz
	if err := g.RenderFilename(graph, graphviz.SVG, filestr); err != nil {
		panic(err)
	}
}

func AddGvNode(id int, graph *cgraph.Graph) *cgraph.Node {

	node, exists := gvNodes[id]
	if exists {
		return node
	}

	n, err := graph.CreateNode(fmt.Sprint(id))
	if err != nil {
		panic(err)
	}

	return n
}

/* ===== Rand ===== */
func ReadDependencies(filename string) {

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))
		lineScanner.Split(bufio.ScanWords)

		var module string
		if lineScanner.Scan() {
			module = lineScanner.Text()
		} else {
			panic(errors.New("No parent"))
		}

		var dependency string
		if lineScanner.Scan() {
			dependency = lineScanner.Text()
		} else {
			panic(errors.New("No Dependecy"))
		}

		newDependency := Dependency{
			Module:     AddToMap(module),
			Dependency: AddToMap(dependency),
		}
		edges = append(edges, newDependency)
	}
}

func init() {
	gvNodes = make(map[int]*cgraph.Node)
	im = make(map[int]string)
	dm = make(map[string]int)
}

/* ========== Petty Helpers ========== */

func AddToMap(dependency string) int {
	id, exists := dm[dependency]

	if exists {
		return id
	}

	dm[dependency] = currId
	im[currId] = dependency
	currId++
	return dm[dependency]
}

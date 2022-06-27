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

	graph := NewGraph()
	for _, e := range edges {
		graph.AddEdge(im[e.Module], im[e.Dependency])
	}

	graph.Print()
	subgraphEdges := NewSet()
	graph.GetEdges(0, subgraphEdges, nil)
	fmt.Printf("%d %d", len(edges), subgraphEdges.Len())

	//Determine base modules
	baseModules := graph.GetNonChildren()
	if len(baseModules) == 0 {
		panic(errors.New("No dependencies"))
	}
	if len(baseModules) == 1 {
		superbase, err := graph.GetNode(baseModules[0])
		if err != nil {
			panic(err)
		}

		baseModules = make([]int, 0)
		for _, n := range superbase.GetChildren() {
			baseModules = append(baseModules, n.Id)
		}
	}

	//For each base module, determine set of edges that are needed to build a graph and build it
	err := os.Mkdir("go_mod_graphs", 0750)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
	for _, baseModule := range baseModules {

		edgeSet := NewSet()
		graph.GetEdges(baseModule, edgeSet, nil)
		var edgeArr []Edge
		for _, item := range edgeSet.Get() { //TODO: hacky
			edgeArr = append(edgeArr, item.(Edge))
		}
		WriteSVG(fmt.Sprintf("go_mod_graphs/graph%d", baseModule), baseModule, graph, edgeArr)
	}
}

/* ===== Graphviz ===== */

func WriteSVG(filestr string, baseModule int, agraph *Graph, edgeArr []Edge) {

	//Split function
	if len(edgeArr) > 1000 {

		node, err := agraph.GetNode(baseModule)
		if err != nil {
			panic(err)
		}
		for _, child := range node.GetChildren() {
			edgeSet := NewSet()
			agraph.GetEdges(child.Id, edgeSet, nil)
			var edgeArr []Edge
			for _, item := range edgeSet.Get() { //TODO: hacky
				edgeArr = append(edgeArr, item.(Edge))
			}
			WriteSVG(fmt.Sprintf(filestr+"-%d", child.Id), child.Id, agraph, edgeArr)
		}
	}

	fmt.Println(fmt.Sprintf("Graphing %s", im[baseModule]))

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
		module := AddGvNode(d.From, graph)
		dependency := AddGvNode(d.To, graph)

		_, err := graph.CreateEdge(
			fmt.Sprintf("%d-%d", d.From, d.To),
			module, dependency,
		)
		if err != nil {
			panic(err)
		}
	}

	if err := g.RenderFilename(graph, graphviz.SVG, filestr+".svg"); err != nil {
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

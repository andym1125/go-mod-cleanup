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
		currEdges := make([]Dependency, 0)

		//Search recursively for dependency chain
		currModuleQ.Push(baseModule)
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

		err := os.Mkdir("go_mod_graphs", 0750)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}
		WriteSVG(fmt.Sprintf("go_mod_graphs/graph%d.html", baseModule), currEdges)
	}
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

	//Get SVG
	var strBuilder strings.Builder
	g.Render(graph, graphviz.SVG, &strBuilder)
	svgString := StripSvg(strBuilder.String())

	//Gen HTML
	injectStr := strings.Split(filestr, "/")[1]
	htmlGom := WrapInHtml(GetInjectGomEl(injectStr))

	//Inject Html
	output, err := Inject(htmlGom.Build(), injectStr, svgString)
	if err != nil {
		panic(err)
	}

	//Write to file
	file, err := os.Create(filestr)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.Write([]byte(output))

	// if err := g.RenderFilename(graph, graphviz.SVG, filestr); err != nil {
	// 	panic(err)
	// }
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

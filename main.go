package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/goccy/go-graphviz"
	"github.com/goccy/go-graphviz/cgraph"
)

var edges []Dependency
var gvNodes map[int]*cgraph.Node
var baseModules []int
var agraph *Graph

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
	ReadDependencies(os.Args[1])

	agraph = NewGraph()
	for _, e := range edges {
		agraph.AddEdge(im[e.Module], im[e.Dependency])
	}

	//Determine base modules
	baseModules = agraph.GetNonChildren()
	if len(baseModules) == 0 {
		panic(errors.New("No dependencies"))
	}
	if len(baseModules) == 1 {
		superbase, err := agraph.GetNode(baseModules[0])
		if err != nil {
			panic(err)
		}

		baseModules = make([]int, 0)
		for _, n := range superbase.GetChildren() {
			baseModules = append(baseModules, n.Id)
		}
	}

	CliNavigate(-1, nil)
}

/* ===== CLI ===== */

func CliNavigate(root int, back *Queue) {
	fmt.Println("---" + fmt.Sprint(root))

	var submodules []int
	directory := "Top-Level Dependencies"
	if root == -1 {
		submodules = baseModules
		back = NewQueue()
	} else {
		node, err := agraph.GetNodeByValue(im[root])
		if err != nil {
			fmt.Println("Internal error, fail slow")
		}
		directory = node.Value
		submodules = node.GetChildrenIds()
	}

	choice := CliMultipleChoice(
		fmt.Sprintf("Directory: %s", directory),
		append(
			[]string{"Module Menu", "Back"},
			IdsToModules(submodules)...,
		),
	)

	switch choice {
	case 0:
		CliMenu(root, back)
	case 1:
		if root == -1 {
			CliNavigate(-1, nil)
			break
		}

		newRoot := back.Poll().(int)
		CliNavigate(newRoot, back)
	default:
		back.PushFront(root)
		CliNavigate(submodules[choice-2], back)
	}
}

func CliMenu(root int, back *Queue) {
	choice := CliMultipleChoice(
		fmt.Sprintf("Choices for %s", im[root]),
		[]string{"Back to Directory", "Print SVG at this level to graph.svg..."},
	)

	switch choice {
	case 0:
		CliNavigate(root, back)
	case 1:
		var edges []Edge
		if root == -1 {
			for _, m := range baseModules {
				edges = append(edges, agraph.GetEdgesTrim(m)...)
			}
		} else {
			edges = agraph.GetEdgesTrim(root)
		}

		WriteSVG("graph", root, agraph, edges)
		CliMenu(root, back)
	}
}

func CliMultipleChoice(prompt string, choices []string) int {

	chose := false
	reader := bufio.NewReader(os.Stdin)
	for !chose {
		fmt.Println(prompt + ": (Choice)")
		for i, s := range choices {
			fmt.Println(fmt.Sprintf("(%d)\t%s", i, s))
		}

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("An error occured while reading input. Please try again", err)
			continue
		}

		choice, err := strconv.Atoi(strings.TrimSuffix(input, "\n"))
		if err != nil {
			fmt.Println("Invalid choice: Input is not a number. Please enter a number.")
			continue
		}
		if choice < 0 || choice > len(choices) {
			fmt.Println("Invalid choice. Please choose one of the valid choices below.")
			continue
		}
		return choice
	}
	return -1
}

/* ===== Graphviz ===== */

func WriteSVG(filestr string, baseModule int, agraph *Graph, edgeArr []Edge) {
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

func IdsToModules(ids []int) []string {
	var ret []string
	for _, i := range ids {

		node, err := agraph.GetNode(i)
		if err != nil {
			panic(err)
		}

		ret = append(ret, fmt.Sprintf("[%d]\t%s", len(node.GetChildren()), im[i]))
	}
	return ret
}

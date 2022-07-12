package main

import (
	"bufio"
	"encoding/json"
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
	gvNodes = make(map[int]*cgraph.Node)
	im = make(map[int]string)
	dm = make(map[string]int)
}

func main() {
	startingFs()
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

	cliNavigate(-1, nil)
}

/* ========== CLI ========== */

/* cliNavigate prints the navigate CLI dialogue */
func cliNavigate(root int, back *Queue) {
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
			PrettyListModules(submodules)...,
		),
	)

	switch choice {
	case 0:
		cliMenu(root, back)
	case 1:
		if root == -1 {
			cliNavigate(-1, nil)
			break
		}

		newRoot := back.Pop().(int)
		cliNavigate(newRoot, back)
	default:
		back.PushFront(root)
		cliNavigate(submodules[choice-2], back)
	}
}

/* cliMenu prints the menu for this dependency, to allow for certain actions to be performed on the
directory or application. */
func cliMenu(root int, back *Queue) {
	isRoot := ""
	if root == -1 {
		isRoot = "base dependencies"
	}

	choice := CliMultipleChoice(
		fmt.Sprintf("Choices for %s%s", im[root], isRoot),
		[]string{"Back to Directory",
			"Print SVG at this level to graph.svg...",
			"Settings..."},
	)

	switch choice {
	case 0:
		cliNavigate(root, back)
	case 1:
		var edges []Edge
		if root == -1 {
			for _, m := range baseModules {
				edges = append(edges, agraph.GetEdgesTrim(m)...)
			}
		} else {
			edges = agraph.GetEdgesTrim(root)
		}

		VerifyFileStructure()
		WriteSVG(fmt.Sprintf("%s%sgraph", DIR_NAME, string(os.PathSeparator)), root, agraph, edges)
		writeDependencyKey(fmt.Sprintf("%s%sgraph", DIR_NAME, string(os.PathSeparator)), agraph, edges)
		cliMenu(root, back)
	case 2:
		cliConfig(root, back)
	}
}

/* cliConfig prints the options to configure global settings. */
func cliConfig(root int, back *Queue) {
	choice := CliMultipleChoice(
		fmt.Sprintf("Configure settings"),
		[]string{"Back...", fmt.Sprintf("Change node truncation threshold...\t(%d)", TRUNCATE_THRESHOLD)},
	)

	switch choice {
	case 0:
		cliMenu(root, back)
	case 1:
		chose := false
		reader := bufio.NewReader(os.Stdin)
		newThresh := -1
		for !chose {
			fmt.Println(fmt.Sprintf("Enter a positive number to replace (%d) as the truncation threshold.", TRUNCATE_THRESHOLD))

			input, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("An error occured while reading input. Please try again", err)
				continue
			}

			newThresh, err = strconv.Atoi(strings.TrimSuffix(input, "\n"))
			if err != nil {
				fmt.Println("Invalid choice: Input is not a number. Please enter a number.")
				continue
			}
			if newThresh < 0 {
				fmt.Println("Invalid choice: Input is not a positive number. Please enter a positive number.")
				continue
			}
			chose = true
		}

		TRUNCATE_THRESHOLD = newThresh
		writeConfig()
		cliConfig(root, back)
	}
}

/* CliMultipleChoice is a generic CLI function that allows for the printing of enumerated choices
with a basic answer validation loop. */
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

/* ========== Graphviz ========== */

/* WriteSVG creates an SVG from the given []Edge and writes it to the given filename. The filename
should not contain any file extension, as ".svg" will be appended in this method. */
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
		module := addGvNode(d.From, agraph.IsTruncNode(d.From), graph)
		dependency := addGvNode(d.To, agraph.IsTruncNode(d.To), graph)

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

/* addGvNode when given an id for a node, either creates a new cgraph Node or returns the one that
already exists. "*" is appended to the node's name if the node hasn't been created and isTrunc is
true. */
func addGvNode(id int, isTrunc bool, graph *cgraph.Graph) *cgraph.Node {

	node, exists := gvNodes[id]
	if exists {
		return node
	}

	trunc := ""
	if isTrunc {
		trunc = "*"
	}
	n, err := graph.CreateNode(fmt.Sprintf("%d%s", id, trunc))
	if err != nil {
		panic(err)
	}

	return n
}

/* ========== Filesys ========== */
// Only confirmed to work with mac

/* ReadDependencies parses the `go mod graph` output at the given filename for consumption in this
application. */
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

/* VerifyFileStructure ensures that any prerequisites to write commands, such as directories, exists. */
func VerifyFileStructure() {
	err := os.Mkdir(DIR_NAME, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		panic(err)
	}
}

/* writeConfig writes the config to "config.json". Should writew to config any time config is
altered. */
func writeConfig() {
	config := GetConfig()
	configByte, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Dump config raw: ", config)
		panic(err)
	}

	err = os.WriteFile(fmt.Sprintf("%s%sconfig.json", DIR_NAME, string(os.PathSeparator)), configByte, 0666)
	if err != nil {
		fmt.Println("Dump config: ", string(configByte))
		panic(err)
	}
}

/* Read config enters the settings into application memory for use. */
func readConfig() {
	input, err := os.ReadFile(fmt.Sprintf("%s%sconfig.json", DIR_NAME, string(os.PathSeparator)))
	if os.IsNotExist(err) {
		return
	} else if err != nil {
		panic(err)
	}

	config := &Config{}
	err = json.Unmarshal(input, config)
	if err != nil {
		panic(err)
	}
	LoadConfig(config)
}

/* writeDependencyKey writes the dependency key info in JSON format to "dependencies.json" TODO stub*/
func writeDependencyKey(filename string, graph *Graph, edges []Edge) {

}

/* WriteHTML . TODO stub*/
func writeHtml() {

}

/* startingFs performs operations that require interaction with the filesystem on the start of the
application. TODO stub*/
func startingFs() {
	//read
	VerifyFileStructure()
	readConfig()
}

/* ========== Petty Helpers ========== */

/* AddToMap adds the dependency to the dependency-id and id-dependency maps */
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

/* PrettyListModules converts a list of module ids to a list of strings with the following format:
"#<id> [<num children>] <name of module>" */
func PrettyListModules(ids []int) []string {
	var ret []string
	for _, i := range ids {

		node, err := agraph.GetNode(i)
		if err != nil {
			panic(err)
		}

		ret = append(ret, fmt.Sprintf("#%d\t[%d]\t%s\t", node.Id, len(node.GetChildren()), im[i]))
	}
	return ret
}

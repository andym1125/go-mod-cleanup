package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var edges []Dependency

//var dependencyGraph [][]bool
var im map[int]string
var dm map[string]int
var currId int = 0

type Dependency struct {
	Module     string
	Dependency string
}

func init() {
	im = make(map[int]string)
	dm = make(map[string]int)
}

func main() {

	//ReadDependencies(os.Args[1])
	ReadDependencies("input_small.txt")

	//Register edges into maps
	for _, e := range edges {
		fmt.Println(AddToMap(e.Module), " -> ", AddToMap(e.Dependency))
	}

	//From now on, no changes to edges, graph, map, etc

	//Build Graph
	graph := New()
	for _, e := range edges {
		graph.AddDependency(fmt.Sprint(dm[e.Module]), fmt.Sprint(dm[e.Dependency]))
	}

	fmt.Println(StringifyOrderedTier(graph.Tier()))

	GenHTMLFromDependencyGraph(graph)
	//DrawDependencies("Dependencies.html", graph)
	//WriteDependencies("gomod-simple.txt", edges)
}

/* ========== Drawing ========== */

func DrawDependencies(filename string, graph DependencyGraph) {

}

/* ========== File IO ========== */

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

		edges = append(edges, Dependency{
			Module:     module,
			Dependency: dependency,
		})
	}
}

func WriteDependencies(filename string, dependencies []Dependency) {
	if len(dependencies) == 0 {
		fmt.Println("No dependencies. Not writing")
	}

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var currParent string
	for _, d := range dependencies {

		if currParent != d.Module {
			currParent = d.Module
			file.Write([]byte("\n" + d.Module + "\n"))
		}

		file.Write([]byte("\t" + d.Dependency + "\n"))
	}
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

/* ========== Standard Helpers ========== */

func StringifyBoolArr2(arr [][]bool) string {

	ret := "----- Bool 2D Arr -----\n"
	for i := range arr {
		for j := range arr[i] {

			//1/0
			if arr[i][j] {
				ret += "1"
			} else {
				ret += "0"
			}

			//Delimiter
			if j != len(arr[i])-1 {
				ret = ret + " "
			}
		}
		ret += "\n"
	}

	return ret + "----------END----------\n"
}

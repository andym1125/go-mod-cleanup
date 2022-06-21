package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

var edges []Dependency
var dependencyGraph [][]bool
var idMap map[int]string
var dependencyMap map[string]int
var currId int

type Dependency struct {
	Module     string
	Dependency string
}

func main() {

	ReadDependencies(os.Args[1])

	//Register edges into maps
	for _, e := range edges {
		fmt.Println(AddToMap(e.Dependency))
		fmt.Println(AddToMap(e.Module))
	}

	WriteDependencies("gomod-simple.txt", edges)
}

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

func AddToMap(dependency string) int {
	id, exists := dependencyMap[dependency]

	if exists {
		return id
	}

	dependencyMap[dependency] = currId
	idMap[currId] = dependency
	currId++
	return dependencyMap[dependency]
}

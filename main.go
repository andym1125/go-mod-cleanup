package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Dependency struct {
	Parent string
	Names  []string
}

func main() {
	filename := os.Args[1]
	//filename := "input.txt"
	var dependencies []Dependency

	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lineScanner := bufio.NewScanner(strings.NewReader(scanner.Text()))
		lineScanner.Split(bufio.ScanWords)

		var parent string
		if lineScanner.Scan() {
			parent = lineScanner.Text()
		} else {
			panic(errors.New("No parent"))
		}

		var names []string
		for lineScanner.Scan() {
			names = append(names, lineScanner.Text())
		}

		dependencies = append(dependencies, Dependency{
			Parent: parent,
			Names:  names,
		})
	}

	WriteDependencies("gomod-simple.txt", dependencies)
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

		if currParent != d.Parent {
			currParent = d.Parent
			file.Write([]byte("\n" + d.Parent + "\n"))
		}

		file.Write([]byte("\t" + d.Names[0] + "\n"))
	}
}

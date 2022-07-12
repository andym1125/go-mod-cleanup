package main

import (
	"fmt"
	"os"
)

func DependencySearch(rootDir string) {
	file, err := os.Open(rootDir)
	if err != nil {
		panic(err)
	}

	filestat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	if !filestat.IsDir() {
		//SkimDependency(file)
	} else {

	}
}

/* RecursiveTouch iterates through all the files rootDir recursively, executing the callback on
each non-directory file. */
func RecursiveTouch(rootDir string, callback func(file *os.File)) {
	file, err := os.Open(rootDir)
	if err != nil {
		panic(err)
	}

	filestat, err := file.Stat()
	if err != nil {
		panic(err)
	}

	if filestat.IsDir() {
		dirs, err := file.ReadDir(0)
		if err != nil {
			panic(err)
		}

		for _, d := range dirs {
			fmt.Println(d.Name())
			RecursiveTouch(rootDir+string(os.PathSeparator)+d.Name(), callback)
		}
	} else {
		callback(file)
	}
}

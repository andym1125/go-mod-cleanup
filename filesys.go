package main

import (
	"encoding/json"
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"regexp"
)

var importRegex *regexp.Regexp
var importsRegex *regexp.Regexp
var nativeImportRegex *regexp.Regexp
var githubImportRegex *regexp.Regexp
var stripQuotesRegex *regexp.Regexp

func init() {
	importRegex = regexp.MustCompile(`import\s*\(\s*(?:[^"]+"([^"]+)"\s*)+\)`)
	importsRegex = regexp.MustCompile(`"[^"]+"`)
	nativeImportRegex = regexp.MustCompile(`^[^\/]+$`)
	githubImportRegex = regexp.MustCompile(`github.com\/[^\/]+\/[^\/]+`)
	stripQuotesRegex = regexp.MustCompile(`[^"]+`)
}

func ParseImports(file *os.File) []string {
	parsedFile, err := parser.ParseFile(token.NewFileSet(), file.Name(), nil, parser.ImportsOnly)
	if err != nil {
		panic(err)
	}

	ret := make([]string, 0)
	for _, a := range parsedFile.Imports {
		ret = append(ret, a.Path.Value)
	}

	return ret
}

func DependencySearch(rootDir string) {
	output := RecursiveTouch(rootDir, func(file *os.File) string {
		ret, err := json.Marshal(ParseImports(file))
		if err != nil {
			panic(err)
		}

		return string(ret)
	})

	for _, arrstr := range output {
		var arr *[]string
		json.Unmarshal([]byte(arrstr), arr)
		fmt.Println(arr)
	}
}

/* RecursiveTouch iterates through all the files rootDir recursively, executing the callback on
each non-directory file. Callback allows for a string return which allows data to be transferred
out of the function. */
func RecursiveTouch(rootDir string, callback func(file *os.File) string) []string {
	var data []string
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
			data = append(data, RecursiveTouch(rootDir+string(os.PathSeparator)+d.Name(), callback)...)
		}
	} else {
		data = append(data, callback(file))
	}

	return data
}

/* ScrapeImports takes in an import string and adds it to the Set that was passed in. If set is
nil, a Set will be created for use only internally. */
func ScrapeImports(input string, set *Set[string]) {
	if set == nil {
		set = NewSet[string]()
	}

	importStatement := importRegex.FindString(input)
	imports := importsRegex.FindAllString(importStatement, -1)
	for _, a := range imports {
		if strippedImport := StripImport(a); strippedImport != "" {
			set.Add(strippedImport)
		}
	}
}

/* StripImport strips an import string down so it can be better compared to a dependency string.
For native packages (ie "errors"), none of the import will be preserved.
For github imports, only "github.com", the username, and the repo name will be preserved.
For all other imports, the entire string will be preserved. */
func StripImport(input string) string {
	if output := nativeImportRegex.FindString(input); output != "" {
		return ""
	} else if output := githubImportRegex.FindString(input); output != "" {
		return stripQuotesRegex.FindString(output)
	}

	return stripQuotesRegex.FindString(input)
}

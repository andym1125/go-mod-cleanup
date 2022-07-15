package main

import (
	"fmt"
	"os"
	"sort"
	"testing"
)

func TestRecursiveTouch(t *testing.T) {

	strs := RecursiveTouch("../sonr/internal/motor", func(file *os.File) string {
		filestat, _ := file.Stat()
		return "==" + filestat.Name()
	})

	fmt.Println(strs)
}

func TestParseImports(t *testing.T) {
	file, err := os.Open("/Users/andy/Documents/workshop/sonr/internal/motor/motor.go")
	if err != nil {
		panic(err)
	}
	ParseImports(file)
}

func TestDependencySearch(t *testing.T) {
	DependencySearch("/Users/andy/Documents/workshop/sonr/internal/motor")
}

func TestScrapeImports(t *testing.T) {
	set := NewSet[string]()

	expected := []string{"github.com/kataras/golog", "github.com/libp2p/go-libp2p-core",
		"github.com/sonr-io/sonr", "github.com/kataroas/golog", "random.com/sonr-io/sonr/pkg/fs"}
	sort.Strings(expected)

	ScrapeImports(`import ( "errors"  "github.com/kataras/golog" "github.com/libp2p/go-libp2p-core/protocol" device "github.com/sonr-io/sonr/pkg/fs")`, set)
	ScrapeImports(`import ( "errors"  "github.com/kataroas/golog" "github.com/libp2p/go-libp2p-core/hey" device "random.com/sonr-io/sonr/pkg/fs")`, set)
	ScrapeImports(`import ( "errors"  "github.com/kataras/golog" "github.com/libp2p/go-libp2p-core/protocol" device "github.com/sonr-io/sonr/pkg/fs")`, set)
	actual := set.ToArray()
	sort.Strings(actual)

	fmt.Println(actual)
	fmt.Println(expected)

	if len(actual) != len(expected) {
		t.FailNow()
	}
	for i := 0; i < len(actual); i++ {
		if actual[i] != expected[i] {
			t.FailNow()
		}
	}
}

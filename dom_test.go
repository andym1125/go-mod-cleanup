package main

import "testing"

func TestGom(t *testing.T) {

	dom := WrapInHtml(GetInjectGomEl("SVG-Graph-1")).Build() + GetInjectGomEl("SVG-Graph-1").Build()
	inject, _ := Inject(dom, "SVG-Graph-1", "<svg></svg>")
	t.Logf(inject)

	t.FailNow()
}

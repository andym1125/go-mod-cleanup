package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hadihammurabi/gom"
)

type BuiltElement struct {
	gom.Element
	rawHtml string
}

func NewBuiltElement(rawHtml string) BuiltElement {
	return BuiltElement{rawHtml: rawHtml}
}

func (el BuiltElement) Build() (html string) {
	return el.rawHtml
}

/* ===== Specific DOM funcs ===== */
//<button onclick="baseModule=3;changeBaseModule()">3</button>

func GenerateBasicHtml() *gom.Element {
	div := gom.Div.A(gom.Attribute{Name: "class", Value: "flex-container"}).C(
		gom.Div.A(gom.Attribute{Name: "class", Value: "base-container"}).C(
			NewBaseModuleButton(1),
		),
		gom.Div.A(gom.Attribute{Name: "class", Value: "svg-container"}).C(
			gom.H("img").A(gom.Attribute{Name: "src", Value: "https://i.etsystatic.com/14652458/r/il/32b5cc/1542847490/il_fullxfull.1542847490_snfb.jpg"}),
		),
		gom.Div.A(gom.Attribute{Name: "class", Value: "key-container"}),
	)
	return div
}

func NewBaseModuleButton(moduleId int) *gom.Element {
	button := gom.H("button").A(
		gom.Attribute{Name: "onclick", Value: fmt.Sprintf("baseModule=%d;changeBaseModule()", moduleId)},
	)
	return button
}

var StripSvgRegex *regexp.Regexp

func init() {
	StripSvgRegex = regexp.MustCompile(`<\s*svg[^<>]*>`)
}

func Inject(builtHtml string, injectTag string, injectHtml string) (string, error) {

	injectRegx := regexp.MustCompile(fmt.Sprintf(`<[^<>]*%s[^<>]*>[^<>]*<[^<>]*/[^<>]*%s[^<>]*>`, injectTag, injectTag))
	strs := injectRegx.Split(builtHtml, -1)

	if len(strs) == 0 {
		return "", errors.New("Tag to inject to not found")
	}

	ret := strs[0]
	for i := 1; i < len(strs); i++ {
		ret += injectHtml + strs[i]
	}

	return ret, nil
}

func GetInjectGomEl(injectTag string) *gom.Element {
	return gom.H(injectTag)
}

func StripSvg(svg string) string {
	loc := StripSvgRegex.FindStringIndex(svg)
	newSvg := svg[loc[0]:]
	return newSvg
}

func WrapInHtml(el *gom.Element) *gom.Element {
	dom := gom.HTML.C(
		gom.Body.C(el),
	)

	return dom
}

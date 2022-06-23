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

/* ========== */

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

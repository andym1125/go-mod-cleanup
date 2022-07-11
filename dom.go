package main

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/hadihammurabi/gom"
)

/* StripSvgRegex matches SVG tags. */
var StripSvgRegex *regexp.Regexp

func init() {
	StripSvgRegex = regexp.MustCompile(`<\s*svg[^<>]*>`)
}

/* Inject searches a builtHtml string for an injectTag, and replaces each instance with injectHtml
content. This allows for the replacing of certain tags with other content. */
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

/* GetInjectGomEl creates and returns a GOM element with the element type of injectTag. This
function allows for placeholding in GOM objects, for later injection of HTML code that, for
whatever reason, cannot be used with GOM and must be plaintext. Use this in conjunction with
Inject() */
func GetInjectGomEl(injectTag string) *gom.Element {
	return gom.H(injectTag)
}

/* StripSvg strips the SVG tags from a SVG HTML string */
func StripSvg(svg string) string {
	loc := StripSvgRegex.FindStringIndex(svg)
	newSvg := svg[loc[0]:]
	return newSvg
}

/* WrapInHtml wraps the given GOM element in a top-level HTML GOM element. */
func WrapInHtml(el *gom.Element) *gom.Element {
	dom := gom.HTML.C(
		gom.Body.C(el),
	)

	return dom
}

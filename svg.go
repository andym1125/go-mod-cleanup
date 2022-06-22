package main

import (
	"fmt"

	"github.com/hadihammurabi/gom"
)

func GenHTMLFromDependencyGraph(graph *DependencyGraph) string {

	tier, tierSize := graph.Tier()
	nodeArr := GenNodeArr(tier, tierSize)
	fmt.Println(nodeArr)

	return ""
}

func GenNodeArr(tier map[int][]string, numTiers int) [][]string {

	//Find size of largest tier
	largestTierSize := 0
	for i := 0; i <= numTiers; i++ {
		if len(tier[i]) > largestTierSize {
			largestTierSize = len(tier[i])
		}
	}

	//Hydrate Node arr, centering nodes as best as possible
	nodeArr := make([][]string, numTiers+1) //[numTiers][largestTierSize]
	for i := 0; i <= numTiers; i++ {

		var tempArr []string
		start := int((largestTierSize - len(tier[i])) / 2)
		for j := 0; j < start; j++ {
			tempArr = append(tempArr, "")
		}
		tempArr = append(tempArr, tier[i]...)
		for j := len(tempArr); j < largestTierSize; j++ {
			tempArr = append(tempArr, "")
		}
		nodeArr = append(nodeArr, tempArr)
	}

	for i := range nodeArr {
		for j := range nodeArr[i] {
			fmt.Print(nodeArr[i][j] + "\t")
		}
		fmt.Println()
	}

	return nil
}

func WrapInHtml(el *gom.Element) *gom.Element {
	html := gom.H("html").C(
		gom.H("body").C(
			el,
		),
	)

	return html
}

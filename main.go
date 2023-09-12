package main

import (
	tr "dynalist_to_markdown/transformer"
	"dynalist_to_markdown/util"
	"fmt"
	"strings"
)

func main() {
	refined := []string{}
	imageOrder := 0

	// Read target txt file.
	target, err := util.ReadLines("./mock_data.txt")
	if err != nil {
		panic(err.Error())
	}

	inCodeBlock := false
	for i, row := range target {
		refinedRow := row
		fmt.Println(row)

		// Remove title
		if i == 0 {
			continue
		}

		isStartedToMark := tr.CheckStartedToMark(row)
		isCodeMarkLine := tr.CheckCodeMarkLine(row)
		if isCodeMarkLine {
			inCodeBlock = !inCodeBlock
		}

		// Refine indent of shift+ctrl+enter new line.
		if !isStartedToMark && !isCodeMarkLine && !inCodeBlock {
			beforeRow := target[i-1]
			refinedRow = tr.AddIndentsLikeBefore(beforeRow, row)
		}

		// Remove first indent to all of contents.
		refinedRow = tr.RemoveIndents(refinedRow, 4)

		// Add double space to the end of each row
		refinedRow = tr.AddLastDoubleSpace(refinedRow)

		// Convert image phrase format
		refinedRow, imageOrder = tr.RefineImagePhrase(refinedRow, imageOrder)

		// Add '{:target="_blank"}' to the end of link phrase
		refinedRow = tr.RefineLinkPhrase(refinedRow)

		refined = append(refined, refinedRow+"\n")

	}
	fmt.Println("--------")
	fmt.Println(strings.Join(refined, ""))

	// Add meta tag to top of txt file.

}

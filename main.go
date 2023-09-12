package main

import (
	"dynalist_to_markdown/config"
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

		// Remove title and add meta phrase
		if i == 0 {
			refinedRow = strings.TrimLeft(refinedRow, " ")
			refinedRow = strings.Replace(refinedRow, config.Mark+" ", "", 1)
			refinedRow = tr.AddMeta(refinedRow)
			refined = append(refined, refinedRow+"\n")
			continue
		}

		// Add sub title
		isSubTitle := tr.CheckSubTitle(refinedRow)
		if isSubTitle {
			refinedRow = strings.Replace(refinedRow, config.Mark, "", 1)
			refinedRow = tr.RemoveIndents(refinedRow, 8)
			refinedRow = tr.AddSubTitle(refinedRow)
			refined = append(refined, "\n")
			refined = append(refined, refinedRow+"\n")
			continue
		}

		isStartedToMark := tr.CheckStartedToMark(refinedRow)
		isCodeMarkLine := tr.CheckCodeMarkLine(refinedRow)
		if isCodeMarkLine {
			inCodeBlock = !inCodeBlock
		}

		// Refine indent of shift+ctrl+enter new line.
		if !isStartedToMark && !isCodeMarkLine && !inCodeBlock {
			beforeRow := target[i-1]
			refinedRow = tr.AddIndentsLikeBefore(beforeRow, refinedRow)
		}

		// Remove 8 indents to all of contents.
		refinedRow = tr.RemoveIndents(refinedRow, 8)

		// Convert image phrase format
		refinedRow, imageOrder = tr.RefineImagePhrase(refinedRow, imageOrder)

		// Add '{:target="_blank"}' to the end of link phrase
		refinedRow = tr.RefineLinkPhrase(refinedRow)

		// Add double space to the end of each row
		refinedRow = tr.AddLastDoubleSpace(refinedRow)

		refined = append(refined, refinedRow+"\n")

	}
	fmt.Println("--------")
	fmt.Println(strings.Join(refined, ""))

	// Add meta tag to top of txt file.

}

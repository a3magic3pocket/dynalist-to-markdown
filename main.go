package main

import (
	"dynalist_to_markdown/config"
	tr "dynalist_to_markdown/transformer"
	"dynalist_to_markdown/util"
	"fmt"
	"os"
	"strings"
)

func refinedDynalist(rows []string) []string {
	refined := []string{}
	imageOrder := 0

	inCodeBlock := false
	for i, row := range rows {
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
		isRefBeforeRefinedRow := false
		if !isStartedToMark && !isCodeMarkLine && !inCodeBlock {
			beforeRefinedRow := refined[len(refined)-1]
			beforeRow := rows[i-1]
			if len(beforeRow) < len(beforeRefinedRow) {
				isRefBeforeRefinedRow = true
				beforeRow = beforeRefinedRow
			}
			refinedRow = tr.AddIndentsLikeBefore(beforeRow, refinedRow)
		}

		// Remove 8 indents to all of contents.
		if !isRefBeforeRefinedRow {
			refinedRow = tr.RemoveIndents(refinedRow, 8)
		}

		// Convert image phrase format
		if !inCodeBlock {
			refinedRow, imageOrder = tr.RefineImagePhrase(refinedRow, imageOrder)
		}

		// Add '{:target="_blank"}' to the end of link phrase
		if !inCodeBlock {
			refinedRow = tr.RefineLinkPhrase(refinedRow)
		}

		// Escape double brace row
		refinedRow = tr.EscapeDoubleBrace(refinedRow)

		// Add double space to the end of each row
		refinedRow = tr.AddLastDoubleSpace(refinedRow)

		refined = append(refined, refinedRow+"\n")
	}

	return refined
}

func help() {
	fmt.Println(`- Hwo to use?
	- build
		- go build -o dynalist_to_markdown
	- run build file
		- ./dynalist_to_markdown [dynalist file path]
	- run immediately
		- go run main.go [dynalist file path]
	- What is dynalist-to-markdown?
		- Convert the file of dynalist exported to markdown`)
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		help()
		return
	}

	sourcePath := args[0]

	// Read dynalist file exported.
	target, err := util.ReadLines(sourcePath)
	if err != nil {
		panic(err.Error())
	}

	refined := refinedDynalist(target)

	// Write refined
	resultPath := util.GetResultFilePath(sourcePath)
	err = util.WriteLines(refined, resultPath)
	if err != nil {
		panic(err.Error())
	}
}

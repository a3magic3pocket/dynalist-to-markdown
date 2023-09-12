package transformer

import (
	"dynalist_to_markdown/config"
	"dynalist_to_markdown/util"
	"fmt"
	"regexp"
	"strings"
)

func AddMeta(title string) string {
	timePhrase := util.GetTimeNowMetaFormat()

	return fmt.Sprintf(`---
title: %s
date: %s
categories: [your-category]
tags: [your-tag1, your-tage2]    # TAG names should always be lowercase
---`, title, timePhrase)
}

func CheckCodeMarkLine(row string) bool {
	return strings.Contains(row, " ```")
}

func CheckStartedToMark(row string) bool {
	trimmed := strings.TrimLeft(row, " ")

	return trimmed[0:1] == config.Mark
}

func GetNotSpaceIndex(row string) int {
	return strings.IndexFunc(row, func(r rune) bool {
		return r != ' '
	})
}

func AddIndentsLikeBefore(beforeRow string, targetRow string) string {
	beforeRow = strings.ReplaceAll(beforeRow, config.Mark, " ")
	notSpaceIndex := GetNotSpaceIndex(beforeRow)

	buf := []string{}
	for i := 0; i < notSpaceIndex; i++ {
		buf = append(buf, " ")
	}
	buf = append(buf, strings.TrimLeft(targetRow, " "))
	joined := strings.Join(buf, "")

	return joined
}

func RemoveIndents(row string, numSpace int) string {
	base := numSpace
	notSpaceIndex := GetNotSpaceIndex(row)
	if notSpaceIndex < numSpace {
		base = notSpaceIndex
	}

	return row[base:]
}

func AddLastDoubleSpace(row string) string {
	return strings.Join([]string{row, " ", " "}, "")
}

func RefineImagePhrase(row string, imageOrder int) (string, int) {
	pattern := `!\[(.*?)]\((.+?)\)`
	re := regexp.MustCompile(pattern)
	isMatched := re.MatchString(row)

	if isMatched {
		refinedRow := row
		newImageOrder := imageOrder

		for _, submatchedInfo := range re.FindAllStringSubmatch(row, -1) {
			if len(submatchedInfo) < 3 {
				continue
			}
			submatched := submatchedInfo[0]
			link := submatchedInfo[2]
			imagePath := fmt.Sprintf("%s/%s/%02d-%s", config.ImageDirPath, config.Endpoint, imageOrder, link)
			result := fmt.Sprintf(
				"<a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>",
				imagePath,
				imagePath,
				config.ImageWidth,
				config.ImageHeight)

			refinedRow = strings.Replace(refinedRow, submatched, result, 1)
			newImageOrder += 1
		}

		return refinedRow, newImageOrder
	}

	return row, imageOrder
}

func AddBlankPhrase(row string) string {
	return row + "{:target=\"_blank\"}"
}

func RefineLinkPhrase(row string) string {
	pattern := `[^!]\[.*?\]\(.+?\)|^\[.*?\]\(.+?\)`
	patternWithBlank := `[^!]\[.*?\]\(.+?\)(\s*{:target\s*=\s*["']_blank["']\s*})|^\[.*?\]\(.+?\)(\s*{:target\s*=\s*["']_blank["']\s*})`
	re := regexp.MustCompile(pattern)
	re2 := regexp.MustCompile(patternWithBlank)
	isMatched := re.MatchString(row)
	isMatchedWithBlank := re2.MatchString(row)

	// Remove blank phrases
	if isMatchedWithBlank {
		for _, subMatchInfo := range re2.FindAllStringSubmatch(row, -1) {
			if len(subMatchInfo) < 3 {
				continue
			}

			blankPhrase := subMatchInfo[1]
			if blankPhrase == "" {
				blankPhrase = subMatchInfo[2]
			}

			row = strings.ReplaceAll(row, blankPhrase, "")
		}
	}

	// Add blank phrases
	if isMatched {
		refinedRow := row

		for _, submatchedInfo := range re.FindAllStringSubmatch(row, -1) {
			if len(submatchedInfo) < 1 {
				continue
			}
			linkPhrase := submatchedInfo[0]
			refinedRow = strings.Replace(refinedRow, linkPhrase, AddBlankPhrase(linkPhrase), 1)
		}

		return refinedRow
	}

	return row
}

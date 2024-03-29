package transformer

import (
	"dynalist_to_markdown/config"
	"fmt"
	"strings"
	"testing"
)

func TestCheckCodeMarkLine(t *testing.T) {
	// case 1
	codeMark := " ```"
	target := fmt.Sprintf("%sbash", codeMark)
	result := CheckCodeMarkLine(target)
	if !result {
		t.Errorf("'%s' is contained\n", codeMark)
	}

	// case 2
	target = fmt.Sprintf("%s한글", codeMark)
	result = CheckCodeMarkLine(target)
	if !result {
		t.Errorf("'%s' is contained\n", codeMark)
	}

	// case 3
	target = "it is not contained"
	result = CheckCodeMarkLine(target)
	if result {
		t.Errorf("'%s' is not contained\n", codeMark)
	}
}

func TestCheckStartedToMark(t *testing.T) {
	// case 1
	target := fmt.Sprintf("%s it is title", config.Mark)
	result := CheckStartedToMark(target)
	if !result {
		t.Errorf("'%s' starts on a row\n", config.Mark)
	}

	// case 2
	target = fmt.Sprintf("%s제목제목제목", config.Mark)
	result = CheckStartedToMark(target)
	if !result {
		t.Errorf("'%s' starts on a row\n", config.Mark)
	}

	// case 3
	target = "it is not contained"
	result = CheckStartedToMark(target)
	if result {
		t.Errorf("'%s' not exists\n", config.Mark)
	}

	// case 4
	target = fmt.Sprintf("the end of row %s", config.Mark)
	result = CheckStartedToMark(target)
	if result {
		t.Errorf("'%s' does not start on a row\n", config.Mark)
	}

	// case 5
	target = fmt.Sprintf("the middle of row %s |||", config.Mark)
	result = CheckStartedToMark(target)
	if result {
		t.Errorf("'%s' does not start on a row\n", config.Mark)
	}

	// case 6
	target = "          "
	result = CheckStartedToMark(target)
	if result {
		t.Errorf("'%s' does not start on a row\n", config.Mark)
	}
}

func TestGetNotSpaceIndex(t *testing.T) {
	// case 1
	target := "       answer is 7"
	result := GetNotSpaceIndex(target)
	if result != 7 {
		t.Errorf("NotSpaceIndex is not %d\n", 7)
	}

	// case 2
	target = "answer is 0"
	result = GetNotSpaceIndex(target)
	if result != 0 {
		t.Error("NotSpaceIndex is not 0")
	}
}

func TestAddIndentsLikeBefore(t *testing.T) {
	// case 1
	targetBeforeRow := fmt.Sprintf("    %s 4 indents with mark", config.Mark)
	targetRow := "                     dummy"
	result := AddIndentsLikeBefore(targetBeforeRow, targetRow)
	answer := "      dummy"
	if result != answer {
		t.Errorf("failed to run AddIndentsLikeBefore with '4 indents with mark'\n")
	}

	// case 2
	targetBeforeRow = "    4 indents without mark"
	targetRow = "                     dummy"
	result = AddIndentsLikeBefore(targetBeforeRow, targetRow)
	answer = "    dummy"
	if result != answer {
		t.Errorf("failed to run AddIndentsLikeBefore with '4 indents without mark'\n")
	}

	// case 3
	targetBeforeRow = "  2 indents(less than 4) without mark"
	targetRow = "                     dummy"
	result = AddIndentsLikeBefore(targetBeforeRow, targetRow)
	answer = "  dummy"
	if result != answer {
		t.Errorf("failed to run AddIndentsLikeBefore with '2 indents(less than 4) without mark'\n")
	}
}

func TestRemoveIndents(t *testing.T) {
	// case 1
	target := "       7 indents"
	result := RemoveIndents(target, 4)
	answer := "   7 indents"
	if result != answer {
		t.Errorf("failed to run TestRemoveIndents with '7 indents'\n")
	}

	// case 2
	target = "  2 indents(less than base 4)"
	result = RemoveIndents(target, 4)
	answer = "2 indents(less than base 4)"
	if result != answer {
		t.Errorf("failed to run TestRemoveIndents with '2 indents(less than base 4)'\n")
	}

	// case 3
	target = "     5 indents and base is 5"
	result = RemoveIndents(target, 5)
	answer = "5 indents and base is 5"
	if result != answer {
		t.Errorf("failed to run TestRemoveIndents with '5 indents and base is 5'\n")
	}

	// case 5
	target = "       "
	result = RemoveIndents(target, 5)
	answer = "       "
	if result != answer {
		t.Errorf("failed to run TestRemoveIndents with 'all spaces'\n")
	}
}

func TestAddLastDoubleSpace(t *testing.T) {
	// case 1
	target := "without double space"
	result := AddLastDoubleSpace(target)
	answer := "without double space  "
	if result != answer {
		t.Errorf("failed to run TestRemoveIndents with 'without double space'\n")
	}
}

func TestRefineImagePhrase(t *testing.T) {
	// case 1
	target := `![test_image](https://case-1.jpeg)`
	imageOrder := 3
	result, newImageOrder := RefineImagePhrase(target, imageOrder)
	imagePath := fmt.Sprintf("%s/%s/%02d-https://case-1.jpeg",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	answer := fmt.Sprintf(
		"<a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>",
		imagePath,
		imagePath,
		config.ImageWidth,
		config.ImageHeight)
	if result != answer {
		t.Errorf("failed to run RefineImagePhrase with normal case'\n")
	}
	if newImageOrder != imageOrder+1 {
		t.Errorf("imageOrder is not updated on normal case'\n")
	}

	// case 2
	target = "without image"
	imageOrder = 1
	result, newImageOrder = RefineImagePhrase(target, imageOrder)
	if result != target {
		t.Errorf("failed to run RefineImagePhrase with 'without image'\n")
	}
	if newImageOrder != imageOrder {
		t.Errorf("imageOrder is updated on 'without image'\n")
	}

	// case 3
	target = `IN MIDDLE || ![test_image](https://case-1.jpeg)||`
	imageOrder = 4
	result, newImageOrder = RefineImagePhrase(target, imageOrder)
	imagePath = fmt.Sprintf("%s/%s/%02d-https://case-1.jpeg",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	answer = fmt.Sprintf(
		"IN MIDDLE || <a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>||",
		imagePath,
		imagePath,
		config.ImageWidth,
		config.ImageHeight)
	if result != answer {
		t.Errorf("failed to run RefineImagePhrase with 'image in middle of row'\n")
	}
	if newImageOrder != imageOrder+1 {
		t.Errorf("imageOrder is not updated on 'image in middle of row'\n")
	}

	// case 4
	target = `IN MIDDLE || ![test_image1](https://case-1.jpeg)|| ![test_image2](https://case-2.png)`
	imageOrder = 5
	result, newImageOrder = RefineImagePhrase(target, imageOrder)
	imagePath1 := fmt.Sprintf("%s/%s/%02d-https://case-1.jpeg",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	imagePath2 := fmt.Sprintf("%s/%s/%02d-https://case-2.png",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	answer = fmt.Sprintf(
		"IN MIDDLE || <a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>|| <a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>",
		imagePath1,
		imagePath1,
		config.ImageWidth,
		config.ImageHeight,
		imagePath2,
		imagePath2,
		config.ImageWidth,
		config.ImageHeight)
	if result != answer {
		t.Errorf("failed to run RefineImagePhrase with 'multiple images in row'\n")
	}
	if newImageOrder != imageOrder+2 {
		t.Errorf("imageOrder is not updated on 'multiple images in row'\n")
	}

	// case 5
	target = `IN MIDDLE || ![test_image1](https://case-1.jpeg)|| ![test_image2](https://case-2.png)|| [test_linl](https://hello-world.com)`
	imageOrder = 5
	result, newImageOrder = RefineImagePhrase(target, imageOrder)
	imagePath1 = fmt.Sprintf("%s/%s/%02d-https://case-1.jpeg",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	imagePath2 = fmt.Sprintf("%s/%s/%02d-https://case-2.png",
		config.ImageDirPath,
		config.Endpoint,
		imageOrder)
	answer = fmt.Sprintf(
		"IN MIDDLE || <a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>|| <a href='%s' target='_blank'><img src='%s' width='%d%%' height='%d%%'></a>|| [test_linl](https://hello-world.com)",
		imagePath1,
		imagePath1,
		config.ImageWidth,
		config.ImageHeight,
		imagePath2,
		imagePath2,
		config.ImageWidth,
		config.ImageHeight)
	if result != answer {
		t.Errorf("failed to run RefineImagePhrase with 'multiple images with link in row'\n")
	}
	if newImageOrder != imageOrder+2 {
		t.Errorf("imageOrder is not updated on 'multiple images with link in row'\n")
	}
}

func TestAddBlankPhrase(t *testing.T) {
	// case 1
	target := "hello world"
	result := AddBlankPhrase(target)
	answer := "hello world{:target=\"_blank\"}"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase'\n")
	}
}

func TestRefineLinkPhrase(t *testing.T) {
	// case 1
	target := "[test_link](https://hello-world.com)"
	result := RefineLinkPhrase(target)
	answer := "[test_link](https://hello-world.com){:target=\"_blank\"}"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase without blank phrase'\n")
	}

	// case 2
	target = "[test_link](https://hello-world.com){:target"
	result = RefineLinkPhrase(target)
	answer = "[test_link](https://hello-world.com){:target=\"_blank\"}{:target"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase with the part of blank phrase'\n")
	}

	// case 3
	target = "[test_link](https://hello-world.com){:target=\"_blank\"}"
	result = RefineLinkPhrase(target)
	answer = "[test_link](https://hello-world.com){:target=\"_blank\"}"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase with blank phrase'\n")
	}

	// case 4
	target = "![test_image](https://hello-world.jpg)"
	result = RefineLinkPhrase(target)
	answer = "![test_image](https://hello-world.jpg)"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase with image'\n")
	}

	// case 5
	target = "asdf[test_link1](https://hello-world1.com){:target=\"_blank\"}||[test_link2](https://hello-world2.com)bbb"
	result = RefineLinkPhrase(target)
	answer = "asdf[test_link1](https://hello-world1.com){:target=\"_blank\"}||[test_link2](https://hello-world2.com){:target=\"_blank\"}bbb"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase with multiple blank phrases'\n")
	}

	// case 6
	target = "asdf[test_link1](https://hello-world1.com){:target=\"_blank\"}||[test_link2](https://hello-world2.com)bbb|Zxc![test_image](https://hello-world.jpg)"
	result = RefineLinkPhrase(target)
	answer = "asdf[test_link1](https://hello-world1.com){:target=\"_blank\"}||[test_link2](https://hello-world2.com){:target=\"_blank\"}bbb|Zxc![test_image](https://hello-world.jpg)"
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase with multiple blank phrases and image'\n")
	}
}

func TestAddMeta(t *testing.T) {
	// case 1
	result := AddMeta("hello")
	splitted := strings.Split(result, "\n")
	fmt.Println("result", result)
	fmt.Println("splitted", splitted)
	if len(splitted) != 6 {
		t.Error("the length of splitted must be 6")
	}

	isFirstLine := splitted[0] == "---"
	if !isFirstLine {
		t.Error("first line is wrong")
	}

	isOkTitle := splitted[1] == "title: hello"
	if !isOkTitle {
		t.Error("title is wrong")
	}

	isDateOk := splitted[2][0:6] == "date: "
	if !isDateOk {
		t.Error("date is wrong")
	}

	isOkCate := splitted[3] == "categories: [your-category]"
	if !isOkCate {
		t.Error("category is wrong")
	}

	isOkTags := splitted[4] == "tags: [your-tag1, your-tage2]    # TAG names should always be lowercase"
	if !isOkTags {
		t.Error("tags is wrong")
	}
}

func TestCheckSubTitle(t *testing.T) {
	// case 1
	target := "asdfa"
	isSubTitle := CheckSubTitle(target)
	if isSubTitle {
		t.Errorf("%s is not subtitle", target)
	}

	target = " - asdfa"
	isSubTitle = CheckSubTitle(target)
	if isSubTitle {
		t.Errorf("%s is not subtitle", target)
	}

	target = "               - asdfa"
	isSubTitle = CheckSubTitle(target)
	if isSubTitle {
		t.Errorf("%s is not subtitle", target)
	}

	target = "    - asdfa"
	isSubTitle = CheckSubTitle(target)
	if !isSubTitle {
		t.Errorf("%s is subtitle", target)
	}
}

func TestAddSubTitle(t *testing.T) {
	// case 1
	target := "asdf"
	result := AddSubTitle(target)
	answer := "## asdf"
	if result != answer {
		t.Errorf("%s is not %s", target, answer)
	}
}

func TestEscapeDoubleBrace(t *testing.T) {
	// case 1
	target := "style={{ backgroundImage: `url('${session.user.image}')`}}"
	result := EscapeDoubleBrace(target)
	answer := "{% raw %}style={{ backgroundImage: `url('${session.user.image}')`}}{% endraw %}"
	if result != answer {
		t.Errorf("%s is not %s", target, answer)
	}

	// case 2
	target = "is not contained double brace"
	result = EscapeDoubleBrace(target)
	answer = "is not contained double brace"
	if result != answer {
		t.Errorf("%s is not %s", target, answer)
	}

	// case 3
	target = "{{ left side double brace"
	result = EscapeDoubleBrace(target)
	answer = "{% raw %}{{ left side double brace{% endraw %}"
	if result != answer {
		t.Errorf("%s is not %s", target, answer)
	}

	// case 4
	target = "right side double brace}}"
	result = EscapeDoubleBrace(target)
	answer = "{% raw %}right side double brace}}{% endraw %}"
	if result != answer {
		t.Errorf("%s is not %s", target, answer)
	}
}

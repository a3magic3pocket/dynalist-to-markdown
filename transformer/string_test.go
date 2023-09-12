package transformer

import (
	"dynalist_to_markdown/config"
	"fmt"
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
	// target := "[test_linl](https://hello-world.com)"
	// target := "[test_linl](https://hello-world.com){:target"
	target := "[test_linl](https://hello-world.com){:target=\"_blank\"}"
	result := RefineLinkPhrase(target)
	answer := "[test_linl](https://hello-world.com){:target=\"_blank\"}"
	fmt.Println("target", target)
	fmt.Println("answer", answer)
	if result != answer {
		t.Errorf("failed to run AddBlankPhrase'\n")
	}
	t.Error("asdfasfd")
}

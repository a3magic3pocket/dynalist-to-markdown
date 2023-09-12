package util

import (
	"dynalist_to_markdown/config"
	"fmt"
	"os"
	"testing"
)

func TestReadLinesAndWriteLines(t *testing.T) {
	// case 1
	targetPath := "./test_dummy.txt"
	targetData := []string{
		"hello\n",
		"world\n",
	}
	err := WriteLines(targetData, targetPath)
	if err != nil {
		t.Error("failed to create dummy file")
	}

	result, err := ReadLines(targetPath)
	if err != nil {
		t.Error("failed to run ReadLines")
	}
	if result[0] != "hello" {
		t.Error("failed to create dummy file")
	}

	err = os.Remove(targetPath)
	if err != nil {
		t.Error("failed to remove dummy file")
	}
}

func TestGetResultFilePath(t *testing.T) {
	// case 1
	target := "/hello/world/myfile.txt"
	result := GetResultFilePath(target)
	answer := fmt.Sprintf("/hello/world/myfile%s.md", config.ResultFileSuffix)
	if result != answer {
		t.Error("failed to run GetResultFilePath")
	}

	// case 2
	target = "/hello/world/myfile"
	result = GetResultFilePath(target)
	answer = fmt.Sprintf("/hello/world/myfile%s.md", config.ResultFileSuffix)
	if result != answer {
		t.Error("failed to run GetResultFilePath")
	}
}

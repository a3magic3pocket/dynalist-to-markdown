package util

import (
	"bufio"
	"os"
)

type FileLine struct {
	Line string
	Err  error
}

func ReadLines(filePath string) (result []string, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}, err
	}

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, nil
}

package util

import (
	"bufio"
	"dynalist_to_markdown/config"
	"fmt"
	"os"
	"strings"
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
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, nil
}

func WriteLines(rows []string, filePath string) (err error) {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		_, err = writer.WriteString(row)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetResultFilePath(filePath string) string {
	splitted := strings.Split(filePath, "/")
	parentPath := strings.Join(splitted[:len(splitted)-1], "/")

	originFileName := splitted[len(splitted)-1]
	splitted = strings.Split(originFileName, ".")

	fileNameWithoutExt := splitted[0]
	if len(splitted) > 1 {
		fileNameWithoutExt = strings.Join(splitted[:len(splitted)-1], ".")
	}

	ext := "md"

	return fmt.Sprintf("%s/%s%s.%s", parentPath, fileNameWithoutExt, config.ResultFileSuffix, ext)
}

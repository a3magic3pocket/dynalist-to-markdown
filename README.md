# dynalist-to-markdown
Convert the file of dynalist exported to markdown

## How to run
- ```bash
	# build
	go build -o dynalist_to_markdown

	# run build file
	./dynalist_to_markdown [dynalist file path]

	# run immediately
	go run main.go [dynalist file path]
    ```

## Install dependencies
- ```bash
    go mod tidy
    ```

## Test
- ```bash
    go test ./...
    ```
default: build

build:
	go build -o scrapenstein main.go

lint: fmt

fmt:
	gofmt -w .

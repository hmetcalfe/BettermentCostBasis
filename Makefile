export GO111MODULE=on

all : build lint test clean

build:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 go build -o bin/bettermentCostBasis.mac cmd/bettermentCostBasis.go
	GOOS=linux GOARCH=amd64 go build -o bin/bettermentCostBasis cmd/bettermentCostBasis.go
	GOOS=windows GOARCH=amd64 go build -o bin/bettermentCostBasis.exe cmd/bettermentCostBasis.go

tidy:
	go run tidy

lint:
	golangci-lint run --verbose --fix --deadline=5m

test:
	go test -coverprofile=coverage.out ./...

clean:
	go clean

.PHONY: all build lint test clean
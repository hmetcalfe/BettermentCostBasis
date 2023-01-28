GO = GO111MODULE=on go

all : build lint test clean

build:
	mkdir -p bin
	GOOS=darwin GOARCH=amd64 $(GO) build -o bin/bettermentCostBasis.mac cmd/bettermentCostBasis.go
	GOOS=linux GOARCH=amd64 $(GO) build -o bin/bettermentCostBasis cmd/bettermentCostBasis.go
	GOOS=windows GOARCH=amd64 $(GO) build -o bin/bettermentCostBasis.exe cmd/bettermentCostBasis.go

tidy:
	@$(GO) mod tidy

lint:
	golangci-lint run --verbose --fix --deadline=5m

test:
	@$(GO) test -coverprofile=coverage.out ./...

clean:
	@$(GO) clean

.PHONY: all build lint test clean
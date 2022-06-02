GOLINT := golangci-lint

.PHONY: build test

build:
	go build -o conduit-connector-http-server cmd/http-server/main.go

test:
	go test $(GOTEST_FLAGS) -race ./...

lint:
	$(GOLINT) run --timeout=5m -c .golangci.yml


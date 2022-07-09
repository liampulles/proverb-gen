# Init variables
GOBIN := $(shell go env GOPATH)/bin

# Keep this at the top so that it is default when `make` is called.
# This is used by Travis CI.
coverage.txt:
	go test -race -coverprofile=coverage.txt.tmp -covermode=atomic ./...

view-cover: clean coverage.txt
	go tool cover -html=coverage.txt
test: build
	go test ./...
build:
	go build ./...
install: build
	go install ./...
inspect: build $(GOBIN)/golangci-lint
	$(GOBIN)/golangci-lint run
update:
	go get -u ./...
pre-commit: update clean coverage.txt inspect
	go mod tidy
clean:
	rm -f coverage.txt $(GOBIN)/proverb-gen

# Needed tools
$(GOBIN)/golangci-lint:
	$(MAKE) install-tools
install-tools:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOBIN) v1.46.2
	rm -rf ./v1.46.2
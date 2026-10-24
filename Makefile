.PHONY: help build test run fmt check demo install

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

build: ## build the binary
	mkdir -p target
	go build -o target/jira-issue-selector ./cmd/jira-issue-selector

test: ## run all tests
	go test -v ./...

run: build ## run the built binary
	./target/jira-issue-selector

fmt: ## format code
	go fmt ./...

check: build test ## full verification (build + test)
	@echo "ALL CHECKS PASSED"

demo: build ## run binary against local demo server
	JIRA_USER=user JIRA_API_KEY=key JIRA_HOSTNAME=http://localhost:8090 ./target/jira-issue-selector

install: build ## build and install to ~/.local/bin
	mkdir -p "$${HOME}/.local/bin"
	cp target/jira-issue-selector "$${HOME}/.local/bin"
	@echo "Installed to $${HOME}/.local/bin"

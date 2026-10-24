.PHONY: help help-all build test run fmt check demo install

.DEFAULT_GOAL := help

help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "  check    - Full verification (build + test)"
	@echo "  build    - Build the binary"
	@echo "  test     - Run all tests"
	@echo "  run      - Build and run the binary"
	@echo "  install  - Build and install to ~/.local/bin"
	@echo ""
	@echo "Run 'make help-all' for the full list, including internal targets."

help-all:
	@echo "Public targets:"
	@echo "  check    - Full verification (build + test)"
	@echo "  build    - Build the binary"
	@echo "  test     - Run all tests"
	@echo "  run      - Build and run the binary"
	@echo "  install  - Build and install to ~/.local/bin"
	@echo ""
	@echo "Internal targets (typically called via check or build):"
	@echo "  fmt      - Format code"
	@echo "  demo     - Run binary against local demo server"

build:
	mkdir -p target
	go build -o target/jira-issue-selector ./cmd/jira-issue-selector

test:
	go test -v ./...

run: build
	./target/jira-issue-selector

fmt:
	go fmt ./...

check: build test
	@echo "ALL CHECKS PASSED"

demo: build
	JIRA_USER=user JIRA_API_KEY=key JIRA_HOSTNAME=http://localhost:8090 ./target/jira-issue-selector

install: build
	mkdir -p "$${HOME}/.local/bin"
	cp target/jira-issue-selector "$${HOME}/.local/bin"
	@echo "Installed to $${HOME}/.local/bin"

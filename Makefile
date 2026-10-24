.PHONY: help help-all build test run fmt check demo install

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		grep -E "^(check|build|test|run|install):" | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""
	@echo "Run 'make help-all' for the full list, including internal targets."

help-all:
	@echo "Public targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		grep -E "^(check|build|test|run|install):" | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'
	@echo ""
	@echo "Internal targets (typically called via check or build):"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		grep -vE "^(help|help-all|check|build|test|run|install):" | \
		grep -E '^[a-zA-Z_-]+:.*?## .*$$' | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

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

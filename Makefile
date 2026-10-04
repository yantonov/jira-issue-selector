.PHONY: help build test run fmt check demo install

.DEFAULT_GOAL := help

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-12s %s\n", $$1, $$2}'

build: ## build the binary
	bin/build.sh

test: ## run all tests
	bin/test.sh

run: ## run the built binary
	bin/run.sh

fmt: ## format code
	bin/fmt.sh

check: build test ## full verification (build + test)
	@echo "ALL CHECKS PASSED"

demo: ## start demo JIRA server
	bin/start-demo.sh

install: ## build and install to ~/.local/bin
	bin/install-from-source.sh

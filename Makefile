.PHONY: build test run fmt check demo install

build:
	bin/build.sh

test:
	bin/test.sh

run:
	bin/run.sh

fmt:
	bin/fmt.sh

# full verification: build + test
check: build test
	@echo "ALL CHECKS PASSED"

demo:
	bin/start-demo.sh

install:
	bin/install-from-source.sh

.PHONY: run-examples
run-examples:
	@go test -tags all -v -run "^Test_Example*"  ./...

.PHONY: run-tests
run-tests:
	@go test -tags all -v  ./...

.PHONY: run-lint
run-lint:
	@golangci-lint run
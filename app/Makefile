.PHONY: test-analyzer
test-analyzer:
	go test -v ./analyzer/...

# NOTE: golangci-lint はバイナリインストールが推奨されている
# https://golangci-lint.run/docs/welcome/install/#install-from-sources
.PHONY: lint-analyzer
lint-analyzer:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	golangci-lint run

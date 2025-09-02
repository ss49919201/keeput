.PHONY: test-go
test-go:
	go test ./analyzer/...

# NOTE: golangci-lint はバイナリインストールが推奨されている
# https://golangci-lint.run/docs/welcome/install/#install-from-sources
.PHONY: lint-go
lint-go:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	golangci-lint run

AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query "Account" --output text)
ANALYZER_FUNCTION_NAME ?= keeput-analyzer-sandbox
ANALYZER_REPOSITORY ?= keeput-analyzer-lambda-sandbox
ANALYZER_IMAGE_TAG ?= latest
ANALYZER_IMAGE ?= ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/${ANALYZER_REPOSITORY}:${ANALYZER_IMAGE_TAG}
ANALYZER_LAMBROLL_TFSTATE_URL ?= s3://keeput-tf-state/sandbox/analyzer/service/terraform.tfstate

.PHONY: test-go
test-go:
	go test ./analyzer/...

# NOTE: golangci-lint はバイナリインストールが推奨されている
# https://golangci-lint.run/docs/welcome/install/#install-from-sources
.PHONY: lint-go
lint-go:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	golangci-lint run

.PHONY: build-analyzer
build-analyzer:
	docker buildx build --platform linux/amd64 \
	--secret type=env,id=AWS_ACCESS_KEY_ID,env=AWS_ACCESS_KEY_ID \
	--secret type=env,id=AWS_SECRET_ACCESS_KEY,env=AWS_SECRET_ACCESS_KEY \
	--secret type=env,id=AWS_SESSION_TOKEN,env=AWS_SESSION_TOKEN \
	--provenance=false \
	-t ${ANALYZER_IMAGE} .

.PHONY: deploy-analyzer
deploy-analyzer: build-analyzer
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com
	docker push ${ANALYZER_IMAGE}
	FUNCTION_NAME=${ANALYZER_FUNCTION_NAME} IMAGE_URI=${ANALYZER_IMAGE} lambroll deploy --tfstate ${ANALYZER_LAMBROLL_TFSTATE_URL}

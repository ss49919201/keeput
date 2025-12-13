AWS_ACCOUNT_ID := $(shell aws sts get-caller-identity --query "Account" --output text)
FUNCTION_NAME ?= keeput-analyzer-sandbox
REPOSITORY ?= keeput-analyzer-lambda-sandbox
IMAGE_TAG ?= latest
IMAGE ?= ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com/${REPOSITORY}:${IMAGE_TAG}
LAMBROLL_TFSTATE_URL ?= s3://keeput-tf-state/sandbox/analyzer/service/terraform.tfstate

.PHONY: test
test:
	go test ./...

# NOTE: golangci-lint はバイナリインストールが推奨されている
# https://golangci-lint.run/docs/welcome/install/#install-from-sources
.PHONY: lint
lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/HEAD/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v2.4.0
	golangci-lint run

.PHONY: build
build:
	docker buildx build --platform linux/amd64 \
	--secret type=env,id=AWS_ACCESS_KEY_ID,env=AWS_ACCESS_KEY_ID \
	--secret type=env,id=AWS_SECRET_ACCESS_KEY,env=AWS_SECRET_ACCESS_KEY \
	--secret type=env,id=AWS_SESSION_TOKEN,env=AWS_SESSION_TOKEN \
	--provenance=false \
	-t ${IMAGE} .

.PHONY: deploy
deploy: build
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-1.amazonaws.com
	docker push ${IMAGE}
	FUNCTION_NAME=${FUNCTION_NAME} IMAGE_URI=${IMAGE} lambroll deploy --tfstate ${LAMBROLL_TFSTATE_URL}

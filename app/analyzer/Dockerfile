FROM golang:1.25.5-alpine3.23 AS build
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -ldflags '-w' -trimpath -o main ./analyzer/cmd/awslambda/main.go

FROM amazon/aws-cli AS otelcollectorlayer
ARG LAYER_ARN=arn:aws:lambda:ap-northeast-1:901920570463:layer:aws-otel-collector-amd64-ver-0-117-0:1
ENV AWS_DEFAULT_REGION=ap-northeast-1
WORKDIR /opt
RUN --mount=type=secret,id=AWS_ACCESS_KEY_ID,env=AWS_ACCESS_KEY_ID \
    --mount=type=secret,id=AWS_SECRET_ACCESS_KEY,env=AWS_SECRET_ACCESS_KEY \
    --mount=type=secret,id=AWS_SESSION_TOKEN,env=AWS_SESSION_TOKEN \
    yum -y install unzip && \
    curl $(aws lambda get-layer-version-by-arn --arn ${LAYER_ARN} --query 'Content.Location' --output text) --output layer.zip && \
    unzip layer.zip && rm -f layer.zip

FROM gcr.io/distroless/static-debian12
WORKDIR /opt
COPY --from=otelcollectorlayer /opt .
COPY otel-config.yaml /opt/collector-config/config.yaml
COPY --from=build /app/main /main
ENTRYPOINT [ "/main" ]

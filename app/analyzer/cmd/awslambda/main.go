package main

import (
	"context"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
)

func handleRequest(ctx context.Context) error {
	log.Print("Hello")
	return nil
}

func main() {
	lambda.Start(handleRequest)
}

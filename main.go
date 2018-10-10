package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mozamimy/lambda-github-issue-opener/handler"
)

func main() {
	lambda.Start(handler.HandleRequest)
}

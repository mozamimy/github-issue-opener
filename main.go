package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/mozamimy/github-issue-opener/handler"
)

func main() {
	lambda.Start(handler.HandleRequest)
}

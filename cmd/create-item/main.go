package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/handlers"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	lambda.Start(handlers.HandleCreateItem(logger))
}

package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/handlers"
)

var (
	logger         *slog.Logger
	dynamodbclient *dynamodb.Client
	cfg            handlers.HandleReadItemConfig
)

const (
	envTableName = "TABLE_NAME"
)

func init() {
	var missingcfg []string

	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load default AWS config
	// If loading fails log an error and exit
	awscfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load default aws config: %s", err.Error())
		os.Exit(1)
		return
	}

	// Initialize dynamodb client
	dynamodbclient = dynamodb.NewFromConfig(awscfg)

	// Load environment variables
	cfg.TableName = os.Getenv(envTableName)
	if cfg.TableName == "" {
		missingcfg = append(missingcfg, envTableName)
	}

	// If any environment variables are missing log an error and exit
	if len(missingcfg) > 0 {
		fmt.Fprintf(
			os.Stderr,
			"failed to load the following environment variables: %s",
			strings.Join(missingcfg, ", "),
		)
		os.Exit(1)
		return
	}
}

func main() {
	lambda.Start(handlers.HandleReadItem(
		cfg,
		logger,
		dynamodbclient,
	))
}

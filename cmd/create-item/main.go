package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/config"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/handlers"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/middleware"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/services"
)

const (
	envTableName = "TABLE_NAME"
	resolverName = "createItem"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "run failed: %s", err.Error())
	}
}

func run() error {
	var cfg config.Config
	var missingcfg []string

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Load default AWS config
	// If loading fails log an error and exit
	awscfg, err := awsconfig.LoadDefaultConfig(context.Background())
	if err != nil {
		return fmt.Errorf("failed to load default aws config: %w", err)
	}

	// Initialize dynamodb client
	dynamodbclient := dynamodb.NewFromConfig(awscfg)

	// Load environment variables
	cfg.TableName = os.Getenv(envTableName)
	if cfg.TableName == "" {
		missingcfg = append(missingcfg, envTableName)
	}

	// If any environment variables are missing log an error and exit
	if len(missingcfg) > 0 {
		return fmt.Errorf(
			"failed to load the following environment variables: %s",
			strings.Join(missingcfg, ", "),
		)
	}

	service := services.NewItem(
		logger,
		dynamodbclient,
		cfg,
	)

	lambda.Start(
		middleware.Apply(
			handlers.HandleCreateItem(service),
			middleware.WithRecovery[handlers.Request[handlers.CreateItemRequest], handlers.CreateItemResponse](
				logger,
				resolverName,
			),
			middleware.WithErrorLogging[handlers.Request[handlers.CreateItemRequest], handlers.CreateItemResponse](
				logger,
				resolverName,
			),
		))

	return nil
}

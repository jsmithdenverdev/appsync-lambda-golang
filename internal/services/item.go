package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/config"
	"github.com/jsmithdenverdev/appsync-lambda-golang/internal/models"
)

type Item struct {
	logger         *slog.Logger
	dynamodbclient *dynamodb.Client
	config         config.Config
}

func NewItem(
	logger *slog.Logger,
	dynamodbclient *dynamodb.Client,
	config config.Config) Item {
	return Item{
		logger:         logger,
		dynamodbclient: dynamodbclient,
		config:         config,
	}
}

type CreateItemRequest struct {
	Name string
	Tags []string
}

func (service Item) CreateItem(ctx context.Context, request CreateItemRequest) (models.Item, error) {
	var item models.Item
	id := uuid.New().String()
	av, err := attributevalue.MarshalMap(models.Item{
		ID:   id,
		Name: request.Name,
		Tags: request.Tags,
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "failed to marshal item to attribute value map", "error", err)
		return item, fmt.Errorf("failed to marshal item to attribute value map: %w", err)
	}

	_, err = service.dynamodbclient.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(service.config.TableName),
		Item:      av,
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "failed to put item", "error", err)
		return item, fmt.Errorf("failed to put item: %w", err)
	}

	row, err := service.dynamodbclient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(service.config.TableName),
		Key: map[string]types.AttributeValue{
			"id": &types.AttributeValueMemberS{
				Value: id,
			},
		},
	})

	if err != nil {
		service.logger.ErrorContext(ctx, "failed to get item", "error", err)
		return item, fmt.Errorf("failed to get item: %w", err)
	}

	if err = attributevalue.UnmarshalMap(row.Item, &item); err != nil {
		return item, fmt.Errorf("failed to unmarshal dynamodb attribute value map: %w", err)
	}

	return item, nil
}

type BatchReadItemRequest struct {
	ID string
}

func (service Item) BatchReadItem() {}

func (service Item) BatchListItems() {}

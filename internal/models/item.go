package models

type Item struct {
	ID   string `dynamodbav:"id"`
	Name string `dynamodbav:"name"`
}

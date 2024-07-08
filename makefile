build-GraphQLLambdaFunction:
	GOOS=linux GOARCH=amd64 go build -o ./cmd/function/bootstrap ./cmd/function/main.go

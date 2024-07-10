package requests

import "context"

type CreateItem struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (request CreateItem) Valid(ctx context.Context) (problems map[string]string) {
	problems = make(map[string]string)
	if request.Name == "" {
		problems["Name"] = "request name cannot be empty"
	}
	return problems
}

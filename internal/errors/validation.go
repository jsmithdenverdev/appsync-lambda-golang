package errors

type Validation struct {
	Problems []struct {
		Key   string
		Value string
	}
}

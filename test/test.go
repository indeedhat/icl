package test

type marshalTest struct {
	target any
	output string
	error  string
}

type unmarshalTest struct {
	document string
	output   any
	error    string
}

func ptr[T any](v T) *T {
	return &v
}

package dtos

type ErrorDTO struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorsDTO struct {
	Errors []ErrorDTO `json:"errors"`
}

package dtos

type ErrorDTO struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorsDTO struct {
	Errors []ErrorDTO `json:"errors"`
}

type ErrorResponse struct {
	Message string     `json:"message"`
	Errors  []ErrorDTO `json:"errors"`
}

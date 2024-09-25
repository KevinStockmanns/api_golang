package models

type ErrorWrapper struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

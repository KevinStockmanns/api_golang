package wrapper

type PageWrapper struct {
	Page          int         `json:"page"`
	Size          int         `json:"size"`
	TotalPage     int         `json:"totalPage"`
	TotalElements int64       `json:"totalElements"`
	Content       interface{} `json:"content"`
}

type ErrorWrapper struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

type ErrorsWrapper struct {
	Errors []ErrorWrapper `json:"errors"`
}

func (newError *ErrorsWrapper) AddNewError(field string, errorText string) {
	newError.Errors = append(newError.Errors, ErrorWrapper{Field: field, Error: errorText})
}

func (newError *ErrorsWrapper) AddEror(error ErrorWrapper) {
	newError.Errors = append(newError.Errors, error)
}

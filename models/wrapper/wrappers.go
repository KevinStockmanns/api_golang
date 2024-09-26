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

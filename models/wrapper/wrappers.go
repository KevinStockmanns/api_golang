package wrapper

type PageResponse struct {
	Page          int         `json:"page"`
	Size          int         `json:"size"`
	TotalPage     int         `json:"totalPage"`
	TotalElements int64       `json:"totalElements"`
	Content       interface{} `json:"content"`
}

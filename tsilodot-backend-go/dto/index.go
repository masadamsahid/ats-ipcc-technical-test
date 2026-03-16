package dto

type Pagination struct {
	CurrentPage int `json:"current_page"`
	TotalPages  int `json:"total_pages"`
	TotalItems  int `json:"total_items"`
}

type ResponseGeneric[T any] struct {
	Message    string      `json:"message"`
	Pagination *Pagination `json:"pagination,omitempty"`
	Data       *T          `json:"data,omitempty"`
	Errors     any         `json:"errors,omitempty"`
}

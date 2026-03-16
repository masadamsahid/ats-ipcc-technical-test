package dto

type ResponseGeneric[T any] struct {
	Message string `json:"message"`
	Data    *T     `json:"data,omitempty"`
	Errors  any    `json:"errors,omitempty"`
}

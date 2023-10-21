package api

import "time"

type Response[T any] struct {
	StatusCode int       `json:"status_code"`
	Data       T         `json:"data"`
	Date       time.Time `json:"date"`
}

func NewResponse[T any](statusCode int, data T) Response[T] {
	return Response[T]{
		StatusCode: statusCode,
		Data:       data,
		Date:       time.Now(),
	}
}

type ErrorResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

package api

import "time"

type Response struct {
	StatusCode int       `json:"status_code"`
	Data       any       `json:"data"`
	Date       time.Time `json:"date"`
}

func NewResponse(statusCode int, data any) Response {
	return Response{
		StatusCode: statusCode,
		Data:       data,
		Date:       time.Now(),
	}
}

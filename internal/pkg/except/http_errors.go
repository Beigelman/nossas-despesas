package except

import "net/http"

func NotFoundError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusNotFound, message...)
}

func BadRequestError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, message...)
}

func UnprocessableEntityError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusUnprocessableEntity, message...)
}

func InternalServerError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, message...)
}

func UnauthorizedError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusUnauthorized, message...)
}

func ForbiddenError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusForbidden, message...)
}

func ConflictError(message ...interface{}) *HTTPError {
	return NewHTTPError(http.StatusConflict, message...)
}

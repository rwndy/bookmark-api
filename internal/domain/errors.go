package domain

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func ErrNotFound(msg string) *AppError {
	return &AppError{Code: 404, Message: msg}
}

func ErrForbidden(msg string) *AppError {
	return &AppError{Code: 403, Message: msg}
}

func ErrBadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func ErrUnauthorized(msg string) *AppError {
	return &AppError{Code: 401, Message: msg}
}
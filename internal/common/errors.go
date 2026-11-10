package common

import (
	"errors"
	"net/http"
)

// AppError là lỗi nghiệp vụ với mã HTTP tương ứng.
type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error { return e.Err }

// BusinessRuleException lỗi nghiệp vụ (400).
func BusinessRuleException(message string) *AppError {
	return &AppError{Code: http.StatusBadRequest, Message: message}
}

// ConflictException lỗi xung đột (409).
func ConflictException(message string) *AppError {
	return &AppError{Code: http.StatusConflict, Message: message}
}

// ResourceNotFoundException lỗi không tìm thấy (404).
func ResourceNotFoundException(message string) *AppError {
	return &AppError{Code: http.StatusNotFound, Message: message}
}

// FromError chuyển error thành *AppError nếu có thể.
func FromError(err error) *AppError {
	if err == nil {
		return nil
	}
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return &AppError{Code: http.StatusInternalServerError, Message: err.Error(), Err: err}
}

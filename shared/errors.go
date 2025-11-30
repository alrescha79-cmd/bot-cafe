package shared

import "fmt"

// Error codes
const (
	ErrCodeInvalidInput   = "ERR_INVALID_INPUT"
	ErrCodeNotFound       = "ERR_NOT_FOUND"
	ErrCodeUnauthorized   = "ERR_UNAUTHORIZED"
	ErrCodeDatabaseError  = "ERR_DATABASE"
	ErrCodeInternalError  = "ERR_INTERNAL"
	ErrCodeDuplicateEntry = "ERR_DUPLICATE"
	ErrCodeServiceError   = "ERR_SERVICE"
)

// AppError represents application error
type AppError struct {
	Code    string
	Message string
	Err     error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// NewError creates a new application error
func NewError(code, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// NewInvalidInputError creates invalid input error
func NewInvalidInputError(message string) *AppError {
	return &AppError{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

// NewNotFoundError creates not found error
func NewNotFoundError(resource string) *AppError {
	return &AppError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s tidak ditemukan", resource),
	}
}

// NewUnauthorizedError creates unauthorized error
func NewUnauthorizedError() *AppError {
	return &AppError{
		Code:    ErrCodeUnauthorized,
		Message: "Akses tidak diizinkan",
	}
}

// NewDatabaseError creates database error
func NewDatabaseError(err error) *AppError {
	return &AppError{
		Code:    ErrCodeDatabaseError,
		Message: "Terjadi kesalahan database",
		Err:     err,
	}
}

// NewInternalError creates internal error
func NewInternalError(err error) *AppError {
	return &AppError{
		Code:    ErrCodeInternalError,
		Message: "Terjadi kesalahan internal",
		Err:     err,
	}
}

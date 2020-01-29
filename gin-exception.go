package ahlib_gin_gorm

import "gopkg.in/go-playground/validator.v9"

type ServerError struct {
	error
	Code    int
	Message string
}

func NewServerError(code int, message string) ServerError {
	return ServerError{
		Code:    code,
		Message: message,
	}
}

func IsValidationFormatError(err error) bool {
	// validator.ValidationErrors
	// *errors.errorString
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return false
	}

	for _, field := range errs {
		if field.Tag() == "required" {
			return false
		}
	}
	return true
}

package errors

type AppError struct {
	Code    int
	Message string
	Cause   error
	Details any
}

func (e *AppError) Error() string {
	return e.Message
}

func New(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

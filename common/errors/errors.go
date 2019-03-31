package errors

type AuthError interface {
	error
	Code() ErrorCodeType
}

type MyError struct {
	code    ErrorCodeType
	message string
}

func (e *MyError) Error() string {
	return e.message
}

func (e *MyError) Code() ErrorCodeType {
	return e.code
}

// New create a interface type error
func New(code ErrorCodeType, message string) AuthError {
	return &MyError{code: code, message: message}
}

package liberr

// Error is an implementation of error interface which has extra
// error information such as Code and Field
type Error struct {
	Message string
	Code    string
	Field   string
}

// NewError implements a new error instance
func NewError(msg, code, field string) *Error {
	return &Error{
		Message: msg,
		Code:    code,
		Field:   field,
	}
}

// Error returns error message
func (e *Error) Error() string {
	return e.Message
}

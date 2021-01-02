package errors

// TollErr ..
type TollErr struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// NewErrorWithCode ..
func NewErrorWithCode(code, message string) TollErr {
	return TollErr{
		Code:    code,
		Message: message,
	}
}

// Error ..
func (e TollErr) Error() string {
	return e.Code + " : " + e.Message
}

var (
	// ErrUnprocessableEntity ...
	ErrUnprocessableEntity = NewErrorWithCode(
		"ERR.HTTP.UNPROCESSABLEENTITY",
		"Invalid Json",
	)
)

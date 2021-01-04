package errors

// TollErr is a toll error type
type TollErr struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

// NewErrorWithCode returns a new instance of toll error type
func NewErrorWithCode(code, message string) TollErr {
	return TollErr{
		Code:    code,
		Message: message,
	}
}

// Error returns a string of error
func (e TollErr) Error() string {
	return e.Code + " : " + e.Message
}

// ToTollError takes error as argument and returns TollErr
func ToTollError(err error) *TollErr {
	if err == nil {
		return nil
	}
	appErr, ok := err.(TollErr)
	if ok {
		return &appErr
	}
	return &TollErr{
		Code:    "ERR",
		Message: err.Error(),
	}
}

var (
	// ErrUnprocessableEntity is error indicating the data is non processible
	ErrUnprocessableEntity = NewErrorWithCode(
		"ERR.HTTP.UNPROCESSABLEENTITY",
		"Invalid Json",
	)

	// ErrMissingFields is error indicating there is some missing fields to process the request
	ErrMissingFields = NewErrorWithCode(
		"ERR.APP.MISSING_FIELDS",
		"There are some missing fields that are required",
	)
)

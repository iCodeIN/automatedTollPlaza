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

// ToTollError takes error as argument and returns LeAPIError
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
	// ErrUnprocessableEntity ...
	ErrUnprocessableEntity = NewErrorWithCode(
		"ERR.HTTP.UNPROCESSABLEENTITY",
		"Invalid Json",
	)

	// ErrMissingFields ...
	ErrMissingFields = NewErrorWithCode(
		"ERR.APP.MISSING_FIELDS",
		"There are some missing fields that are required",
	)
)

package customerror

type CustomError interface {
	Error() string
	Code() int
}

type customError struct {
	Err        string
	StatusCode int
}

func (e *customError) Error() string {
	return e.Err
}

func (e *customError) Code() int {
	return e.StatusCode
}

func NewCustomError(err string, statusCode int) CustomError {
	return &customError{
		Err:        err,
		StatusCode: statusCode,
	}
}

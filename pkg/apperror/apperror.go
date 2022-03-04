package apperror

import "errors"

//AppError - custom error
type AppError struct {
	text   string
	status int
}

//ErrInternalServer - internal server error
var ErrInternalServer = errors.New("internal server error")

func (a AppError) Error() string {
	return a.text
}

//Type - returns status code of error
func (a AppError) Type() int {
	return a.status
}

//BadRequestError - returns bad request error
func BadRequestError(text string) error {
	return &AppError{
		text:   text,
		status: 400,
	}
}

//NotFoundError - returns not found error
func NotFoundError(text string) error {
	return &AppError{
		text:   text,
		status: 404,
	}
}

func Unauthorized() error {
	return &AppError{
		text:   "You are not authorized to access this service",
		status: 401,
	}
}

//IsAppError - returns apperror value and true if err is of type apperror
func IsAppError(err error) (*AppError, bool) {
	appError := new(AppError)
	if errors.As(err, &appError) {
		return appError, true
	}
	return nil, false
}

//IsNotFoundError -
func IsNotFoundError(err error) bool {
	a, ok := IsAppError(err)
	if !ok {
		return false
	}
	if a.Type() == 404 {
		return true
	}
	return false
}

//IsBadRequestError -
func IsBadRequestError(err error) bool {
	a, ok := IsAppError(err)
	if !ok {
		return false
	}
	if a.Type() == 400 {
		return true
	}
	return false
}

//UserFriendlyError - returns error of type apperror
func UserFriendlyError(text string, status int) error {
	return &AppError{
		text:   text,
		status: status,
	}
}

//CouldNotCompleteRequest -
func CouldNotCompleteRequest() error {
	return UserFriendlyError("Could not complete request", 500)
}

//IsAppErrorBool - returns apperror value and true if err is of type apperror
func IsAppErrorBool(err error) bool {
	appError := new(AppError)
	return errors.As(err, &appError)
}

//ErrorCode - returns error code if error is of type AppError else returns 0
func ErrorCode(err error) int {
	e, ok := IsAppError(err)
	if ok == false {
		return 0
	}
	return e.Type()
}

package errors

type (
	AppErrorInterface interface {
		Error() string
		Unwrap() error
		Message() string
		StatusCode() int
	}
)

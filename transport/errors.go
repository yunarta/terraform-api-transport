package transport

type BadRequestError struct {
	code  int
	error string
}

func (e BadRequestError) Error() string {
	return e.error
}

type MockRequestError struct {
	error string
	path  string
}

func (e MockRequestError) Error() string {
	return e.error
}

var _ error = BadRequestError{}
var _ error = MockRequestError{}

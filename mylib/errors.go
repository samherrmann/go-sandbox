package mylib

// UnauthorizedError ...
type UnauthorizedError struct{}

// Unauthorized ...
func (err *UnauthorizedError) Unauthorized() bool {
	return true
}

func (err *UnauthorizedError) Error() string {
	return "unauthorized"
}

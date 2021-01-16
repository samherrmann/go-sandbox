package mylib

// DoSomething performs some action.
func DoSomething() error {
	return new(UnauthorizedError)
}

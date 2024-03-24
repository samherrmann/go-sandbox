package view

type ViewData[T any] struct {
	Title string
	Path  string
	Error string
	Main  T
}

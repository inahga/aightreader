// Package ui provides an interface for using aightreader with various different
// interfaces.
package ui

type UI interface {
	Start() error

	// Main must be called last from the program main function. Gio requires this
	// for cross platform purposes.
	Main()
}

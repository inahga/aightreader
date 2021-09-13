package device

import "context"

type Device interface {
	// Listen asynchronously listens for key presses.
	Listen(context.Context) error

	// NoteOn sets a callback function for when a key is pressed.
	NoteOn(func(key uint8))

	// NoteOn sets a callback function for when a key is released.
	NoteOff(func(key uint8))

	Close() error
}

// Package ui provides an interface for using aightreader with various different
// interfaces.
package ui

import "github.com/inahga/aightreader/device"

type UI interface {
	Start() error

	// Main must be called last from the program main function. Gio requires this
	// for cross platform purposes.
	Main()

	// ChooseDevice presents devices and requests the user chooses one. It blocks
	// until a choice is made. If no devices are available, it will block indefinitely
	// until one is available. The UI should periodically recheck for new devices
	// and rerender itself as necessary.
	ChooseDevice() device.Device
	// ChooseDeviceError raises an error on the ChooseDevice screen. It does not
	// block.
	ChooseDeviceError(error)
	// DidChooseDevice sets a callback for when the user chooses a new device.
	DidChooseDevice(func(device.Device) error)
}

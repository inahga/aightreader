package device

type driver interface {
	ListDevices() ([]Device, error)
}

var drivers []driver

type Device interface {
	Description() string

	// Listen asynchronously listens for key presses.
	Listen() error

	// NoteOn sets a callback function for when a key is pressed.
	NoteOn(func(key uint8))

	// NoteOn sets a callback function for when a key is released.
	NoteOff(func(key uint8))

	// SetDebug enables or disables debug logging.
	SetDebug(bool)

	Close() error
}

func ListDevices() (devices []Device, err error) {
	for _, driver := range drivers {
		d, err := driver.ListDevices()
		if err != nil {
			return nil, err
		}
		devices = append(devices, d...)
	}
	return
}

package device

import (
	"fmt"
)

func init() {
	drivers = append(drivers, &virtualDriver{})
}

type virtualDriver struct {
}

func (m *virtualDriver) ListDevices() (devices []Device, err error) {
	return nil, nil
}

type virtualDevice struct {
	debug           bool
	noteOn, noteOff func(key uint8)
}

func (v *virtualDevice) Listen() error {
	panic(fmt.Errorf("not implemented"))
}

func (v *virtualDevice) Description() string {
	return "Virtual Keyboard"
}

func (v *virtualDevice) NoteOn(fn func(key uint8)) {
	v.noteOn = fn
}

func (v *virtualDevice) NoteOff(fn func(key uint8)) {
	v.noteOff = fn
}

func (v *virtualDevice) Close() error {
	return nil
}

func (v *virtualDevice) SetDebug(b bool) {
	v.debug = b
}

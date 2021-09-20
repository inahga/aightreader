package device

func init() {
	drivers = append(drivers, &dummyDriver{})
}

type dummyDriver struct {
}

func (m *dummyDriver) ListDevices() (devices []Device, err error) {
	return []Device{&dummyDevice{}}, nil
}

type dummyDevice struct {
	debug           bool
	noteOn, noteOff func(key uint8)
}

func (v *dummyDevice) Listen() error {
	return nil
}

func (v *dummyDevice) Description() string {
	return "Dummy Keyboard"
}

func (v *dummyDevice) NoteOn(fn func(key uint8)) {
	v.noteOn = fn
}

func (v *dummyDevice) NoteOff(fn func(key uint8)) {
	v.noteOff = fn
}

func (v *dummyDevice) Close() error {
	return nil
}

func (v *dummyDevice) SetDebug(b bool) {
	v.debug = b
}

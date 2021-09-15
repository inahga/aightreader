package device

import (
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

func init() {
	drv, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}
	drivers = append(drivers, &midiDriver{Driver: drv})
}

type midiDriver struct {
	midi.Driver
}

func (m *midiDriver) ListDevices() (devices []Device, err error) {
	ins, err := m.Driver.Ins()
	if err != nil {
		return
	}
	for _, in := range ins {
		devices = append(devices, &midiDevice{in: in})
	}
	return
}

type midiDevice struct {
	in              midi.In
	noteOn, noteOff func(key uint8)
	debug           bool
}

func (m *midiDevice) Listen() error {
	if err := m.in.Open(); err != nil {
		return err
	}
	opts := []func(r *reader.Reader){
		reader.NoteOn(func(_ *reader.Position, _, key, _ uint8) {
			m.noteOn(key)
		}),
		reader.NoteOff(func(_ *reader.Position, _, key, _ uint8) {
			m.noteOff(key)
		}),
	}
	if !m.debug {
		opts = append(opts, reader.NoLogger())
	}

	rd := reader.New(opts...)
	if err := rd.ListenTo(m.in); err != nil {
		return fmt.Errorf("failed to listen on MIDI device: %w", err)
	}
	return nil
}

func (m *midiDevice) Description() string {
	return m.in.String()
}

func (m *midiDevice) NoteOn(fn func(key uint8)) {
	m.noteOn = fn
}

func (m *midiDevice) NoteOff(fn func(key uint8)) {
	m.noteOff = fn
}

func (m *midiDevice) Close() error {
	return m.in.Close()
}

func (m *midiDevice) SetDebug(b bool) {
	m.debug = b
}

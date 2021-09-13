package device

import (
	"context"
	"fmt"

	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

type MIDIDevice struct {
	in              midi.In
	noteOn, noteOff func(key uint8)
	debug           bool
}

var midiDriver midi.Driver

func init() {
	drv, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}
	midiDriver = drv
}

func ListMIDIDevices() ([]midi.In, error) {
	return midiDriver.Ins()
}

func NewMIDIDevice(dev int, debug bool) (Device, error) {
	ins, err := midiDriver.Ins()
	if err != nil {
		return nil, err
	}
	in := ins[dev]

	if err := in.Open(); err != nil {
		return nil, err
	}

	return &MIDIDevice{in: in}, nil
}

func (m *MIDIDevice) Listen(_ context.Context) error {
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

func (m *MIDIDevice) NoteOn(fn func(key uint8)) {
	m.noteOn = fn
}

func (m *MIDIDevice) NoteOff(fn func(key uint8)) {
	m.noteOff = fn
}

func (m *MIDIDevice) Close() error {
	return m.in.Close()
}

package text

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/inahga/aightreader/device"
)

// TextUI communicates to the user strictly through text, asynchronously. It is not
// to be confused with a TUI, which draws graphics using terminal elements.
//
// This UI is meant for debugging.
type TextUI struct {
	out *os.File
	in  *bufio.Reader
}

func New() (*TextUI, error) {
	return &TextUI{
		out: os.Stdout,
		in:  bufio.NewReader(os.Stdin),
	}, nil
}

func (t *TextUI) Start() error { return nil }

func (t *TextUI) Main() {
	select {}
}

func (t *TextUI) ChooseDevice(devices []device.Device) device.Device {
	devices, err := device.ListDevices()
	if err != nil {
		panic(fmt.Errorf("bug: cannot list available devices: %w", err))
	} else if len(devices) == 0 {
		fmt.Fprintln(t.out, "no available devices!")
	}
	for index, device := range devices {
		fmt.Fprintf(t.out, "%d: %s\n", index, device.Description())
	}

	for {
		fmt.Fprint(t.out, "Choose a device: ")

		input, err := t.in.ReadString('\n')
		if err != nil {
			panic(fmt.Errorf("bug: failed to read from stdin: %w", err))
		}

		choice, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil {
			continue
		}
		if 0 > choice || choice > len(devices) {
			continue
		}
		return devices[choice]
	}
}

func (t *TextUI) ChooseDeviceError(err error) {
	fmt.Fprintln(t.out, err)
}

func (t *TextUI) DidChooseDevice(func(device.Device) error) {
}

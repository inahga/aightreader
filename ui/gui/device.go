package gui

import (
	"fmt"

	"github.com/inahga/aightreader/device"
)

func (t *GUI) ChooseDevice() device.Device {
	devices, err := device.ListDevices()
	if err != nil {
		panic(fmt.Errorf("bug: cannot list available devices: %w", err))
	} else if len(devices) == 0 {
		fmt.Println("no available devices!")
		select {}
	}
	return devices[0]
}

func (t *GUI) ChooseDeviceError(err error) {

}

func (t *GUI) DidChooseDevice(func(device.Device) error) {

}

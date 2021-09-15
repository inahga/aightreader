package gui

import "github.com/inahga/aightreader/device"

func (t *GUI) ChooseDevice(devices []device.Device) device.Device {
	return nil
}

func (t *GUI) ChooseDeviceError(err error) {

}

func (t *GUI) DidChooseDevice(func(device.Device) error) {

}

package main

import (
	"context"

	"github.com/inahga/aightreader/device"
	"github.com/inahga/aightreader/game"
)

func main() {
	const dev = 1

	// create the MIDI device
	midi, err := device.NewMIDIDevice(dev, true)
	if err != nil {
		panic(err)
	}

	// start the game
	if err := game.Start(context.Background(), midi); err != nil {
		panic(err)
	}
}

package main

import (
	"context"
	"flag"

	"github.com/inahga/aightreader/device"
	"github.com/inahga/aightreader/game"
	"github.com/inahga/aightreader/ui"
	"github.com/inahga/aightreader/ui/gui"
	"github.com/inahga/aightreader/ui/text"
)

func main() {
	textUI := flag.Bool("text", false, "run with text UI")
	flag.Parse()

	const dev = 1

	midi, err := device.NewMIDIDevice(dev, true)
	if err != nil {
		panic(err)
	}

	// Start game state
	go func() {
		if err := game.Start(context.Background(), midi); err != nil {
			panic(err)
		}
	}()

	var u ui.UI
	if *textUI {
		u, err = text.New()
	} else {
		u, err = gui.New()
	}
	if err != nil {
		panic(err)
	}
	u.Start()
	u.Main()
}

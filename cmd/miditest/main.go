// Command miditest lists available MIDI devices, or listens for messages on one
// of them.
package main

import (
	"fmt"
	"os"
	"strconv"

	"gitlab.com/gomidi/midi/reader"
	"gitlab.com/gomidi/rtmididrv"
)

func main() {
	drv, err := rtmididrv.New()
	if err != nil {
		panic(err)
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		dev, err := strconv.Atoi(os.Args[1])
		if err != nil {
			panic(err)
		}

		in := ins[dev]
		if err := in.Open(); err != nil {
			panic(err)
		}
		defer in.Close()

		rd := reader.New()
		if err := rd.ListenTo(in); err != nil {
			panic(err)
		}
		<-make(chan struct{})
	} else {
		for index, in := range ins {
			fmt.Printf("%d: %s\n", index, in.String())
		}
	}
}

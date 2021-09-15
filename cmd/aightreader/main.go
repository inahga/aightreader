package main

import (
	"context"
	"flag"

	"github.com/inahga/aightreader/game"
	"github.com/inahga/aightreader/ui"
	"github.com/inahga/aightreader/ui/gui"
	"github.com/inahga/aightreader/ui/text"
)

func main() {
	textUI := flag.Bool("text", false, "run with text UI")
	flag.Parse()

	var (
		u   ui.UI
		err error
	)
	if *textUI {
		u, err = text.New()
	} else {
		u, err = gui.New()
	}
	if err != nil {
		panic(err)
	}
	u.Start()

	go func() {
		if err := game.Start(context.Background(), u); err != nil {
			panic(err)
		}
	}()

	u.Main()
}

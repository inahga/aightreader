package gui

import "gioui.org/app"

type GUI struct{}

func New() (*GUI, error) {
	return &GUI{}, nil
}

func (g *GUI) Start() error {
	go func() {
		w := app.NewWindow(
			app.Title("aightreader"),
		)
		for range w.Events() {
		}
	}()
	return nil
}

func (g *GUI) Main() {
	app.Main()
}

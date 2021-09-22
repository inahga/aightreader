package gui

import (
	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
)

type (
	C = layout.Context
	D = layout.Dimensions

	GUI struct{}
)

func New() (*GUI, error) {
	return &GUI{}, nil
}

func (g *GUI) Start() error {
	go func() {
		w := app.NewWindow(
			app.Title("aightreader"),
		)
		if err := drawWindow(w); err != nil {
			panic(err)
		}
	}()
	return nil
}

func drawWindow(w *app.Window) error {
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			staff := GrandStaff{}
			staff.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func (g *GUI) Main() {
	app.Main()
}

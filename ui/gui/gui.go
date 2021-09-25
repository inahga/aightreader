package gui

import (
	"image"
	"image/draw"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"

	T "github.com/inahga/aightreader/theory"
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

	staff := GrandStaff{
		StaffLineWeight: 1,
		StaffLines:      31,
		TopStaffLine:    9,
		BottomStaffLine: 24,
		TimeSignature:   T.TimeSignature{3, 4},
		Key:             T.MustGetKeyByClass(T.G, T.Natural, T.Minor),
	}

	for e := range w.Events() {
		switch e := e.(type) {
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			staff.Layout(gtx)
			e.Frame(gtx.Ops)
		}
	}
	return nil
}

func (g *GUI) Main() {
	app.Main()
}

func canvas(gtx C) draw.Image {
	c := gtx.Constraints
	return image.NewRGBA(image.Rect(0, 0, c.Max.X, c.Max.Y))
}

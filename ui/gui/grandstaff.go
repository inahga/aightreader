package gui

import (
	"image"

	"gioui.org/op/paint"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

// GrandStaff is a Gio widget that renders the grand staff.
type GrandStaff struct{}

func (g *GrandStaff) Layout(gtx C) D {
	face, err := opentype.NewFace(Leland, &opentype.FaceOptions{
		Size: 200,
		DPI:  60,
	})
	if err != nil {
		panic(err)
	}

	rect := image.Rect(0, 0, 1000, 500)
	dst := image.NewRGBA(rect)
	d := font.Drawer{
		Dst:  dst,
		Src:  image.Black,
		Face: face,
		Dot:  fixed.P(6, 300),
	}
	d.DrawString(string([]rune{'\uE050'}))

	imageOp := paint.NewImageOp(dst)
	imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: image.Point{X: rect.Dx(), Y: rect.Dy()}}
}

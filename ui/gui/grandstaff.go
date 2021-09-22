package gui

import (
	"image"
	"image/draw"

	"gioui.org/op/paint"
)

// GrandStaff is a Gio widget that renders the grand staff.
type GrandStaff struct{}

func (g *GrandStaff) Layout(gtx C) D {
	store, err := newGlyphStore(glyphFonts["leland"], float64(gtx.Constraints.Max.Y), 96)
	if err != nil {
		panic(err)
	}

	rect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	dst := image.NewGray(rect)
	mask := store.MustGetGlyphMask("trebleClef")
	draw.DrawMask(dst, rect, image.White, image.Point{}, mask.mask, mask.maskp, draw.Over)

	imageOp := paint.NewImageOp(dst)
	imageOp.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return D{Size: image.Point{X: rect.Dx(), Y: rect.Dy()}}
}

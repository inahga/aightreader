package gui

import (
	"image"
	"image/draw"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
)

const (
	defaultDPI = 100 // IDK what the correct value should be...
)

// GrandStaff is a Gio widget that renders the grand staff.
type GrandStaff struct {
	StaffLineWeight int // Thickness in pixels of staff lines
	StaffLines      int // Number of staff and ledger lines to make space for.
	TopStaffLine    int // Number of ledger lines before the top staff starts.
	BottomStaffLine int // Number of ledger lines where the bottom staff ends.

	staffLineHeight int
	glyphStore      *glyphStore
	leftOffset      int
}

func (g *GrandStaff) Layout(gtx C) D {
	g.staffLineHeight = gtx.Constraints.Max.Y / g.StaffLines

	// placeholder offset
	g.leftOffset = 5

	store, err := newGlyphStore(glyphFonts["leland"], 2.7*float64(g.staffLineHeight), defaultDPI)
	if err != nil {
		panic(err)
	}
	g.glyphStore = store

	g.drawStaffLines(gtx)
	g.drawClefs(gtx)
	return D{Size: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
}

func (g *GrandStaff) drawStaffLines(gtx C) {
	height, weight := g.staffLineHeight, g.StaffLineWeight
	for _, lines := range []struct{ a, b int }{
		{g.TopStaffLine, g.TopStaffLine + 5},
		{g.BottomStaffLine - 5, g.BottomStaffLine},
	} {
		for i := lines.a; i < lines.b; i++ {
			save := op.Save(gtx.Ops)
			clip.Rect{
				Min: image.Pt(g.leftOffset, height*i),
				Max: image.Pt(gtx.Constraints.Max.X, (height*i)+weight),
			}.Add(gtx.Ops)
			paint.PaintOp{}.Add(gtx.Ops)
			save.Load()
		}
	}

	g.drawVerticalStaffLine(gtx, g.leftOffset, g.StaffLineWeight)
	g.drawVerticalStaffLine(gtx, gtx.Constraints.Max.X-3*g.StaffLineWeight, 3*g.StaffLineWeight)
}

func (g *GrandStaff) drawVerticalStaffLine(gtx C, offset, width int) {
	defer op.Save(gtx.Ops).Load()

	topOffset := g.TopStaffLine * g.staffLineHeight
	bottomOffset := (g.BottomStaffLine - 1) * g.staffLineHeight

	clip.Rect{
		Min: image.Pt(offset, topOffset),
		Max: image.Pt(offset+width, bottomOffset),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func (g *GrandStaff) drawClefs(gtx C) {
	// placeholder offsets
	g.drawGlyph(gtx, "trebleClef", image.Pt(15, (g.TopStaffLine-1)*g.staffLineHeight))
	g.drawGlyph(gtx, "bassClef", image.Pt(15, (g.BottomStaffLine-5)*g.staffLineHeight))
}

func (g *GrandStaff) drawGlyph(gtx C, name string, point image.Point) {
	defer op.Save(gtx.Ops).Load()

	mask := g.glyphStore.MustGetGlyphMask(name)
	rect := image.Rect(0, 0, gtx.Constraints.Max.X, gtx.Constraints.Max.Y)
	dst := image.NewRGBA(rect)
	draw.DrawMask(dst, rect.Add(point), image.Black, image.Point{}, mask.mask, mask.maskp, draw.Over)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

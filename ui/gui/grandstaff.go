package gui

import (
	"fmt"
	"image"
	"image/draw"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	xdraw "golang.org/x/image/draw"
)

const (
	defaultDPI   = 100 // IDK what the correct value should be...
	bracePadding = 2
)

// GrandStaff is a Gio widget that renders the grand staff.
type GrandStaff struct {
	StaffLineWeight int // Thickness in pixels of staff lines
	StaffLines      int // Number of staff and ledger lines to make space for.
	TopStaffLine    int // Number of ledger lines before the top staff starts.
	BottomStaffLine int // Number of ledger lines where the bottom staff ends.

	GlyphStore      *glyphStore
	LargeGlyphStore *glyphStore
	leftOffset      int
}

func (g *GrandStaff) Layout(gtx C) D {
	g.GlyphStore = mustNewGlyphStore(glyphFonts["leland"],
		lelandMagicCoefficient*float64(g.staffLineHeight(gtx)), defaultDPI)
	g.LargeGlyphStore = mustNewGlyphStore(glyphFonts["leland"], 1000, defaultDPI)

	g.leftOffset = g.drawLeftBrace(gtx)
	g.drawStaffLines(gtx)
	g.drawClefs(gtx)
	return D{Size: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
}

func (g *GrandStaff) drawStaffLines(gtx C) {
	height, weight := g.staffLineHeight(gtx), g.StaffLineWeight
	gap := g.numGapLedgerLines()
	for _, lines := range []struct{ a, b int }{
		{g.TopStaffLine, g.TopStaffLine + gap},
		{g.BottomStaffLine - gap, g.BottomStaffLine},
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

	// Align with left.
	g.drawVerticalStaffLine(gtx, g.leftOffset, g.StaffLineWeight)
	// Align to right and thicken the line.
	g.drawVerticalStaffLine(gtx, gtx.Constraints.Max.X-3*g.StaffLineWeight, 3*g.StaffLineWeight)
}

func (g *GrandStaff) drawVerticalStaffLine(gtx C, offset, width int) {
	defer op.Save(gtx.Ops).Load()

	topOffset := g.TopStaffLine * g.staffLineHeight(gtx)
	bottomOffset := (g.BottomStaffLine - 1) * g.staffLineHeight(gtx)

	clip.Rect{
		Min: image.Pt(offset, topOffset),
		Max: image.Pt(offset+width, bottomOffset),
	}.Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
}

func (g *GrandStaff) drawClefs(gtx C) {
	// placeholder offsets
	// Treble clef baseline should be aligned with G.
	a := g.drawGlyph(gtx, "trebleClef", image.Pt(g.leftOffset+10,
		(g.TopStaffLine+3)*g.staffLineHeight(gtx)+(g.StaffLineWeight/2)))
	// Bass clef baseline should be aligned with F.
	b := g.drawGlyph(gtx, "bassClef", image.Pt(g.leftOffset+10,
		(g.BottomStaffLine-4)*g.staffLineHeight(gtx)+(g.StaffLineWeight/2)))
	fmt.Println(a, b)
}

func (g *GrandStaff) drawGlyph(gtx C, name string, point image.Point) int {
	defer op.Save(gtx.Ops).Load()

	mask := g.GlyphStore.MustGetGlyphMask(name)
	dst := canvas(gtx)
	draw.DrawMask(dst, mask.dr.Add(point), image.Black, image.Point{}, mask.mask, mask.maskp, draw.Over)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return mask.dr.Dx()
}

func (g *GrandStaff) drawLeftBrace(gtx C) int {
	defer op.Save(gtx.Ops).Load()
	topOffset, bottomOffset := g.topOffset(gtx), g.bottomOffset(gtx)
	dst := canvas(gtx)

	// The font by default doesn't come with a large enough brace. Start with it
	// very large, then scale it down.
	mask := g.LargeGlyphStore.MustGetGlyphMask("leftBrace")
	unmasked := image.NewRGBA(mask.dr)
	draw.DrawMask(unmasked, mask.dr, image.Black, image.Point{}, mask.mask, mask.maskp, draw.Over)

	dimensions := mask.dr.Max.Sub(mask.dr.Min)
	scaledWidth := int(float64(dimensions.X) * float64(bottomOffset-topOffset) / float64(dimensions.Y))
	dr := image.Rect(0, topOffset, scaledWidth, bottomOffset)

	xdraw.NearestNeighbor.Scale(dst, dr, unmasked, mask.dr, draw.Over, nil)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return scaledWidth + bracePadding
}

// numGapLedgerLines is the number of ledger lines allowable between the top
// and bottom staffs.
func (g *GrandStaff) numGapLedgerLines() int {
	return g.BottomStaffLine - g.TopStaffLine - 10
}

func (g *GrandStaff) staffLineHeight(gtx C) int {
	return gtx.Constraints.Max.Y / g.StaffLines
}

// topOffset returns the number of pixels until the top staff line is reached.
func (g *GrandStaff) topOffset(gtx C) int {
	return g.TopStaffLine * g.staffLineHeight(gtx)
}

// bottomOffset returns the number of pixels until the bottom staff line is reached.
func (g *GrandStaff) bottomOffset(gtx C) int {
	return (g.BottomStaffLine-1)*g.staffLineHeight(gtx) + g.StaffLineWeight
}

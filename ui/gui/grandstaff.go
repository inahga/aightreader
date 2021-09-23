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
	defaultDPI   = 300
	bracePadding = 2
)

type (

	// GrandStaff is a Gio widget that renders the grand staff.
	GrandStaff struct {
		StaffLineWeight int // Thickness in pixels of staff lines
		StaffLines      int // Number of staff and ledger lines to make space for.
		TopStaffLine    int // Number of ledger lines before the top staff starts.
		BottomStaffLine int // Number of ledger lines where the bottom staff ends.
		TimeSignature

		glyphStore      *glyphStore
		largeGlyphStore *glyphStore
		leftOffset      int
	}

	TimeSignature struct {
		BeatsPerBar, BeatUnit int
	}
)

func (g *GrandStaff) Layout(gtx C) D {
	// Problem: these stores are recreated on every frame redraw!
	g.glyphStore = mustNewGlyphStore(glyphFonts["leland"],
		lelandMagicCoefficient*float64(g.staffLineHeight(gtx)), defaultDPI)
	g.largeGlyphStore = mustNewGlyphStore(glyphFonts["leland"], 1000, defaultDPI)

	g.leftOffset = g.drawLeftBrace(gtx)
	g.drawStaffLines(gtx)
	g.leftOffset = g.drawClefs(gtx)
	g.leftOffset = g.drawTimeSignature(gtx)
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

// drawClefs draws the clefs at the current left offset, then returns the new
// left offset.
func (g *GrandStaff) drawClefs(gtx C) int {
	offset := g.leftOffset + g.leftOffset/2

	// Treble clef baseline should be aligned with G.
	t := g.drawGlyph(gtx, "trebleClef", image.Pt(offset,
		(g.TopStaffLine+3)*g.staffLineHeight(gtx)+(g.StaffLineWeight/2)))

	// Bass clef baseline should be aligned with F.
	b := g.drawGlyph(gtx, "bassClef", image.Pt(offset,
		(g.BottomStaffLine-4)*g.staffLineHeight(gtx)+(g.StaffLineWeight/2)))

	if t > b {
		return t
	}
	return b + offset
}

func (g *GrandStaff) drawTimeSignature(gtx C) int {
	offset := g.leftOffset + g.leftOffset/5

	num := fmt.Sprintf("timeSignature%d", g.TimeSignature.BeatsPerBar)
	denom := fmt.Sprintf("timeSignature%d", g.TimeSignature.BeatUnit)

	g.drawGlyph(gtx, num, image.Pt(offset, g.yOffset(gtx, g.TopStaffLine+1)))
	g.drawGlyph(gtx, denom, image.Pt(offset, g.yOffset(gtx, g.TopStaffLine+3)))
	g.drawGlyph(gtx, num, image.Pt(offset, g.yOffset(gtx, g.BottomStaffLine-4)))
	g.drawGlyph(gtx, denom, image.Pt(offset, g.yOffset(gtx, g.BottomStaffLine-2)))
	return 0
}

// drawGlyph returns the width of the drawn glyph.
func (g *GrandStaff) drawGlyph(gtx C, name string, point image.Point) int {
	defer op.Save(gtx.Ops).Load()

	mask := g.glyphStore.MustGetGlyphMask(name)
	dst := canvas(gtx)
	draw.DrawMask(dst, mask.dr.Add(point), image.Black, image.Point{}, mask.mask, mask.maskp, draw.Over)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)
	return mask.dr.Dx()
}

// drawLeftBrace returns the width of the drawn brace.
func (g *GrandStaff) drawLeftBrace(gtx C) int {
	defer op.Save(gtx.Ops).Load()
	topOffset, bottomOffset := g.topOffset(gtx), g.bottomOffset(gtx)
	dst := canvas(gtx)

	// The font by default doesn't come with a large enough brace. Start with it
	// very large, then scale it down.
	mask := g.largeGlyphStore.MustGetGlyphMask("leftBrace")
	unmasked := image.NewRGBA(mask.dr)
	draw.DrawMask(unmasked, mask.dr, image.Black, image.Point{}, mask.mask, mask.maskp, draw.Over)

	dimensions := mask.dr.Max.Sub(mask.dr.Min)
	scaledWidth := int(float64(dimensions.X) * float64(bottomOffset-topOffset) / float64(dimensions.Y))
	dr := image.Rect(0, topOffset, scaledWidth, bottomOffset)

	// Bilinear is slow... but everything else looks terrible on high DPI.
	xdraw.BiLinear.Scale(dst, dr, unmasked, mask.dr, draw.Over, nil)
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

func (g *GrandStaff) yOffset(gtx C, line int) int {
	return line*g.staffLineHeight(gtx) + (g.StaffLineWeight / 2)
}

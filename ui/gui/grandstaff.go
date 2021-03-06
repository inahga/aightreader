package gui

import (
	"fmt"
	"image"
	"image/draw"

	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	xdraw "golang.org/x/image/draw"

	T "github.com/inahga/aightreader/theory"
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
		TopStaffLine    int // Number of ledger lines before the top staff starts (i.e. F).
		BottomStaffLine int // Number of staff and ledger lines where the bottom staff ends (i.e. G).

		T.Time
		T.Key

		glyphStore      *glyphStore
		largeGlyphStore *glyphStore
		leftOffset      int
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
	g.leftOffset = g.drawKeySignature(gtx)
	g.leftOffset = g.drawTimeSignature(gtx)
	return D{Size: image.Point{X: gtx.Constraints.Max.X, Y: gtx.Constraints.Max.Y}}
}

func (g *GrandStaff) drawStaffLines(gtx C) {
	height, weight := g.staffLineHeight(gtx), g.StaffLineWeight
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
	t := g.drawGlyph(gtx, "trebleClef", image.Pt(offset, g.yOffset(gtx, g.TopStaffLine+3)))
	// Bass clef baseline should be aligned with F.
	b := g.drawGlyph(gtx, "bassClef", image.Pt(offset, g.yOffset(gtx, g.BottomStaffLine-4)))

	if t > b {
		return t
	}
	return b + offset
}

func (g *GrandStaff) drawKeySignature(gtx C) int {
	offset := g.leftOffset
	sig := g.Key.Sig

	var glyphName string
	if len(sig) > 0 {
		if sig[0].Accidental == T.Sharp {
			glyphName = "sharpAccidental"
		} else if sig[0].Accidental == T.Flat {
			glyphName = "flatAccidental"
		}
	}

	// Arbitrary padding, adjust as needed.
	offset += g.leftOffset / 10
	for _, sig := range g.Key.Sig {
		// Arbitrary padding, adjust as needed.
		offset += g.leftOffset / 15
		g.drawGlyph(gtx, glyphName, image.Pt(offset, g.bassNoteYOffset(gtx, sig.BassNote())))
		offset += g.drawGlyph(gtx, glyphName, image.Pt(offset, g.trebleNoteYOffset(gtx, sig.TrebleNote())))
	}
	return offset
}

func (g *GrandStaff) drawTimeSignature(gtx C) int {
	// Arbitrary padding, adjust as needed.
	offset := g.leftOffset + g.leftOffset/10

	num := fmt.Sprintf("timeSignature%d", g.Time.BeatsPerBar)
	denom := fmt.Sprintf("timeSignature%d", g.Time.BeatUnit)

	// Center the smaller glyph relative to the bigger one.
	var numOffset, denomOffset, ret int
	numWidth := g.glyphStore.MustGetGlyphMask(num).Dimensions().X
	denomWidth := g.glyphStore.MustGetGlyphMask(denom).Dimensions().X

	if numWidth > denomWidth {
		denomOffset = (numWidth - denomWidth) / 2
		ret = numWidth
	} else {
		numOffset = (denomWidth - numWidth) / 2
		ret = denomWidth
	}

	g.drawGlyph(gtx, num, image.Pt(offset+numOffset, g.yOffset(gtx, g.TopStaffLine+1)))
	g.drawGlyph(gtx, denom, image.Pt(offset+denomOffset, g.yOffset(gtx, g.TopStaffLine+3)))
	g.drawGlyph(gtx, num, image.Pt(offset+numOffset, g.yOffset(gtx, g.BottomStaffLine-4)))
	g.drawGlyph(gtx, denom, image.Pt(offset+denomOffset, g.yOffset(gtx, g.BottomStaffLine-2)))
	return offset + ret
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

	dimensions := mask.Dimensions()
	scaledWidth := int(float64(dimensions.X) * float64(bottomOffset-topOffset) / float64(dimensions.Y))
	dr := image.Rect(0, topOffset, scaledWidth, bottomOffset)

	// Bilinear is slow... but everything else looks terrible on high DPI.
	xdraw.BiLinear.Scale(dst, dr, unmasked, mask.dr, draw.Over, nil)
	paint.NewImageOp(dst).Add(gtx.Ops)
	paint.PaintOp{}.Add(gtx.Ops)

	return scaledWidth + bracePadding
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

func (g *GrandStaff) bassNoteYOffset(gtx C, note T.Note) int {
	// c0 is 9 ledger lines below bass clef staff
	return g.noteYOffset(gtx, note, g.bottomOffset(gtx)+9*g.staffLineHeight(gtx))
}

func (g *GrandStaff) trebleNoteYOffset(gtx C, note T.Note) int {
	// c0 is 19 lines and ledger lines below top line of treble clef staff
	return g.noteYOffset(gtx, note, g.topOffset(gtx)+19*g.staffLineHeight(gtx))
}

func (g *GrandStaff) noteYOffset(gtx C, note T.Note, relative int) int {
	dist := int(note.Octave)*7 + int(note.Class)
	return relative - (dist * g.staffLineHeight(gtx) / 2)
}

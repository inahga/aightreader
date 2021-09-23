package gui

import (
	"fmt"
	"image"
	"sync"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/font/sfnt"
	"golang.org/x/image/math/fixed"
)

type (
	glyphFont struct {
		Name string
		*sfnt.Font
		SymbolMap map[string]rune
	}

	glyphMask struct {
		dr      image.Rectangle
		mask    image.Image
		maskp   image.Point
		advance fixed.Int26_6
	}

	glyphStore struct {
		Font   *glyphFont
		Face   font.Face
		glyphs map[string]*glyphMask
		lock   *sync.Mutex
	}
)

var glyphFonts = map[string]*glyphFont{}

func newGlyphStore(font *glyphFont, size float64, dpi float64) (*glyphStore, error) {
	face, err := opentype.NewFace(font.Font, &opentype.FaceOptions{
		Size: size,
		DPI:  dpi,
	})
	if err != nil {
		return nil, err
	}
	return &glyphStore{
		Font:   font,
		Face:   face,
		glyphs: map[string]*glyphMask{},
		lock:   &sync.Mutex{},
	}, nil
}

func mustNewGlyphStore(font *glyphFont, size float64, dpi float64) *glyphStore {
	store, err := newGlyphStore(font, size, dpi)
	if err != nil {
		return nil
	}
	return store
}

func (g *glyphStore) GetGlyphMask(name string) (*glyphMask, error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if glyph, ok := g.glyphs[name]; ok {
		return glyph, nil
	}

	r, ok := g.Font.SymbolMap[name]
	if !ok {
		return nil, fmt.Errorf("unregistered symbol %s", name)
	}

	dr, mask, maskp, advance, ok := g.Face.Glyph(fixed.P(0, 0), r)
	if !ok {
		return nil, fmt.Errorf("unknown rune for symbol %s", name)
	}

	glyph := &glyphMask{dr: dr, mask: mask, maskp: maskp, advance: advance}
	g.glyphs[name] = glyph
	return glyph, nil
}

func (g *glyphStore) MustGetGlyphMask(name string) *glyphMask {
	glyph, err := g.GetGlyphMask(name)
	if err != nil {
		panic(err)
	}
	return glyph
}

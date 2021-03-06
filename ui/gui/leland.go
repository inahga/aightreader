package gui

import (
	_ "embed"

	"golang.org/x/image/font/opentype"
)

var (
	//go:embed leland.otf
	lelandBytes []byte

	// Leland is the open source MuseScore font, distributed under the OFL-1.1
	// license. See https://github.com/MuseScoreFonts/Leland
	Leland *opentype.Font

	// lelandMagicCoefficient is the a coefficient that when multiplied by the
	// width between staff lines, results in a size value that fits the whole
	// note perfectly between the staff lines.
	lelandMagicCoefficient float64 = .95

	lelandSymbolMap = map[string]rune{
		"leftBrace":  '\uE000',
		"trebleClef": '\uE050',
		"bassClef":   '\uE062',

		"timeSignature0": '\uE080',
		"timeSignature1": '\uE081',
		"timeSignature2": '\uE082',
		"timeSignature3": '\uE083',
		"timeSignature4": '\uE084',
		"timeSignature5": '\uE085',
		"timeSignature6": '\uE086',
		"timeSignature7": '\uE087',
		"timeSignature8": '\uE088',
		"timeSignature9": '\uE089',

		"wholeNote":             '\uE1D2',
		"upwardHalfNote":        '\uE1D3',
		"downwardHalfNote":      '\uE1D4',
		"upwardQuarterNote":     '\uE1D5',
		"downwardQuarterNote":   '\uE1D6',
		"upwardEighthNote":      '\uE1D7',
		"downwardEighthNote":    '\uE1D8',
		"upwardSixteenthNote":   '\uE1D9',
		"downwardSixteenthNote": '\uE1DA',

		"climbingEighthTail":    '\uE220',
		"climbingSixteenthTail": '\uE221',

		"flatAccidental":    '\uE260',
		"naturalAccidental": '\uE261',
		"sharpAccidental":   '\uE262',

		"wholeRest":     '\uE4E3',
		"halfRest":      '\uE4E4',
		"quarterRest":   '\uE4E5',
		"eigthRest":     '\uE4E6',
		"sixteenthRest": '\uE4E7',

		"ottava":      '\uE512',
		"ottavaBassa": '\uE513',

		"piano":      '\uE520',
		"forte":      '\uE522',
		"pianissimo": '\uE52B',
		"mezzoPiano": '\uE52C',
		"mezzoForte": '\uE52D',
		"fortissimo": '\uE52F',

		"trill": '\uE566',

		"pedal":        '\uE650',
		"releasePedal": '\uE655',

		"finger0": '\uE880',
		"finger1": '\uE881',
		"finger2": '\uE882',
		"finger3": '\uE883',
		"finger4": '\uE884',
		"finger5": '\uE885',
		"finger6": '\uE886',
		"finger7": '\uE887',
		"finger8": '\uE888',
		"finger9": '\uE889',

		"fallingEighthTail":    '\uF40F',
		"fallingSixteenthTail": '\uF412',
	}
)

func init() {
	f, err := opentype.Parse(lelandBytes)
	if err != nil {
		panic(err)
	}
	glyphFont := &glyphFont{
		Name:      "leland",
		Font:      f,
		SymbolMap: lelandSymbolMap,
	}
	glyphFonts[glyphFont.Name] = glyphFont
}

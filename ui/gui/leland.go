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

	lelandSymbolMap = map[string]rune{
		"leftBrace":  '\uE000',
		"trebleClef": '\uE050',
		"bassClef":   '\uE062',

		"timeSignatureZero":  '\uE080',
		"timeSignatureOne":   '\uE081',
		"timeSignatureTwo":   '\uE082',
		"timeSignatureThree": '\uE083',
		"timeSignatureFour":  '\uE084',
		"timeSignatureFive":  '\uE085',
		"timeSignatureSix":   '\uE086',
		"timeSignatureSeven": '\uE087',
		"timeSignatureEight": '\uE088',
		"timeSignatureNine":  '\uE089',

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

		"fingerZero":  '\uE880',
		"fingerOne":   '\uE881',
		"fingerTwo":   '\uE882',
		"fingerThree": '\uE883',
		"fingerFour":  '\uE884',
		"fingerFive":  '\uE885',
		"fingerSix":   '\uE886',
		"fingerSeven": '\uE887',
		"fingerEight": '\uE888',
		"fingerNine":  '\uE889',

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

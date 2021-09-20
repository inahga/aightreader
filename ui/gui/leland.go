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
)

const (
	lelandLeftBrace  = '\uE000'
	lelandTrebleClef = '\uE050'
	lelandBassClef   = '\uE062'

	lelandTimeSignatureZero  = '\uE080'
	lelandTimeSignatureOne   = '\uE081'
	lelandTimeSignatureTwo   = '\uE082'
	lelandTimeSignatureThree = '\uE083'
	lelandTimeSignatureFour  = '\uE084'
	lelandTimeSignatureFive  = '\uE085'
	lelandTimeSignatureSix   = '\uE086'
	lelandTimeSignatureSeven = '\uE087'
	lelandTimeSignatureEight = '\uE088'
	lelandTimeSignatureNine  = '\uE089'

	lelandWholeNote             = '\uE1D2'
	lelandUpwardHalfNote        = '\uE1D3'
	lelandDownwardHalfNote      = '\uE1D4'
	lelandUpwardQuarterNote     = '\uE1D5'
	lelandDownwardQuarterNote   = '\uE1D6'
	lelandUpwardEighthNote      = '\uE1D7'
	lelandDownwardEighthNote    = '\uE1D8'
	lelandUpwardSixteenthNote   = '\uE1D9'
	lelandDownwardSixteenthNote = '\uE1DA'

	lelandClimbingEighthTail    = '\uE220'
	lelandClimbingSixteenthTail = '\uE221'

	lelandFlatAccidental    = '\uE260'
	lelandNaturalAccidental = '\uE261'
	lelandSharpAccidental   = '\uE262'

	lelandWholeRest     = '\uE4E3'
	lelandHalfRest      = '\uE4E4'
	lelandQuarterRest   = '\uE4E5'
	lelandEigthRest     = '\uE4E6'
	lelandSixteenthRest = '\uE4E7'

	lelandOttava      = '\uE512'
	lelandOttavaBassa = '\uE513'

	lelandPiano      = '\uE520'
	lelandForte      = '\uE522'
	lelandPianissimo = '\uE52B'
	lelandMezzoPiano = '\uE52C'
	lelandMezzoForte = '\uE52D'
	lelandFortissimo = '\uE52F'

	lelandTrill = '\uE566'

	lelandPedal        = '\uE650'
	lelandReleasePedal = '\uE655'

	lelandFingerZero  = '\uE880'
	lelandFingerOne   = '\uE881'
	lelandFingerTwo   = '\uE882'
	lelandFingerThree = '\uE883'
	lelandFingerFour  = '\uE884'
	lelandFingerFive  = '\uE885'
	lelandFingerSix   = '\uE886'
	lelandFingerSeven = '\uE887'
	lelandFingerEight = '\uE888'
	lelandFingerNine  = '\uE889'

	lelandFallingEighthTail    = '\uF40F'
	lelandFallingSixteenthTail = '\uF412'
)

func init() {
	f, err := opentype.Parse(lelandBytes)
	if err != nil {
		panic(err)
	}
	Leland = f
}

package theory

import (
	"fmt"
)

var (
	sharpOrder = []Signature{
		{Tone{F, Sharp}, 5, 3}, {Tone{C, Sharp}, 5, 3}, {Tone{G, Sharp}, 5, 3},
		{Tone{D, Sharp}, 5, 3}, {Tone{A, Sharp}, 4, 2}, {Tone{E, Sharp}, 5, 3},
		{Tone{B, Sharp}, 4, 2},
	}
	flatOrder = []Signature{
		{Tone{B, Flat}, 4, 2}, {Tone{E, Flat}, 5, 3}, {Tone{A, Flat}, 4, 2},
		{Tone{D, Flat}, 5, 3}, {Tone{G, Flat}, 4, 2}, {Tone{C, Flat}, 5, 3},
		{Tone{F, Flat}, 4, 2},
	}

	keys = []Key{
		{Tone{C, Natural}, Tone{A, Natural}, []Signature{}},

		// Clockwise circle of fifths.
		{Tone{G, Natural}, Tone{E, Natural}, sharpOrder[0:1]},
		{Tone{D, Natural}, Tone{B, Natural}, sharpOrder[0:2]},
		{Tone{A, Natural}, Tone{F, Sharp}, sharpOrder[0:3]},
		{Tone{E, Natural}, Tone{C, Sharp}, sharpOrder[0:4]},
		{Tone{B, Natural}, Tone{G, Sharp}, sharpOrder[0:5]},
		{Tone{F, Sharp}, Tone{D, Sharp}, sharpOrder[0:6]},
		{Tone{C, Sharp}, Tone{A, Sharp}, sharpOrder[0:7]},

		// Counterclockwise circle of fifths.
		{Tone{C, Flat}, Tone{A, Flat}, flatOrder[0:7]},
		{Tone{G, Flat}, Tone{E, Flat}, flatOrder[0:6]},
		{Tone{D, Flat}, Tone{B, Flat}, flatOrder[0:5]},
		{Tone{A, Flat}, Tone{F, Natural}, flatOrder[0:4]},
		{Tone{E, Flat}, Tone{C, Natural}, flatOrder[0:3]},
		{Tone{B, Flat}, Tone{G, Natural}, flatOrder[0:2]},
		{Tone{F, Natural}, Tone{D, Natural}, flatOrder[0:1]},
	}
)

func GetKeyByClass(class Class, adjustment Accidental, mode Mode) (Key, error) {
	for _, key := range keys {
		tone := Tone{class, adjustment}
		if (mode == Major && key.Major == tone) || (mode == Minor && key.Minor == tone) {
			return key, nil
		}
	}
	return Key{}, fmt.Errorf("no such key %s %s %s", class, adjustment, mode)
}

func MustGetKeyByClass(class Class, adjustment Accidental, mode Mode) Key {
	k, err := GetKeyByClass(class, adjustment, mode)
	if err != nil {
		panic(err)
	}
	return k
}

// BassNote returns a note indicating the octave at which the signature should
// be placed when rendering bass clef.
func (s Signature) BassNote() Note {
	return N(s.Class, s.Accidental, s.BassOctave, 0)
}

// BassNote returns a note indicating the octave at which the signature should
// be placed when rendering treble clef.
func (s Signature) TrebleNote() Note {
	return N(s.Class, s.Accidental, s.TrebleOctave, 0)
}

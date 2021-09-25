// Package theory implements music theory. There is also support for associating
// MIDI numbers with music theory.
//
// Note that this is inahga's(tm) music theory. My understanding of actual music
// theory is pretty bad :)
package theory

type (
	Octave     int
	Duration   float64
	Accidental int
	Mode       int
	Class      string

	Tone struct {
		Class
		Accidental
	}

	// Note represents the pitch and duration of a sound.
	Note struct {
		Tone
		Octave
		Duration
	}

	Signature struct {
		Tone
		TrebleOctave Octave // Helps render the time signature for treble clef
		BassOctave   Octave // Helps render the time signature for bass clef
	}

	Key struct {
		Major Tone
		Minor Tone
		Sig   []Signature
	}
)

const (
	Major Mode = iota
	Minor
)

func (m Mode) String() string {
	if m == Major {
		return "Major"
	} else if m == Minor {
		return "Minor"
	} else {
		return ""
	}
}

const (
	Natural Accidental = iota
	Flat
	Sharp
)

func (a Accidental) String() string {
	if a == Flat {
		return "b"
	} else if a == Sharp {
		return "#"
	} else {
		return ""
	}
}

const (
	Whole Duration = 1 << iota
	Half
	Quarter
	Eight
	Sixteenth
	ThirtySecond
)

const (
	A Class = "A"
	B Class = "B"
	C Class = "C"
	D Class = "D"
	E Class = "E"
	F Class = "F"
	G Class = "G"
)

func (t Tone) String() string {
	return string(t.Class) + t.Accidental.String()
}
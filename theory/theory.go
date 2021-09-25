// Package theory implements music theory. There is also support for associating
// MIDI numbers with music theory. It vaguely corresponds to the MusicXML data
// structures.
//
// Note that this is inahga's(tm) music theory. My understanding of actual music
// theory is pretty bad :).
package theory

type (
	Octave     int
	Duration   float64
	Accidental int
	Mode       int
	Class      int
	Clef       int

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

	KeySignature struct {
		Tone
		TrebleOctave Octave // Helps render the time signature for treble clef
		BassOctave   Octave // Helps render the time signature for bass clef
	}

	Key struct {
		Major Tone
		Minor Tone
		Sig   []KeySignature
	}

	Time struct {
		BeatsPerBar, BeatUnit int
	}

	Measure []Note

	Part struct {
		Clef
		Time
		Key
		Measures []Measure
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
	}
	return ""
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
	}
	return ""
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
	C Class = iota
	D
	E
	F
	G
	A
	B
)

const (
	Treble Clef = iota
	Bass
)

var classes = []string{"C", "D", "E", "F", "G", "A", "B"}

func (c Class) String() string {
	return classes[c]
}

func (t Tone) String() string {
	return t.Class.String() + t.Accidental.String()
}

func (c Clef) String() string {
	if c == Treble {
		return "Treble"
	} else if c == Bass {
		return "Bass"
	}
	return ""
}

func N(c Class, a Accidental, o Octave, d Duration) Note {
	return Note{
		Tone: Tone{
			Class:      c,
			Accidental: a,
		},
		Octave:   o,
		Duration: d,
	}
}

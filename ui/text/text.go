package text

// TextUI communicates to the user strictly through text, asynchronously. It is not
// to be confused with a TUI, which draws graphics using terminal elements.
//
// This UI is meant for debugging.
type TextUI struct{}

func New() (*TextUI, error) {
	return &TextUI{}, nil
}

func (t *TextUI) Start() error { return nil }

func (t *TextUI) Main() {
	select {}
}

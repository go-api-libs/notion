package notion

import "strings"

// PlainText returns the plain text of all the rich texts.
func (ts RichTexts) PlainText() string {
	b := &strings.Builder{}
	for _, t := range ts {
		b.WriteString(t.PlainText)
	}

	return b.String()
}

package notion

import "strings"

// plainText returns the plain text of all the rich texts.
func (ts RichTexts) plainText() string {
	ss := make([]string, len(ts))
	for i, t := range ts {
		ss[i] = t.PlainText
	}

	return strings.Join(ss, "")
}

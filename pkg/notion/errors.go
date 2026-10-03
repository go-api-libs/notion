package notion

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Error formats as e.g. "notion: 400 validation_error: body failed validation (key: value)".
func (e *Error) Error() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "notion: %d %s: %s", e.Status, e.Code, e.Message)

	if len(e.AdditionalData) == 0 {
		return b.String()
	}

	b.WriteString(" (")

	// sorted, so the same error always reads the same
	for i, key := range slices.Sorted(maps.Keys(e.AdditionalData)) {
		if i > 0 {
			b.WriteString("; ")
		}

		data := e.AdditionalData[key]
		fmt.Fprintf(b, "%s: %s%s", key, data.String, strings.Join(data.String2, ", "))
	}

	b.WriteByte(')')

	return b.String()
}

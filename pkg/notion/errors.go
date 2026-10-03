package notion

import (
	"fmt"
	"maps"
	"net/http"
	"slices"
	"strings"
)

func (e *ErrorAPI) Error() string { return e.format(e.Status, string(e.Code)) }

func (e *ErrorOauth) Error() string { return e.format(e.Status, string(e.Code)) }

func (e *publicApiCommonErrorResponse) format(status int, code string) string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%d %s - %s: obj - %s; msg - %s", status, http.StatusText(status), code, e.Object, e.Message)

	if e.AdditionalData != nil {
		// sorted, so the same error always reads the same
		for _, key := range slices.Sorted(maps.Keys(*e.AdditionalData)) {
			data := (*e.AdditionalData)[key]
			fmt.Fprintf(b, "; %s - ", key)

			if data.String != nil {
				b.WriteString(*data.String)
			}

			if data.String2 != nil {
				b.WriteString(strings.Join(*data.String2, ", "))
			}
		}
	}

	return b.String()
}

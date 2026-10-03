package notion

import (
	"fmt"
	"net/http"
	"strings"
)

func (e *ErrorAPI) Error() string {
	b := &strings.Builder{}
	fmt.Fprintf(b, "%d %s - %s: obj - %s; msg - %s", e.Status, http.StatusText(e.Status), e.Code, e.Object, e.Message)

	if e.AdditionalData != nil {
		for key, data := range *e.AdditionalData {
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

// TODO: after consolidating error_oauth_*, replace the following with a real error implementation
func (e *error_oauth_400) Error() string { return "" }
func (e *error_oauth_401) Error() string { return "" }
func (e *error_oauth_403) Error() string { return "" }
func (e *error_oauth_500) Error() string { return "" }

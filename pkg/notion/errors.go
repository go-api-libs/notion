package notion

import (
	"fmt"
	"net/http"
	"strings"
)

func (e *error_api_400) Error() string {
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

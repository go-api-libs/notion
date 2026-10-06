package notion

import (
	"encoding/json/v2"
	"testing"
)

func TestFilterFlatten(t *testing.T) {
	const (
		today     = `{"property":"Blch","type":"date","date":{"equals":"today"}}`
		yesterday = `{"property":"Blch","type":"date","date":{"equals":"yesterday"}}`
		notDone   = `{"property":"@F<W","type":"select","select":{"does_not_equal":"Done"}}`
	)

	for _, tc := range []struct{ in, want string }{
		{today, today},
		{`{"or":[` + today + `]}`, today},
		{`{"and":[{"and":[` + today + `,` + notDone + `]}]}`, `{"and":[` + today + `,` + notDone + `]}`},
		// the board view's filter, three levels deep
		{
			`{"or":[{"or":[` + today + `,{"and":[` + yesterday + `,` + notDone + `]}]},{"and":[` + notDone + `]}]}`,
			`{"or":[` + today + `,{"and":[` + yesterday + `,` + notDone + `]},` + notDone + `]}`,
		},
	} {
		var f Filter
		if err := json.Unmarshal([]byte(tc.in), &f); err != nil {
			t.Fatal(err)
		}

		got, err := json.Marshal(f.Flatten())
		if err != nil {
			t.Fatal(err)
		}

		if string(got) != tc.want {
			t.Errorf("Flatten(%s)\n got  %s\n want %s", tc.in, got, tc.want)
		}
	}
}

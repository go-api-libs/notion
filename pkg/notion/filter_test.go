package notion

import (
	"encoding/json/v2"
	"reflect"
	"strings"
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

// documented is, by the types of property they filter, each condition Notion documents for a filter of one.
var documented = []struct{ types, conditions []string }{
	{[]string{"title", "rich_text", "url", "email", "phone_number"}, []string{
		`"equals":"Tuscan kale"`, `"does_not_equal":"Tuscan kale"`, `"contains":"kale"`, `"does_not_contain":"kale"`,
		`"starts_with":"Tuscan"`, `"ends_with":"kale"`, `"is_empty":true`, `"is_not_empty":true`,
	}},
	{[]string{"number"}, []string{
		`"equals":42`, `"does_not_equal":42`, `"greater_than":42`, `"less_than":42`,
		`"greater_than_or_equal_to":42`, `"less_than_or_equal_to":42`, `"is_empty":true`, `"is_not_empty":true`,
	}},
	{[]string{"unique_id"}, []string{`"equals":42`, `"greater_than":42`, `"less_than_or_equal_to":42`}},
	{[]string{"checkbox"}, []string{`"equals":true`, `"does_not_equal":true`}},
	{[]string{"select", "status"}, []string{`"equals":"Done"`, `"does_not_equal":"Done"`, `"is_empty":true`, `"is_not_empty":true`}},
	{[]string{"multi_select"}, []string{`"contains":"A"`, `"does_not_contain":"A"`, `"is_empty":true`, `"is_not_empty":true`}},
	{[]string{"date", "created_time", "last_edited_time"}, []string{
		`"equals":"2021-05-10"`, `"before":"2021-05-10"`, `"after":"2021-05-10"`,
		`"on_or_before":"2021-05-10"`, `"on_or_after":"2021-05-10"`,
		`"this_week":{}`, `"past_week":{}`, `"past_month":{}`, `"past_year":{}`,
		`"next_week":{}`, `"next_month":{}`, `"next_year":{}`, `"is_empty":true`, `"is_not_empty":true`,
	}},
	{[]string{"people", "created_by", "last_edited_by"}, []string{
		`"contains":"6c574cee-ca68-41c8-86e0-1b9e992689fb"`, `"does_not_contain":"6c574cee-ca68-41c8-86e0-1b9e992689fb"`,
		`"is_empty":true`, `"is_not_empty":true`,
	}},
	{[]string{"files"}, []string{`"is_empty":true`, `"is_not_empty":true`}},
	{[]string{"relation"}, []string{
		`"contains":"dd456007-6c66-4bba-957e-ea501dcda3a6"`, `"does_not_contain":"dd456007-6c66-4bba-957e-ea501dcda3a6"`,
		`"is_empty":true`, `"is_not_empty":true`,
	}},
	{[]string{"formula"}, []string{
		`"string":{"contains":"tax"}`, `"checkbox":{"equals":true}`, `"number":{"greater_than":10}`,
		`"date":{"after":"2021-05-10"}`,
	}},
	{[]string{"rollup"}, []string{
		`"any":{"rich_text":{"contains":"Take Fig on a walk"}}`, `"every":{"number":{"greater_than":10}}`,
		`"none":{"checkbox":{"equals":true}}`, `"number":{"equals":10}`, `"date":{"after":"2021-05-10"}`,
	}},
}

// TestFilter_Documented reads every condition Notion documents for each type of property, and writes it back as it
// was, saying its type: the filter a query sends is the one written, whichever of a union's alternatives holds it.
// So it does a timestamp's and a compound filter.
func TestFilter_Documented(t *testing.T) {
	t.Parallel()

	var filters []struct{ in, want string }

	for _, d := range documented {
		for _, typ := range d.types {
			for _, c := range d.conditions {
				filters = append(filters, struct{ in, want string }{
					`{"property":"P","` + typ + `":{` + c + `}}`,
					`{"property":"P","type":"` + typ + `","` + typ + `":{` + c + `}}`,
				})
			}
		}
	}

	const done, a, b = `{"property":"Done","checkbox":{"equals":true}}`,
		`{"property":"Tags","multi_select":{"contains":"A"}}`, `{"property":"Tags","multi_select":{"contains":"B"}}`

	typed := func(f string) string {
		return strings.Replace(strings.Replace(f, `"checkbox":`, `"type":"checkbox","checkbox":`, 1),
			`"multi_select":`, `"type":"multi_select","multi_select":`, 1)
	}

	for _, f := range []string{
		`{"timestamp":"created_time","created_time":{"past_week":{}}}`,
		`{"timestamp":"last_edited_time","last_edited_time":{"after":"2021-05-10"}}`,
	} {
		filters = append(filters, struct{ in, want string }{f, f})
	}

	filters = append(filters, struct{ in, want string }{
		`{"and":[` + done + `,{"or":[` + a + `,` + b + `]}]}`,
		`{"and":[` + typed(done) + `,{"or":[` + typed(a) + `,` + typed(b) + `]}]}`,
	})

	for _, tc := range filters {
		roundTrip[Filter](t, tc.in, tc.want)
	}
}

// roundTrip reads in as a T and writes it back, as want.
func roundTrip[T any](t *testing.T, in, want string) {
	t.Helper()

	var v T
	if err := json.Unmarshal([]byte(in), &v); err != nil {
		t.Errorf("reading %s: %v", in, err)
		return
	}

	out, err := json.Marshal(v)
	if err != nil {
		t.Errorf("writing %s: %v", in, err)
		return
	}

	if !sameJSON(t, want, string(out)) {
		t.Errorf("%s\n is written as %s\n       want %s", in, out, want)
	}
}

// sameJSON reports whether a and b are the same JSON value, whatever the order of their members.
func sameJSON(t *testing.T, a, b string) bool {
	t.Helper()

	var va, vb any
	if err := json.Unmarshal([]byte(a), &va); err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal([]byte(b), &vb); err != nil {
		t.Fatal(err)
	}

	return reflect.DeepEqual(va, vb)
}

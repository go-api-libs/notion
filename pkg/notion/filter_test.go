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

// TestFilter_Documented reads every condition Notion documents for each type of property, and writes it back as it
// was, saying its type: the filter a query sends is the one written, whichever of a union's alternatives holds it.
func TestFilter_Documented(t *testing.T) {
	t.Parallel()

	var filters []struct{ in, want string }

	add := func(property, typ string, conditions ...string) {
		for _, c := range conditions {
			filters = append(filters, struct{ in, want string }{
				`{"property":"` + property + `","` + typ + `":{` + c + `}}`,
				`{"property":"` + property + `","type":"` + typ + `","` + typ + `":{` + c + `}}`,
			})
		}
	}

	text := []string{
		`"equals":"Tuscan kale"`, `"does_not_equal":"Tuscan kale"`, `"contains":"kale"`, `"does_not_contain":"kale"`,
		`"starts_with":"Tuscan"`, `"ends_with":"kale"`, `"is_empty":true`, `"is_not_empty":true`,
	}
	for _, typ := range []string{"title", "rich_text", "url", "email", "phone_number"} {
		add("Name", typ, text...)
	}

	number := []string{
		`"equals":42`, `"does_not_equal":42`, `"greater_than":42`, `"less_than":42`,
		`"greater_than_or_equal_to":42`, `"less_than_or_equal_to":42`, `"is_empty":true`, `"is_not_empty":true`,
	}
	add("Cost", "number", number...)
	add("ID", "unique_id", `"equals":42`, `"greater_than":42`, `"less_than_or_equal_to":42`)

	add("Done", "checkbox", `"equals":true`, `"does_not_equal":true`)
	add("Stage", "select", `"equals":"Done"`, `"does_not_equal":"Done"`, `"is_empty":true`, `"is_not_empty":true`)
	add("Tags", "multi_select", `"contains":"A"`, `"does_not_contain":"A"`, `"is_empty":true`, `"is_not_empty":true`)
	add("State", "status", `"equals":"Done"`, `"does_not_equal":"Done"`, `"is_empty":true`, `"is_not_empty":true`)

	date := []string{
		`"equals":"2021-05-10"`, `"before":"2021-05-10"`, `"after":"2021-05-10"`,
		`"on_or_before":"2021-05-10"`, `"on_or_after":"2021-05-10"`,
		`"this_week":{}`, `"past_week":{}`, `"past_month":{}`, `"past_year":{}`,
		`"next_week":{}`, `"next_month":{}`, `"next_year":{}`, `"is_empty":true`, `"is_not_empty":true`,
	}
	add("Due", "date", date...)

	people := []string{
		`"contains":"6c574cee-ca68-41c8-86e0-1b9e992689fb"`, `"does_not_contain":"6c574cee-ca68-41c8-86e0-1b9e992689fb"`,
		`"is_empty":true`, `"is_not_empty":true`,
	}
	for _, typ := range []string{"people", "created_by", "last_edited_by"} {
		add("Owner", typ, people...)
	}

	add("Files", "files", `"is_empty":true`, `"is_not_empty":true`)
	add("Tasks", "relation", `"contains":"dd456007-6c66-4bba-957e-ea501dcda3a6"`,
		`"does_not_contain":"dd456007-6c66-4bba-957e-ea501dcda3a6"`, `"is_empty":true`, `"is_not_empty":true`)

	add("Score", "formula", `"string":{"contains":"tax"}`, `"checkbox":{"equals":true}`,
		`"number":{"greater_than":10}`, `"date":{"after":"2021-05-10"}`)
	add("Totals", "rollup", `"any":{"rich_text":{"contains":"Take Fig on a walk"}}`,
		`"every":{"number":{"greater_than":10}}`, `"none":{"checkbox":{"equals":true}}`,
		`"number":{"equals":10}`, `"date":{"after":"2021-05-10"}`)

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
		var f Filter
		if err := json.Unmarshal([]byte(tc.in), &f); err != nil {
			t.Errorf("reading %s: %v", tc.in, err)
			continue
		}

		out, err := json.Marshal(f)
		if err != nil {
			t.Errorf("writing %s: %v", tc.in, err)
			continue
		}

		if !sameJSON(t, tc.want, string(out)) {
			t.Errorf("%s\n is written as %s\n       want %s", tc.in, out, tc.want)
		}
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

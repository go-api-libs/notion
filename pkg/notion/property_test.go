package notion

import (
	"encoding/json/v2"
	"testing"
)

// TestPropertyItem_Documented reads a page's property of each type, as Notion answers for one, into the item of its
// type, whichever of the union's alternatives would read it too.
func TestPropertyItem_Documented(t *testing.T) {
	t.Parallel()

	const (
		user = `{"object":"user","id":"6c574cee-ca68-41c8-86e0-1b9e992689fb"}`
		text = `{"type":"text","text":{"content":"kale","link":null},"annotations":{"bold":false,"italic":false,` +
			`"strikethrough":false,"underline":false,"code":false,"color":"default"},"plain_text":"kale","href":null}`
	)

	for _, tc := range []struct {
		typ, value string
		is         func(PropertyItem) bool
	}{
		{"number", `42`, func(p PropertyItem) bool { return p.NumberPropertyItem != nil }},
		{"url", `"https://developers.notion.com"`, func(p PropertyItem) bool { return p.URLPropertyItem != nil }},
		{"select", `{"id":"c:1","name":"Done","color":"green"}`, func(p PropertyItem) bool { return p.SelectPropertyItem != nil }},
		{"multi_select", `[{"id":"c:1","name":"A","color":"red"}]`, func(p PropertyItem) bool { return p.MultiSelectPropertyItem != nil }},
		{"status", `{"id":"s:1","name":"In progress","color":"blue"}`, func(p PropertyItem) bool { return p.StatusPropertyItem != nil }},
		{"date", `{"start":"2021-05-10","end":null,"time_zone":null}`, func(p PropertyItem) bool { return p.DatePropertyItem != nil }},
		{"email", `"ada@example.com"`, func(p PropertyItem) bool { return p.EmailPropertyItem != nil }},
		{"phone_number", `"+1 555 0100"`, func(p PropertyItem) bool { return p.PhoneNumberPropertyItem != nil }},
		{"checkbox", `true`, func(p PropertyItem) bool { return p.CheckboxPropertyItem != nil }},
		{"created_by", user, func(p PropertyItem) bool { return p.CreatedByPropertyItem != nil }},
		{"created_time", `"2021-05-10T12:00:00.000Z"`, func(p PropertyItem) bool { return p.CreatedTimePropertyItem != nil }},
		{"last_edited_by", user, func(p PropertyItem) bool { return p.LastEditedByPropertyItem != nil }},
		{"last_edited_time", `"2021-05-10T12:00:00.000Z"`, func(p PropertyItem) bool { return p.LastEditedTimePropertyItem != nil }},
		{"formula", `{"type":"number","number":84}`, func(p PropertyItem) bool { return p.FormulaPropertyItem != nil }},
		{"unique_id", `{"prefix":"TASK","number":7}`, func(p PropertyItem) bool { return p.UniqueIDPropertyItem != nil }},
		{"title", text, func(p PropertyItem) bool { return p.TitlePropertyItem != nil }},
		{"rich_text", text, func(p PropertyItem) bool { return p.RichTextPropertyItem != nil }},
		{"people", user, func(p PropertyItem) bool { return p.PeoplePropertyItem != nil }},
		{"relation", `{"id":"dd456007-6c66-4bba-957e-ea501dcda3a6"}`, func(p PropertyItem) bool { return p.RelationPropertyItem != nil }},
	} {
		in := `{"object":"property_item","id":"p","type":"` + tc.typ + `","` + tc.typ + `":` + tc.value + `}`

		var item PropertyItem
		if err := json.Unmarshal([]byte(in), &item); err != nil {
			t.Errorf("reading a %s: %v", tc.typ, err)
		} else if !tc.is(item) {
			t.Errorf("a %s is read as another type: %+v", tc.typ, item)
		}
	}
}

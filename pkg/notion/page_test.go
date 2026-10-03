package notion

import (
	"encoding/json/v2"
	"testing"
)

func TestPageTitle(t *testing.T) {
	p := &Page{}
	// type comes first only because the generated decoder requires it, though Notion sends id first
	if err := json.Unmarshal([]byte(`{
		"object": "page",
		"properties": {
			"Status": {"type": "checkbox", "id": "s", "checkbox": true},
			"Name": {"type": "title", "id": "title", "title": [
				{"type": "text", "text": {"content": "Example ", "link": null}, "plain_text": "Example ", "href": null,
				 "annotations": {"bold": false, "italic": false, "strikethrough": false, "underline": false, "code": false, "color": "default"}},
				{"type": "text", "text": {"content": "Page", "link": null}, "plain_text": "Page", "href": null,
				 "annotations": {"bold": true, "italic": false, "strikethrough": false, "underline": false, "code": false, "color": "default"}}
			]}
		}
	}`), p, jsonOpts); err != nil {
		t.Fatal(err)
	}

	if got, want := p.Title(), "Example Page"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

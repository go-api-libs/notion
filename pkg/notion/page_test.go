package notion

import (
	"encoding/json/v2"
	"testing"
)

func TestPageTitle(t *testing.T) {
	p := &Page{}
	if err := json.Unmarshal([]byte(`{
		"object": "page",
		"properties": {
			"Status": {"id": "s", "type": "checkbox", "checkbox": true},
			"Name": {"id": "title", "type": "title", "title": [
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

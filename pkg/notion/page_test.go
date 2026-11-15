package notion

import (
	"encoding/json/v2"
	"net/http"
	"slices"
	"testing"

	"github.com/MarkRosemaker/cassette"
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

// TestPageOrDataSourcePage reads the entries of a recorded query of the board's data source: each is a page in full,
// with its title.
func TestPageOrDataSourcePage(t *testing.T) {
	ias, err := cassette.InteractionsReadFile("../../api/interactions.json")
	if err != nil {
		t.Fatal(err)
	}

	i := slices.IndexFunc(ias, func(ia *cassette.Interaction) bool {
		return ia.Request.Method == http.MethodPost &&
			ia.Request.URL.Path == "/v1/data_sources/9201f6db-2895-4450-9435-041541a181dd/query"
	})
	if i < 0 {
		t.Fatal("no query of the board's data source is recorded")
	}

	var list QueryList
	if err := json.Unmarshal(ias[i].Response.Body, &list, jsonOpts); err != nil {
		t.Fatal(err)
	}

	if len(list.Results) == 0 {
		t.Fatal("the recorded query has no entries")
	}

	for _, r := range list.Results {
		p := r.Page()
		if p == nil {
			t.Fatalf("an entry of a data source is no page: %+v", r)
		}

		if p.Title() == "" {
			t.Errorf("page %s has no title", p.ID)
		}
	}

	if (PageOrDataSource{}).Page() != nil {
		t.Error("a result that is no page is one")
	}

	if (&Page{}).Title() != "" {
		t.Error("a page without a title property has a title")
	}
}

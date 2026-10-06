package notion

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestPostDatabaseQueryAll(t *testing.T) {
	const (
		first  = "00000000-0000-0000-0000-000000000001"
		second = "00000000-0000-0000-0000-000000000002"
		third  = "00000000-0000-0000-0000-000000000003"
	)

	var cursors []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			StartCursor string `json:"start_cursor"`
			PageSize    int    `json:"page_size"`
		}
		if err := json.UnmarshalRead(r.Body, &body); err != nil {
			t.Error(err)
		}

		if body.PageSize != maxPageSize {
			t.Errorf("page size %d, want %d", body.PageSize, maxPageSize)
		}

		cursors = append(cursors, body.StartCursor)

		w.Header().Set("Content-Type", "application/json")

		page := func(id string) string { return `{"object":"page","id":"` + id + `"}` }
		if body.StartCursor == "" {
			_, _ = w.Write([]byte(`{"object":"list","type":"page_or_data_source","page_or_data_source":{},` +
				`"results":[` + page(first) + `,` + page(second) + `],"has_more":true,"next_cursor":"` + third + `"}`))
		} else {
			_, _ = w.Write([]byte(`{"object":"list","type":"page_or_data_source","page_or_data_source":{},` +
				`"results":[` + page(third) + `],"has_more":false,"next_cursor":null}`))
		}
	}))
	defer srv.Close()

	base, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewClient(WithBearer("test"), WithBaseURL(base), WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatal(err)
	}

	var ids []string

	for r, err := range c.PostDatabaseQueryAll(t.Context(), "data-source", nil, PostDatabaseQuery{}) {
		if err != nil {
			t.Fatal(err)
		}

		ids = append(ids, r.PageOrPartial.PartialPageObjectResponse.ID.String())
	}

	if want := []string{first, second, third}; len(ids) != 3 || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Errorf("got %v, want %v", ids, want)
	}

	if len(cursors) != 2 || cursors[0] != "" || cursors[1] != third {
		t.Errorf("sent cursors %q, want none, then %q", cursors, third)
	}
}

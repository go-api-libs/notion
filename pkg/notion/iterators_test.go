package notion

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"uuid"
)

func TestPostDatabaseQueryAll(t *testing.T) {
	first, second, third := uuid.New(), uuid.New(), uuid.New()

	var cursors []uuid.UUID

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			StartCursor uuid.UUID `json:"start_cursor"`
			PageSize    int       `json:"page_size"`
		}
		if err := json.UnmarshalRead(r.Body, &body); err != nil {
			t.Error(err)
		}

		if body.PageSize != maxPageSize {
			t.Errorf("page size %d, want %d", body.PageSize, maxPageSize)
		}

		cursors = append(cursors, body.StartCursor)

		w.Header().Set("Content-Type", "application/json")

		page := func(id uuid.UUID) string { return `{"object":"page","id":"` + id.String() + `"}` }
		if body.StartCursor == uuid.Nil() {
			_, _ = w.Write([]byte(`{"object":"list","type":"page_or_data_source","page_or_data_source":{},` +
				`"results":[` + page(first) + `,` + page(second) + `],"has_more":true,"next_cursor":"` + third.String() + `"}`))
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

	var ids []uuid.UUID

	dataSourceID := uuid.New()
	for r, err := range c.PostDatabaseQueryAll(t.Context(), dataSourceID, nil, PostDatabaseQuery{}) {
		if err != nil {
			t.Fatal(err)
		}

		ids = append(ids, r.PageOrPartial.PartialPage.ID)
	}

	if want := []uuid.UUID{first, second, third}; len(ids) != 3 || ids[0] != want[0] || ids[1] != want[1] || ids[2] != want[2] {
		t.Errorf("got %v, want %v", ids, want)
	}

	if len(cursors) != 2 || cursors[0] != uuid.Nil() || cursors[1] != third {
		t.Errorf("sent cursors %q, want none, then %q", cursors, third)
	}
}

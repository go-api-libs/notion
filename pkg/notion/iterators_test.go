package notion

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"strings"
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

	var ids []uuid.UUID

	dataSourceID := uuid.New()
	for r, err := range testClient(t, srv).PostDatabaseQueryAll(t.Context(), dataSourceID, nil, PostDatabaseQuery{}) {
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

// testClient is a client of the API srv serves.
func testClient(t *testing.T, srv *httptest.Server) *Client {
	t.Helper()

	base, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}

	c, err := NewClient(WithBearer("test"), WithBaseURL(base), WithHTTPClient(srv.Client()))
	if err != nil {
		t.Fatal(err)
	}

	return c
}

// list is a page of a list of objects of the type, with the cursor of the next page if there is one.
func list(typ, object string, ids []uuid.UUID, next uuid.UUID) string {
	items := make([]string, len(ids))
	for i, id := range ids {
		items[i] = `{"object":"` + object + `","id":"` + id.String() + `"}`
	}

	more, cursor := "false", "null"
	if next != uuid.Nil() {
		more, cursor = "true", `"`+next.String()+`"`
	}

	return `{"object":"list","type":"` + typ + `","` + typ + `":{},"results":[` + strings.Join(items, ",") +
		`],"has_more":` + more + `,"next_cursor":` + cursor + `}`
}

func TestGetBlockChildrenAll(t *testing.T) {
	page, first, second, third := uuid.New(), uuid.New(), uuid.New(), uuid.New()

	var cursors []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if want := "/blocks/" + page.String() + "/children"; r.URL.Path != want {
			t.Errorf("path %s, want %s", r.URL.Path, want)
		}

		if got := r.URL.Query().Get("page_size"); got != "100" {
			t.Errorf("page size %s, want 100", got)
		}

		cursor := r.URL.Query().Get("start_cursor")
		cursors = append(cursors, cursor)

		w.Header().Set("Content-Type", "application/json")

		if cursor == "" {
			_, _ = w.Write([]byte(list("block", "block", []uuid.UUID{first, second}, third)))
		} else {
			_, _ = w.Write([]byte(list("block", "block", []uuid.UUID{third}, uuid.Nil())))
		}
	}))
	defer srv.Close()

	var ids []uuid.UUID

	for b, err := range testClient(t, srv).GetBlockChildrenAll(t.Context(), page) {
		if err != nil {
			t.Fatal(err)
		}

		ids = append(ids, b.PartialBlock.ID)
	}

	if !slices.Equal(ids, []uuid.UUID{first, second, third}) {
		t.Errorf("got %v, want %v", ids, []uuid.UUID{first, second, third})
	}

	if !slices.Equal(cursors, []string{"", third.String()}) {
		t.Errorf("sent cursors %q, want none, then %q", cursors, third)
	}
}

func TestPostSearchAll(t *testing.T) {
	first, second := uuid.New(), uuid.New()

	var cursors []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Query       string  `json:"query"`
			StartCursor string  `json:"start_cursor"`
			PageSize    float64 `json:"page_size"`
		}
		if err := json.UnmarshalRead(r.Body, &body); err != nil {
			t.Error(err)
		}

		if body.Query != "kale" || body.PageSize != maxPageSize {
			t.Errorf("searched %q, %v at a time, want %q, %d", body.Query, body.PageSize, "kale", maxPageSize)
		}

		cursors = append(cursors, body.StartCursor)

		w.Header().Set("Content-Type", "application/json")

		if body.StartCursor == "" {
			_, _ = w.Write([]byte(list("page_or_data_source", "page", []uuid.UUID{first}, second)))
		} else {
			_, _ = w.Write([]byte(list("page_or_data_source", "page", []uuid.UUID{second}, uuid.Nil())))
		}
	}))
	defer srv.Close()

	var ids []uuid.UUID

	for r, err := range testClient(t, srv).PostSearchAll(t.Context(), nil, PostSearch{Query: "kale"}) {
		if err != nil {
			t.Fatal(err)
		}

		ids = append(ids, r.PageOrPartial.PartialPage.ID)
	}

	if !slices.Equal(ids, []uuid.UUID{first, second}) {
		t.Errorf("got %v, want %v", ids, []uuid.UUID{first, second})
	}

	if len(cursors) != 2 || cursors[1] != second.String() {
		t.Errorf("sent cursors %q, want none, then %q", cursors, second)
	}
}

// TestAll_Stops: breaking out of the loop fetches no further page, and an error ends it, yielded once.
func TestAll_Stops(t *testing.T) {
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		calls++

		w.Header().Set("Content-Type", "application/json")

		if calls > 1 {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"object":"error","status":400,"code":"validation_error","message":"bad cursor"}`))

			return
		}

		_, _ = w.Write([]byte(list("block", "block", []uuid.UUID{uuid.New(), uuid.New()}, uuid.New())))
	}))
	defer srv.Close()

	c := testClient(t, srv)

	for _, err := range c.GetBlockChildrenAll(t.Context(), uuid.New()) {
		if err != nil {
			t.Fatal(err)
		}

		break
	}

	if calls != 1 {
		t.Fatalf("breaking out after the first block fetched %d pages, want 1", calls)
	}

	calls = 0

	var (
		blocks int
		errs   []error
	)

	for _, err := range c.GetBlockChildrenAll(t.Context(), uuid.New()) {
		if err != nil {
			errs = append(errs, err)
			continue
		}

		blocks++
	}

	if blocks != 2 || len(errs) != 1 {
		t.Fatalf("got %d blocks and %d errors, want 2 blocks, then the error", blocks, len(errs))
	}

	if e, ok := errors.AsType[*Error](errs[0]); !ok || e.Code != "validation_error" {
		t.Errorf("got %v, want Notion's validation_error", errs[0])
	}
}

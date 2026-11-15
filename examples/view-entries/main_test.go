package main

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/MarkRosemaker/cassette"
	"github.com/go-api-libs/notion/pkg/notion"
)

// TestQueryOfView reads the board's view as recorded, and queries its data source with the query it makes of the
// view: three entries of it are the query the recordings hold, filtered and sorted like the view, which Notion
// answered. The client says each filter's type, which the recording leaves out: a type naming the member beside it
// says nothing that member does not.
func TestQueryOfView(t *testing.T) {
	ias, err := cassette.InteractionsReadFile("../../api/interactions.json")
	if err != nil {
		t.Fatal(err)
	}

	t.Setenv("NOTION_API_TOKEN", "test")

	replay := cassette.Replay(ias, cassette.WithMatcher(sameCall))

	c, err := notion.NewClient(notion.WithHTTPClient(&http.Client{Transport: replay}))
	if err != nil {
		t.Fatal(err)
	}

	dataSourceID, query, err := queryOfView(t.Context(), c, viewID)
	if err != nil {
		t.Fatal(err)
	}

	query.PageSize = 3

	list, err := c.PostDatabaseQuery(t.Context(), dataSourceID, nil, query)
	if err != nil {
		t.Fatalf("the view's query is not the one recorded: %v", err)
	}

	if len(list.Results) == 0 {
		t.Error("the view shows no entries")
	}
}

// sameCall reports why req is not the recorded call rec: another method, URL or body, a type that names the member
// beside it aside.
func sameCall(req *http.Request, rec cassette.Request) error {
	if req.Method != rec.Method || req.URL.String() != rec.URL.String() {
		return fmt.Errorf("%s %s is no %s %s", req.Method, req.URL, rec.Method, &rec.URL)
	}

	if req.Body == nil {
		return nil
	}

	sent, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader(sent))

	var a, b any
	if err := json.Unmarshal(sent, &a); err != nil {
		return err
	}

	if err := json.Unmarshal(rec.Body, &b); err != nil {
		return err
	}

	if !reflect.DeepEqual(untyped(a), untyped(b)) {
		return fmt.Errorf("sent %s, recorded %s", sent, rec.Body)
	}

	return nil
}

// untyped is v without the types that name the member beside them.
func untyped(v any) any {
	switch v := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(v))

		for k, x := range v {
			if typ, ok := x.(string); k == "type" && ok && v[typ] != nil {
				continue
			}

			out[k] = untyped(x)
		}

		return out
	case []any:
		out := make([]any, len(v))
		for i, x := range v {
			out[i] = untyped(x)
		}

		return out
	default:
		return v
	}
}

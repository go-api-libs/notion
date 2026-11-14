package main

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"time"
	"uuid"

	"github.com/MarkRosemaker/cassette"
)

var (
	// https://app.notion.com/p/fae-tools/457d00dcfc3a4ddba42ffe376530bee0?v=cb3265bc5d5e488fbd8a69221c107915
	dbID         = uuid.MustParse("457d00dcfc3a4ddba42ffe376530bee0")
	dataSourceID = uuid.MustParse("9201f6db-2895-4450-9435-041541a181dd")
	viewID       = uuid.MustParse("cb3265bc5d5e488fbd8a69221c107915")

	// https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3, a page of every kind of block
	examplePageID = uuid.MustParse("96245c8f178444a482ad1941127c3ec3")
)

const (
	pathInteractions = "api/interactions.json"
	baseURL          = "https://api.notion.com/v1/"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

// run records every call anew: the file holds what this run recorded and nothing else, so that a page the run no
// longer comes to goes with it.
func run(ctx context.Context) error {
	rec := cassette.Replay(nil, cassette.RecordMisses(paced{}))
	recorded := &http.Client{Transport: sign{rec}}

	record := func(r cassette.Request) (*cassette.Interaction, error) {
		if _, _, err := send(ctx, recorded, r); err != nil {
			return nil, err
		}

		ias := rec.Interactions()

		return ias[len(ias)-1], nil
	}

	view, err := record(cassette.Request{
		Method: http.MethodGet,
		URL:    mustParseURL("https://api.notion.com/v1/views/" + viewID.String()),
	})
	if err != nil {
		return err
	}

	// the entries of the view, as its filter and sorts select and order them
	viewQuery, err := queryOfView(view.Response.Body)
	if err != nil {
		return err
	}

	if _, err := record(viewQuery); err != nil {
		return err
	}

	for _, r := range []cassette.Request{
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/views?database_id=" + dbID.String())},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/databases/" + dbID.String())},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/data_sources/" + dataSourceID.String())},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/data_sources/" + dataSourceID.String() + "/templates")},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/blocks/" + dbID.String())},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/users/me")},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/users?page_size=3")},
		{Method: http.MethodPost, URL: mustParseURL("https://api.notion.com/v1/search"), Body: []byte(`{"page_size":3}`)},
		// an error, for the order of its members
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/pages/" + uuid.Nil().String())},
	} {
		if _, err := record(r); err != nil {
			return err
		}
	}

	query, err := record(cassette.Request{
		Method: http.MethodPost,
		URL:    mustParseURL("https://api.notion.com/v1/data_sources/" + dataSourceID.String() + "/query"),
		Body:   []byte(`{"page_size":3}`),
	})
	if err != nil {
		return err
	}

	// the first page of the data source, to record a page with its properties and blocks
	var list struct {
		Results []struct {
			ID string `json:"id"`
		} `json:"results"`
	}
	if err := json.Unmarshal(query.Response.Body, &list); err != nil {
		return fmt.Errorf("decoding query: %w", err)
	}

	if len(list.Results) == 0 {
		return errors.New("data source has no pages")
	}

	pageID := list.Results[0].ID

	for _, r := range []cassette.Request{
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/pages/" + pageID)},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/pages/" + pageID + "/properties/title")},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/blocks/" + pageID + "/children?page_size=3")},
		{Method: http.MethodGet, URL: mustParseURL("https://api.notion.com/v1/comments?block_id=" + pageID)},
	} {
		if _, err := record(r); err != nil {
			return err
		}
	}

	if err := recordExample(ctx, record, &http.Client{Transport: sign{paced{}}}); err != nil {
		return fmt.Errorf("example page: %w", err)
	}

	ias, err := once(ctx, rec.Interactions())
	if err != nil {
		return err
	}

	ias.TrimResponseHeaders()

	return ias.WriteFile(pathInteractions)
}

// once is the interactions without those a later one records the same call of, as cassette matches a call to a
// recording: the last of each call is kept, where it was.
func once(ctx context.Context, ias cassette.Interactions) (cassette.Interactions, error) {
	match := cassette.MatchRequest()

	var out cassette.Interactions

	for i, ia := range ias {
		req, err := ia.Request.Create(ctx)
		if err != nil {
			return nil, err
		}

		if !slices.ContainsFunc(ias[i+1:], func(later *cassette.Interaction) bool { return match(req, later.Request) == nil }) {
			out = append(out, ia)
		}
	}

	return out, nil
}

func mustParseURL(s string) url.URL {
	u, err := url.Parse(s)
	if err != nil {
		panic(err)
	}

	return *u
}

// queryOfView is the query of the data source the view shows, with the view's filter and sorts.
func queryOfView(view []byte) (cassette.Request, error) {
	var v struct {
		DataSourceID string         `json:"data_source_id"`
		Filter       jsontext.Value `json:"filter"`
		Sorts        jsontext.Value `json:"sorts"`
	}
	if err := json.Unmarshal(view, &v); err != nil {
		return cassette.Request{}, fmt.Errorf("decoding view: %w", err)
	}

	var filter any
	if len(v.Filter) > 0 {
		if err := json.Unmarshal(v.Filter, &filter); err != nil {
			return cassette.Request{}, fmt.Errorf("decoding filter: %w", err)
		}
	}

	body, err := json.Marshal(struct {
		Filter   any            `json:"filter,omitzero"`
		Sorts    jsontext.Value `json:"sorts,omitzero"`
		PageSize int            `json:"page_size"`
	}{unnest(filter), v.Sorts, 3}, json.Deterministic(true))
	if err != nil {
		return cassette.Request{}, err
	}

	return cassette.Request{
		Method: http.MethodPost,
		URL:    mustParseURL("https://api.notion.com/v1/data_sources/" + v.DataSourceID + "/query"),
		Body:   body,
	}, nil
}

// unnest lists the filters of a compound filter within one of the same kind, an "or" in an "or" or an "and" in an
// "and", in the outer one, which allows the same entries. A view nests filters deeper than a query may, two levels.
func unnest(filter any) any {
	f, ok := filter.(map[string]any)
	if !ok || len(f) != 1 {
		return filter
	}

	for op, children := range f {
		cs, ok := children.([]any)
		if op != "and" && op != "or" || !ok {
			return filter
		}

		var flat []any

		for _, c := range cs {
			c = unnest(c)
			if inner, ok := c.(map[string]any); ok && len(inner) == 1 {
				if same, ok := inner[op].([]any); ok {
					flat = append(flat, same...)
					continue
				}
			}

			flat = append(flat, c)
		}

		return map[string]any{op: flat}
	}

	return filter
}

// recordExample records the example page, a block of each shape in it and its subpages, and each database in it.
// Its lists of blocks are walked without being recorded, since a recording keeps only three items of each list.
func recordExample(ctx context.Context, record func(cassette.Request) (*cassette.Interaction, error), walk *http.Client) error {
	page := examplePageID.String()

	for _, r := range []cassette.Request{
		get("pages/" + page),
		get("pages/" + page + "/properties/title"),
		get("blocks/" + page + "/children?page_size=3"),
		get("comments?block_id=" + page),
	} {
		if _, err := record(r); err != nil {
			return err
		}
	}

	w := walker{ctx: ctx, record: record, walk: walk, seen: map[string]bool{}}

	return w.children(page)
}

// walker records a block of each shape it comes across.
type walker struct {
	ctx    context.Context
	record func(cassette.Request) (*cassette.Interaction, error)
	walk   *http.Client    // what lists the blocks it walks, unrecorded
	seen   map[string]bool // the shapes recorded
}

// children walks the blocks within the block or page id.
func (w *walker) children(id string) error {
	for cursor := ""; ; {
		url := baseURL + "blocks/" + id + "/children?page_size=100"
		if cursor != "" {
			url += "&start_cursor=" + cursor
		}

		status, body, err := send(w.ctx, w.walk, cassette.Request{Method: http.MethodGet, URL: mustParseURL(url)})
		if err != nil {
			return err
		}

		if status != http.StatusOK {
			return fmt.Errorf("GET %s: %d %s", url, status, body)
		}

		var list struct {
			Results    []map[string]any `json:"results"`
			HasMore    bool             `json:"has_more"`
			NextCursor string           `json:"next_cursor"`
		}
		if err := json.Unmarshal(body, &list); err != nil {
			return fmt.Errorf("decoding children of %s: %w", id, err)
		}

		for _, b := range list.Results {
			if err := w.block(b); err != nil {
				return err
			}
		}

		if !list.HasMore {
			return nil
		}

		cursor = list.NextCursor
	}
}

// block records b if no block of its shape was, and what it holds.
func (w *walker) block(b map[string]any) error {
	id, _ := b["id"].(string)
	typ, _ := b["type"].(string)

	if shape, err := json.Marshal(shapeOf(b), json.Deterministic(true)); err != nil {
		return err
	} else if !w.seen[string(shape)] {
		w.seen[string(shape)] = true

		if _, err := w.record(get("blocks/" + id)); err != nil {
			return err
		}
	}

	switch typ {
	case "child_page":
		if _, err := w.record(get("pages/" + id)); err != nil {
			return err
		}
	case "child_database":
		return w.database(id)
	}

	if b["has_children"] == true || typ == "child_page" {
		return w.children(id)
	}

	return nil
}

// database records the database id, and each of its data sources with a few of its entries.
func (w *walker) database(id string) error {
	ia, err := w.record(get("databases/" + id))
	if err != nil {
		return err
	}

	var db struct {
		DataSources []struct {
			ID string `json:"id"`
		} `json:"data_sources"`
	}
	if err := json.Unmarshal(ia.Response.Body, &db); err != nil {
		return fmt.Errorf("decoding database %s: %w", id, err)
	}

	for _, ds := range db.DataSources {
		if _, err := w.record(get("data_sources/" + ds.ID)); err != nil {
			return err
		}

		query, err := w.record(cassette.Request{
			Method: http.MethodPost,
			URL:    mustParseURL(baseURL + "data_sources/" + ds.ID + "/query"),
			Body:   []byte(`{"page_size":3}`),
		})
		if err != nil {
			return err
		}

		if err := w.properties(query.Response.Body); err != nil {
			return err
		}

		if err := w.templates(ds.ID); err != nil {
			return err
		}
	}

	return nil
}

// templates records the templates of the data source, and walks each like a subpage.
func (w *walker) templates(dataSourceID string) error {
	ia, err := w.record(get("data_sources/" + dataSourceID + "/templates"))
	if err != nil {
		return err
	}

	var list struct {
		Templates []struct {
			ID string `json:"id"`
		} `json:"templates"`
	}
	if err := json.Unmarshal(ia.Response.Body, &list); err != nil {
		return fmt.Errorf("decoding templates of %s: %w", dataSourceID, err)
	}

	for _, t := range list.Templates {
		if _, err := w.record(get("pages/" + t.ID)); err != nil {
			return err
		}

		if err := w.children(t.ID); err != nil {
			return err
		}
	}

	return nil
}

// properties records each property of the first entry of a query, as the endpoint for one property returns it.
// Notion gives property IDs already escaped for a URL.
func (w *walker) properties(query []byte) error {
	var list struct {
		Results []struct {
			ID         string `json:"id"`
			Properties map[string]struct {
				ID string `json:"id"`
			} `json:"properties"`
		} `json:"results"`
	}
	if err := json.Unmarshal(query, &list); err != nil {
		return fmt.Errorf("decoding query: %w", err)
	}

	if len(list.Results) == 0 {
		return nil
	}

	page := list.Results[0]

	for _, name := range slices.Sorted(maps.Keys(page.Properties)) {
		if _, err := w.record(get("pages/" + page.ID + "/properties/" + page.Properties[name].ID)); err != nil {
			return err
		}
	}

	return nil
}

// noise are the members whose values tell blocks apart but not their shape: what each says, not how.
var noise = map[string]bool{
	"id": true, "created_time": true, "last_edited_time": true, "created_by": true, "last_edited_by": true,
	"parent": true, "plain_text": true, "content": true, "url": true, "expiry_time": true, "expression": true,
	"title": true, "caption": true, "name": true, "page_id": true, "database_id": true, "block_id": true,
	"start": true, "end": true, "emoji": true,
}

// shapeOf is v without its noise, and with each of its arrays' items once.
func shapeOf(v any) any {
	switch v := v.(type) {
	case map[string]any:
		m := make(map[string]any, len(v))
		for k, x := range v {
			if noise[k] {
				m[k] = nil
			} else {
				m[k] = shapeOf(x)
			}
		}

		return m
	case []any:
		var items []any

		seen := map[string]bool{}

		for _, x := range v {
			item := shapeOf(x)
			if b, err := json.Marshal(item, json.Deterministic(true)); err == nil && !seen[string(b)] {
				seen[string(b)] = true

				items = append(items, item)
			}
		}

		return items
	}

	return v
}

// get is a GET request of the path below the API's URL.
func get(path string) cassette.Request {
	return cassette.Request{
		Method: http.MethodGet,
		URL:    mustParseURL(baseURL + path),
	}
}

// send sends r through c, and returns the status and body of the response.
func send(ctx context.Context, c *http.Client, r cassette.Request) (int, []byte, error) {
	req, err := r.Create(ctx)
	if err != nil {
		return 0, nil, err
	}

	rsp, err := c.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("%s %s: %w", r.Method, &r.URL, err)
	}
	defer rsp.Body.Close()

	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("%s %s: %w", r.Method, &r.URL, err)
	}

	return rsp.StatusCode, body, nil
}

// sign gives each request the credentials and version of the API Notion asks for, before next, the recorder, sees
// it: a recording keeps the headers a call was made with, the credentials masked.
type sign struct{ next http.RoundTripper }

// RoundTrip implements [http.RoundTripper].
func (s sign) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Set("Authorization", "Bearer "+os.Getenv("NOTION_API_KEY"))
	req.Header.Set("Notion-Version", "2026-03-11")

	if req.Body != nil && req.Body != http.NoBody {
		req.Header.Set("Content-Type", "application/json")
	}

	return s.next.RoundTrip(req)
}

// pace spaces requests out to the three a second Notion allows on average.
var pace = time.Tick(time.Second / 3)

// paced sends a request to Notion at the pace it allows, and again after the wait it asks for when it limits the
// rate.
type paced struct{}

// RoundTrip implements [http.RoundTripper].
func (paced) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()

	for {
		select {
		case <-pace:
		case <-ctx.Done():
			return nil, ctx.Err()
		}

		rsp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil || rsp.StatusCode != http.StatusTooManyRequests {
			return rsp, err
		}

		wait, _ := strconv.Atoi(rsp.Header.Get("Retry-After"))
		rsp.Body.Close() //nolint:errcheck // only its status was read

		if req.GetBody != nil {
			if req.Body, err = req.GetBody(); err != nil {
				return nil, err
			}
		}

		select {
		case <-time.After(time.Duration(max(wait, 1)) * time.Second):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

package main

import (
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/MarkRosemaker/cassette"
)

// notion is a workspace as Notion answers for it, for the recorder to walk: the example page holds two paragraphs of
// one shape, a heading, a subpage and a database, whose data source has an entry and a template. The first call of
// the bot's user is answered with a 429.
type notion struct{ limited atomic.Bool }

const (
	childDB    = "11111111-1111-1111-1111-111111111111"
	childDS    = "22222222-2222-2222-2222-222222222222"
	row        = "33333333-3333-3333-3333-333333333333"
	template   = "44444444-4444-4444-4444-444444444444"
	subpage    = "55555555-5555-5555-5555-555555555555"
	boardEntry = "66666666-6666-6666-6666-666666666666"
)

func (n *notion) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	write := func(body string) { _, _ = rec.WriteString(body) }

	rec.Header().Set("Content-Type", "application/json")

	path := strings.TrimPrefix(req.URL.Path, "/v1/")

	block := func(id, typ, content string, children bool) string {
		return `{"object":"block","id":"` + id + `","type":"` + typ + `","has_children":` +
			map[bool]string{true: "true", false: "false"}[children] + `,"` + typ + `":` + content + `}`
	}

	switch {
	case path == "users/me" && !n.limited.Swap(true):
		rec.Header().Set("Retry-After", "0")
		rec.WriteHeader(http.StatusTooManyRequests)
		write(`{"object":"error","status":429,"code":"rate_limited"}`)
	case path == "views/"+viewID.String():
		write(`{"object":"view","data_source_id":"` + dataSourceID.String() + `",` +
			`"filter":{"or":[{"or":[{"property":"Done","checkbox":{"equals":true}}]}]},"sorts":[]}`)
	case path == "data_sources/"+dataSourceID.String()+"/query":
		write(`{"results":[{"id":"` + boardEntry + `","properties":{"Name":{"id":"title"}}}]}`)
	case path == "data_sources/"+childDS+"/query":
		write(`{"results":[{"id":"` + row + `","properties":{"Name":{"id":"title"},"Done":{"id":"a%3Fb"}}}]}`)
	case path == "databases/"+childDB:
		write(`{"object":"database","data_sources":[{"id":"` + childDS + `"}]}`)
	case path == "data_sources/"+childDS+"/templates":
		write(`{"templates":[{"id":"` + template + `"}]}`)
	case path == "blocks/"+examplePageID.String()+"/children" && req.URL.Query().Get("page_size") == "100":
		write(`{"results":[` + strings.Join([]string{
			block("aaaaaaaa-0000-0000-0000-000000000001", "paragraph", `{"rich_text":[{"plain_text":"one"}]}`, false),
			block("aaaaaaaa-0000-0000-0000-000000000002", "paragraph", `{"rich_text":[{"plain_text":"two"}]}`, false),
			block("aaaaaaaa-0000-0000-0000-000000000003", "heading_1", `{"rich_text":[],"is_toggleable":false}`, false),
			block(subpage, "child_page", `{"title":"Sub"}`, false),
			block(childDB, "child_database", `{"title":"Tasks"}`, false),
		}, ",") + `],"has_more":false}`)
	case strings.HasSuffix(path, "/children") && req.URL.Query().Get("page_size") == "100":
		write(`{"results":[],"has_more":false}`)
	default:
		write(`{"object":"list","results":[]}`)
	}

	rsp := rec.Result()
	rsp.Request = req

	return rsp, nil
}

// TestRun walks the workspace and records each call once: a block of each shape, the subpage, the database with its
// data source, the properties of its entry and its template, and the board's view with its filter as a query takes
// it, unnested. The lists of blocks walked are not recorded, and a call Notion limits is made again.
func TestRun(t *testing.T) {
	transport := http.DefaultTransport
	http.DefaultTransport = &notion{}
	tick := pace
	pace = time.Tick(time.Millisecond)

	t.Cleanup(func() { http.DefaultTransport, pace = transport, tick })

	t.Chdir(t.TempDir())

	if err := os.Mkdir("api", 0o755); err != nil {
		t.Fatal(err)
	}

	if err := run(t.Context()); err != nil {
		t.Fatal(err)
	}

	ias, err := cassette.InteractionsReadFile(pathInteractions)
	if err != nil {
		t.Fatal(err)
	}

	calls := map[string]*cassette.Interaction{}

	for _, ia := range ias {
		call := ia.Request.Method + " " + strings.TrimPrefix(ia.Request.URL.String(), baseURL) + " " + string(ia.Request.Body)
		if calls[call] != nil {
			t.Errorf("%s is recorded twice", call)
		}

		calls[call] = ia
	}

	recorded := func(method, path string) *cassette.Interaction {
		for call, ia := range calls {
			if strings.HasPrefix(call, method+" "+path+" ") {
				return ia
			}
		}

		return nil
	}

	for _, path := range []string{
		"blocks/aaaaaaaa-0000-0000-0000-000000000001",
		"blocks/aaaaaaaa-0000-0000-0000-000000000003",
		"pages/" + subpage,
		"databases/" + childDB,
		"data_sources/" + childDS,
		"data_sources/" + childDS + "/templates",
		"pages/" + row + "/properties/title",
		"pages/" + row + "/properties/a%3Fb",
		"pages/" + template,
	} {
		if recorded(http.MethodGet, path) == nil {
			t.Errorf("GET %s is not recorded", path)
		}
	}

	if recorded(http.MethodGet, "blocks/aaaaaaaa-0000-0000-0000-000000000002") != nil {
		t.Error("a second paragraph of the same shape is recorded")
	}

	for call := range calls {
		if strings.Contains(call, "page_size=100") {
			t.Errorf("a list of blocks walked is recorded: %s", call)
		}
	}

	me := recorded(http.MethodGet, "users/me")
	if me == nil || me.Response.StatusCode != http.StatusOK {
		t.Errorf("the call Notion limited is recorded as %+v, want its answer when made again", me)
	}

	if v := me.Request.Headers.Get("Notion-Version"); v == "" {
		t.Error("a recorded call does not say the version of the API it was made with")
	}

	if auth := me.Request.Headers.Get("Authorization"); strings.Trim(auth, "Bearer *") != "" {
		t.Errorf("a recorded call keeps its credentials: %q", auth)
	}

	// the board's data source is queried as its view shows it, and for a few of its entries
	var query *cassette.Interaction

	for call, ia := range calls {
		if strings.HasPrefix(call, http.MethodPost+" data_sources/"+dataSourceID.String()+"/query ") && strings.Contains(call, "filter") {
			query = ia
		}
	}

	if query == nil {
		t.Fatal("the board's view's query is not recorded")
	}

	var body, want any
	if err := json.Unmarshal(query.Request.Body, &body); err != nil {
		t.Fatal(err)
	}

	if err := json.Unmarshal([]byte(`{"filter":{"or":[{"property":"Done","checkbox":{"equals":true}}]},"sorts":[],"page_size":3}`), &want); err != nil {
		t.Fatal(err)
	}

	if got, _ := json.Marshal(body, json.Deterministic(true)); string(got) != string(must(json.Marshal(want, json.Deterministic(true)))) {
		t.Errorf("the view's query is %s, want its filter unnested: %s", query.Request.Body, must(json.Marshal(want)))
	}
}

func must(b []byte, err error) []byte {
	if err != nil {
		panic(err)
	}

	return b
}

package main

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"uuid"

	"github.com/MarkRosemaker/openapi-enrich/cassette"
)

var (
	// https://app.notion.com/p/fae-tools/457d00dcfc3a4ddba42ffe376530bee0?v=cb3265bc5d5e488fbd8a69221c107915
	dbID         = uuid.MustParse("457d00dcfc3a4ddba42ffe376530bee0")
	dataSourceID = uuid.MustParse("9201f6db-2895-4450-9435-041541a181dd")
	viewID       = uuid.MustParse("cb3265bc5d5e488fbd8a69221c107915")
)

const (
	pathInteractions = "api/interactions.json"
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

func run(ctx context.Context) error {
	ias, err := cassette.InteractionsReadFile(pathInteractions)
	if err != nil {
		return err
	}

	ias = slices.DeleteFunc(ias, eqTo(cassette.Interaction{}))

	record := func(r cassette.Request) (cassette.Interaction, error) {
		ia, err := do(ctx, r)
		if err != nil {
			return ia, fmt.Errorf("%s %s: %w", r.Method, r.URL, err)
		}

		ias = slices.DeleteFunc(ias, eqTo(ia))
		ias = append(ias, ia)

		return ia, nil
	}

	view, err := record(cassette.Request{Method: http.MethodGet, URL: "https://api.notion.com/v1/views/" + viewID.String()})
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
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/views?database_id=" + dbID.String()},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/databases/" + dbID.String()},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/data_sources/" + dataSourceID.String()},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/blocks/" + dbID.String()},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/users/me"},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/users?page_size=3"},
		{Method: http.MethodPost, URL: "https://api.notion.com/v1/search", Body: []byte(`{"page_size":3}`)},
		// an error, for the order of its members
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/pages/" + uuid.Nil().String()},
	} {
		if _, err := record(r); err != nil {
			return err
		}
	}

	query, err := record(cassette.Request{
		Method: http.MethodPost,
		URL:    "https://api.notion.com/v1/data_sources/" + dataSourceID.String() + "/query",
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
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/pages/" + pageID},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/pages/" + pageID + "/properties/title"},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/blocks/" + pageID + "/children?page_size=3"},
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/comments?block_id=" + pageID},
	} {
		if _, err := record(r); err != nil {
			return err
		}
	}

	ias.Mask()
	ias.TrimResponseHeaders()
	ias.TrimResponseBodies(3)

	return ias.WriteFile(pathInteractions)
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

	body, err := json.Marshal(struct {
		Filter   jsontext.Value `json:"filter,omitzero"`
		Sorts    jsontext.Value `json:"sorts,omitzero"`
		PageSize int            `json:"page_size"`
	}{v.Filter, v.Sorts, 3})
	if err != nil {
		return cassette.Request{}, err
	}

	return cassette.Request{
		Method: http.MethodPost,
		URL:    "https://api.notion.com/v1/data_sources/" + v.DataSourceID + "/query",
		Body:   body,
	}, nil
}

// do sends r to Notion and records the interaction.
func do(ctx context.Context, r cassette.Request) (cassette.Interaction, error) {
	var reqBody io.Reader
	if len(r.Body) > 0 {
		reqBody = bytes.NewReader(r.Body)
	}

	req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, reqBody)
	if err != nil {
		return cassette.Interaction{}, err
	}

	if len(r.Headers) > 0 {
		req.Header = r.Headers.Clone()
	}

	apiKey := os.Getenv("NOTION_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Notion-Version", "2026-03-11")

	if reqBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	r, err = cassette.NewRequest(req)
	if err != nil {
		return cassette.Interaction{}, err
	}

	ia := cassette.Interaction{Request: r}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ia, err
	}

	ia.Response, err = cassette.NewResponse(rsp)
	if err != nil {
		return ia, fmt.Errorf("recording response: %w", err)
	}

	return ia, nil
}

// eqTo reports whether an interaction is a's: the same request URL and body, so that queries of the same URL are
// recorded apart.
func eqTo(a cassette.Interaction) func(cassette.Interaction) bool {
	return func(b cassette.Interaction) bool {
		return a.Request.URL == b.Request.URL && bytes.Equal(compact(a.Request.Body), compact(b.Request.Body))
	}
}

// compact is body without the whitespace the file of interactions is indented with.
func compact(body cassette.Body) []byte {
	v := jsontext.Value(bytes.Clone(body))
	if v.Compact() != nil {
		return body
	}

	return v
}

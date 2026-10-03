package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"slices"
	"strings"
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

	for _, r := range []cassette.Request{
		{Method: http.MethodGet, URL: "https://api.notion.com/v1/views?database_id=" + dbID.String()},
		// TODO: add more here
	} {
		var reqBody io.Reader
		if len(r.Body) > 0 {
			reqBody = bytes.NewReader(r.Body)
		}

		req, err := http.NewRequestWithContext(ctx, r.Method, r.URL, reqBody)
		if err != nil {
			return err
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
			return err
		}

		ia := cassette.Interaction{Request: r}

		rsp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		ia.Response, err = cassette.NewResponse(rsp)
		if err != nil {
			return fmt.Errorf("recording response: %w", err)
		}

		ias = slices.DeleteFunc(ias, eqTo(ia))
		ias = append(ias, ia)
	}

	ias.Mask()
	ias.TrimResponseHeaders()
	ias.TrimResponseBodies(3)

	return ias.WriteFile(pathInteractions)
}

func cmp(a, b cassette.Interaction) int { return strings.Compare(a.Request.URL, b.Request.URL) }
func eqTo(a cassette.Interaction) func(cassette.Interaction) bool {
	return func(b cassette.Interaction) bool { return cmp(a, b) == 0 }
}

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"uuid"

	"github.com/MarkRosemaker/openapi-enrich/cassette"
)

var (
	// https://app.notion.com/p/fae-tools/457d00dcfc3a4ddba42ffe376530bee0?v=cb3265bc5d5e488fbd8a69221c107915
	dbID         = uuid.MustParse("457d00dcfc3a4ddba42ffe376530bee0")
	dataSourceID = uuid.MustParse("9201f6db-2895-4450-9435-041541a181dd")
	viewID       = uuid.MustParse("cb3265bc5d5e488fbd8a69221c107915")
)

func main() {
	if err := run(context.Background()); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Done!")
}

var reqBody io.Reader

func run(ctx context.Context) error {
	url := "https://api.notion.com/v1/views?database_id=" + dbID.String()

	req, err := http.NewRequestWithContext(ctx, "GET", url, reqBody)
	if err != nil {
		return err
	}

	apiKey := os.Getenv("NOTION_API_KEY")
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Notion-Version", "2026-03-11")

	if reqBody != nil {
		req.Header.Add("Content-Type", "application/json")
	}

	r, err := cassette.NewRequest(req)
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

	if err := cassette.AddInteraction("api/interactions.json", ia); err != nil {
		return err
	}

	return nil
}

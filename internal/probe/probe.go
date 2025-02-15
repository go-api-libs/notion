package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/go-api-libs/notion/pkg/notion"
	"github.com/google/uuid"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
)

var examplePageID = uuid.MustParse("96245c8f178444a482ad1941127c3ec3")

// u := baseURL.JoinPath(fmt.Sprintf("/blocks/%s/children", id))

// probe calls the API server to check what we can do
func probe() error {
	ctx := context.Background()

	c, err := notion.NewClient()
	if err != nil {
		return err
	}

	p, err := c.GetPage(ctx, examplePageID)
	if err != nil {
		return fmt.Errorf("getting page %q: %w", examplePageID, err)
	}

	if want := examplePageID; p.ID != want {
		return fmt.Errorf("unexpected page ID: got %v, want %v", p.ID, want)
	}

	if want := "Example Page"; p.Title() != want {
		return fmt.Errorf("unexpected page title: got %q, want %q", p.Title(), want)
	}

	list, err := c.GetBlocks(ctx, examplePageID, &notion.GetBlocksParams{
		PageSize: 10,
	})
	if err != nil {
		return err
	}

	fmt.Printf("num results: %v\n", len(list.Results))

	bearer := os.Getenv("NOTION_API_KEY")
	if bearer == "" {
		return fmt.Errorf("missing bearer token, set NOTION_API_KEY")
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		bearer = "Bearer " + bearer
	}

	reqEditor := func(req *http.Request) {
		req.Header.Set("Authorization", bearer)
		// Specifies the Notion API version
		req.Header.Set("Notion-Version", "2022-06-28")
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.notion.com/v1/blocks/"+examplePageID.String()+"/children?page_size=10", nil)
	if err != nil {
		return err
	}

	reqEditor(req)

	_, err = http.DefaultClient.Do(req)
	return err
}

// mask any secrets the API might return, e.g. in the response body
func maskSecrets(i *cassette.Interaction) error {
	return nil
}

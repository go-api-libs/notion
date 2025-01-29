package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
)

var examplePageID = uuid.MustParse("96245c8f178444a482ad1941127c3ec3")

// u := baseURL.JoinPath(fmt.Sprintf("/blocks/%s/children", id))

// probe calls the API server to check what we can do
func probe() error {
	bearer := os.Getenv("NOTION_TOKEN")
	if bearer == "" {
		return fmt.Errorf("missing bearer token, set NOTION_TOKEN")
	}

	if !strings.HasPrefix(bearer, "Bearer ") {
		bearer = "Bearer " + bearer
	}

	reqEditor := func(req *http.Request) {
		req.Header.Set("Authorization", bearer)
		// Specifies the Notion API version
		req.Header.Set("Notion-Version", "2022-06-28")
	}

	req, err := http.NewRequest(http.MethodGet, "https://api.notion.com/v1/pages/"+examplePageID.String(), nil)
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

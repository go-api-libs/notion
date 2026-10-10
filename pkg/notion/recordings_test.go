package notion

import (
	"encoding/json/v2"
	"net/http"
	"strings"
	"testing"

	"github.com/MarkRosemaker/openapi-enrich/cassette"
)

// TestRecordedResponses decodes every successful response testing/main.go recorded into the type its endpoint
// returns, so that a response Notion sends but the generated types cannot hold is noticed.
func TestRecordedResponses(t *testing.T) {
	ias, err := cassette.InteractionsReadFile("../../api/interactions.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, ia := range ias {
		if ia.Response.StatusCode != http.StatusOK {
			continue
		}

		v := responseOf(ia.Request.URL)
		if v == nil {
			continue
		}

		t.Run(ia.Request.Method+" "+ia.Request.URL, func(t *testing.T) {
			if err := json.Unmarshal(ia.Response.Body, v); err != nil {
				t.Error(err)
			}
		})
	}
}

// responseOf is a new value of the type the endpoint at url returns, or nil for one not covered.
func responseOf(url string) any {
	switch path := strings.TrimPrefix(url, "https://api.notion.com/v1/"); {
	case strings.HasSuffix(path, "/templates"):
		return new(ListDataSourceTemplatesOk)
	case strings.HasPrefix(path, "comments"):
		return new(CommentList)
	case strings.HasSuffix(path, "/query"), path == "search":
		return new(PageOrDataSourceList)
	case strings.HasPrefix(path, "blocks/") && strings.Contains(path, "/children"):
		return new(BlockList)
	case strings.HasPrefix(path, "blocks/"):
		return new(BlockOrPartial)
	case strings.Contains(path, "/properties/"):
		return new(RetrieveAPagePropertyOk)
	case strings.HasPrefix(path, "pages/"):
		return new(PageOrPartial)
	case strings.HasPrefix(path, "databases/"):
		return new(DatabaseOrPartial)
	case strings.HasPrefix(path, "data_sources/"):
		return new(DataSourceOrPartial)
	case strings.HasPrefix(path, "views/"):
		return new(DataSourceViewOrPartial)
	case strings.HasPrefix(path, "views?"):
		return new(ViewList)
	case path == "users/me":
		return new(UserObjectResponse)
	case strings.HasPrefix(path, "users?"):
		return new(UserList)
	default:
		return nil
	}
}

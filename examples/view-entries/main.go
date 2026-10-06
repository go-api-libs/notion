// Command view-entries prints the entries of a Notion database as one of its views shows them: filtered and sorted
// like the view.
//
// It reads the view once, for its data source, filter and sorts, and then needs one request per 100 entries each time
// it fetches them. Run it with NOTION_API_TOKEN set to an integration's token that can read the database.
package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-api-libs/notion/pkg/notion"
)

// viewID is the board view of https://app.notion.com/p/fae-tools/457d00dcfc3a4ddba42ffe376530bee0?v=cb3265bc5d5e488fbd8a69221c107915.
const viewID notion.IDRequest = "cb3265bc-5d5e-488f-bd8a-69221c107915"

func main() {
	ctx := context.Background()

	c, err := notion.NewClient() // reads NOTION_API_TOKEN
	if err != nil {
		log.Fatal(err)
	}

	dataSourceID, query, err := queryOfView(ctx, c, viewID)
	if err != nil {
		log.Fatal(err)
	}

	// each time the entries are needed: pass &notion.PostDatabaseQueryParams{FilterProperties: ids} instead of nil to
	// have Notion return only some properties
	for result, err := range c.PostDatabaseQueryAll(ctx, dataSourceID, nil, query) {
		if err != nil {
			log.Fatal(err)
		}

		if page := result.Page(); page != nil {
			fmt.Println(page.ID, page.Title())
		}
	}
}

// queryOfView returns the data source the view shows, and the query of it that filters and sorts its entries like
// the view. It stays the same until the view's filter or sorts are edited, so read it once: Notion evaluates relative
// dates such as "today" with each query.
func queryOfView(ctx context.Context, c *notion.Client, id notion.IDRequest) (notion.IDRequest, notion.PostDatabaseQuery, error) {
	v, err := c.RetrieveAView(ctx, id)
	if err != nil {
		return "", notion.PostDatabaseQuery{}, err
	}

	view := v.DataSourceViewObjectResponse
	if view == nil {
		return "", notion.PostDatabaseQuery{}, errors.New("no access to the view")
	}

	query := notion.PostDatabaseQuery{Sorts: view.Sorts}
	if view.Filter != nil {
		query.Filter = view.Filter.Flatten() // a view may nest its filter deeper than a query allows
	}

	return notion.IDRequest(view.DataSourceID), query, nil
}

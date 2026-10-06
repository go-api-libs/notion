package notion

import (
	"context"
	"iter"
	"uuid"
)

// maxPageSize is the most items Notion returns in one page of a list.
const maxPageSize = 100

// PostDatabaseQueryAll yields every entry of the data source the query selects, in its order, fetching them a page of
// 100 at a time. It stops at the first error, which it yields.
func (c *Client) PostDatabaseQueryAll(
	ctx context.Context, dataSourceID IDRequest, params *PostDatabaseQueryParams, body PostDatabaseQuery,
) iter.Seq2[PageOrDataSource, error] {
	if body.PageSize == nil {
		body.PageSize = new(float64(maxPageSize))
	}

	return all(func(cursor string) ([]PageOrDataSource, string, bool, error) {
		body.StartCursor = cursor

		list, err := c.PostDatabaseQuery(ctx, dataSourceID, params, body)
		if err != nil {
			return nil, "", false, err
		}

		return list.Results, list.NextCursor, list.HasMore, nil
	})
}

// PostSearchAll yields every page and data source the search finds, fetching them a page of 100 at a time. It stops
// at the first error, which it yields.
func (c *Client) PostSearchAll(ctx context.Context, body PostSearch) iter.Seq2[PageOrDataSource, error] {
	if body.PageSize == nil {
		body.PageSize = new(float64(maxPageSize))
	}

	return all(func(cursor string) ([]PageOrDataSource, string, bool, error) {
		start, err := startCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}

		body.StartCursor = start

		list, err := c.PostSearch(ctx, body)
		if err != nil {
			return nil, "", false, err
		}

		return list.Results, list.NextCursor, list.HasMore, nil
	})
}

// GetBlockChildrenAll yields every block within the block or page, in its order, fetching them a page of 100 at a
// time. It stops at the first error, which it yields.
func (c *Client) GetBlockChildrenAll(ctx context.Context, blockID IDRequest) iter.Seq2[BlockOrPartial, error] {
	return all(func(cursor string) ([]BlockOrPartial, string, bool, error) {
		start, err := startCursor(cursor)
		if err != nil {
			return nil, "", false, err
		}

		list, err := c.GetBlockChildren(ctx, blockID, &GetBlockChildrenParams{StartCursor: start, PageSize: maxPageSize})
		if err != nil {
			return nil, "", false, err
		}

		return list.Results, list.NextCursor, list.HasMore, nil
	})
}

// all yields the items of each page of a list, which page returns with the cursor of the next and whether there is
// one, given the cursor of its own, empty for the first.
func all[T any](page func(cursor string) (items []T, next string, more bool, err error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		for cursor := ""; ; {
			items, next, more, err := page(cursor)
			if err != nil {
				var zero T

				yield(zero, err)

				return
			}

			for _, item := range items {
				if !yield(item, nil) {
					return
				}
			}

			if !more {
				return
			}

			cursor = next
		}
	}
}

// startCursor is cursor as the endpoints that take it as a UUID do, the nil UUID for the first page.
func startCursor(cursor string) (uuid.UUID, error) {
	if cursor == "" {
		return uuid.Nil(), nil
	}

	return uuid.Parse(cursor)
}

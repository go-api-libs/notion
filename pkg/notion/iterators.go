package notion

import (
	"context"
	"iter"
)

// maxPageSize is the most items Notion returns in one page of a list.
const maxPageSize = 100

// PostDatabaseQueryAll yields every entry of the data source the query selects, in its order, fetching them a page of
// 100 at a time. It stops at the first error, which it yields.
func (c *Client) PostDatabaseQueryAll(
	ctx context.Context, dataSourceID ID, params *PostDatabaseQueryParams, body PostDatabaseQuery,
) iter.Seq2[PageOrDataSource, error] {
	if body.PageSize == 0 {
		body.PageSize = maxPageSize
	}

	return all(func(cursor QueryCursor) ([]PageOrDataSource, QueryCursor, bool, error) {
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
func (c *Client) PostSearchAll(ctx context.Context, params *PostSearchParams, body PostSearch) iter.Seq2[PageOrDataSource, error] {
	if body.PageSize == nil {
		body.PageSize = new(float64(maxPageSize))
	}

	return all(func(cursor ID) ([]PageOrDataSource, ID, bool, error) {
		body.StartCursor = cursor

		list, err := c.PostSearch(ctx, params, body)
		if err != nil {
			return nil, ID{}, false, err
		}

		return list.Results, list.NextCursor, list.HasMore, nil
	})
}

// GetBlockChildrenAll yields every block within the block or page, in its order, fetching them a page of 100 at a
// time. It stops at the first error, which it yields.
func (c *Client) GetBlockChildrenAll(ctx context.Context, blockID ID) iter.Seq2[BlockOrPartial, error] {
	return all(func(cursor ID) ([]BlockOrPartial, ID, bool, error) {
		list, err := c.GetBlockChildren(ctx, blockID, &GetBlockChildrenParams{StartCursor: cursor, PageSize: maxPageSize})
		if err != nil {
			return nil, ID{}, false, err
		}

		return list.Results, list.NextCursor, list.HasMore, nil
	})
}

// all yields the items of each page of a list, which page returns with the cursor of the next and whether there is
// one, given the cursor of its own, the zero cursor for the first.
func all[T any, C comparable](page func(cursor C) (items []T, next C, more bool, err error)) iter.Seq2[T, error] {
	return func(yield func(T, error) bool) {
		var cursor C

		for {
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

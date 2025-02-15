package notion

import (
	"context"
	"fmt"
	"iter"

	"github.com/google/uuid"
)

// BlocksIter is an iterator for blocks in a Notion page.
type BlocksIter struct {
	c    *Client   // Client to interact with the Notion API
	id   uuid.UUID // ID of the Notion page
	next uuid.UUID // Cursor for the next page of results
	err  error     // Error encountered during iteration
}

// ListBlocks creates a new BlocksIter for the given page ID.
func (c *Client) ListBlocks(id uuid.UUID) BlocksIter {
	return BlocksIter{c: c, id: id}
}

const maxPageSizeInt = 100 // Maximum number of blocks per page

// All returns an iterator that yields all blocks in the Notion page.
func (it *BlocksIter) All(ctx context.Context) iter.Seq2[int, Block] {
	it.err = nil // Reset the error

	i := 0
	return func(yield func(int, Block) bool) {
		// Reset the next cursor when the iterator is done
		// so that the iterator can be reused
		defer func() { it.next = uuid.Nil }()

		for page := 0; ; page++ {
			// Get a page of blocks from the Notion API
			list, err := it.c.GetBlocks(ctx, it.id, &GetBlocksParams{
				PageSize:    maxPageSizeInt,
				StartCursor: it.next,
			})
			if err != nil {
				it.err = fmt.Errorf("page %d of getting blocks for %s: %w", page, it.id, err)
				return
			}

			// Yield each block to the caller
			for _, block := range list.Results {
				if !yield(i, block) {
					return
				}

				i++
			}

			// If there are no more blocks, stop iterating
			if !list.HasMore {
				return
			}

			// Update the cursor for the next page of results
			it.next = list.NextCursor
		}
	}
}

// Err returns the error encountered during iteration, if any.
// It must be inspected after All returns.
func (it *BlocksIter) Err() error { return it.err }

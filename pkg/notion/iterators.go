package notion

// import (
// 	"context"
// 	"fmt"
// 	"iter"

// 	"github.com/google/uuid"
// )

// type BlocksIterator struct {
// 	cli  *Client
// 	id   uuid.UUID
// 	next *string
// 	err  error
// }

// func (c *Client) BlocksIterator(id uuid.UUID) BlocksIterator {
// 	return BlocksIterator{cli: c, id: id}
// }

// func (it *BlocksIterator) All(ctx context.Context) iter.Seq2[int, Block] {
// 	it.err = nil // reset the error

// 	i := 0
// 	return func(yield func(int, Block) bool) {
// 		// reset the next cursor when the iterator is done
// 		// so that the iterator can be reused
// 		defer func() { it.next = nil }()

// 		for page := 0; ; page++ {
// 			list, err := it.cli.GetBlocks(ctx, it.id, &GetBlocksParams{
// 				PageSize:    &maxPageSizeInt,
// 				StartCursor: it.next,
// 			})
// 			if err != nil {
// 				it.err = fmt.Errorf("page %d of getting blocks for %s: %w", page, it.id, err)
// 				return
// 			}

// 			for _, block := range list.Results {
// 				if !yield(i, block) {
// 					return
// 				}

// 				i++
// 			}

// 			if !list.HasMore {
// 				return
// 			}

// 			s := list.NextCursor.String()
// 			it.next = &s
// 		}
// 	}
// }

// func (it *BlocksIterator) Err() error { return it.err }

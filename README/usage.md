To install the library, use the following command:

```shell
go get github.com/go-api-libs/notion/pkg/notion
```


### Example 1: Retrieve a Page

```go
package main

import (
	"context"

	"github.com/go-api-libs/notion/pkg/notion"
	"github.com/google/uuid"
)

func main() {
	c, err := notion.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	page, err := c.GetPage(ctx, uuid.MustParse("96245c8f-1784-44a4-82ad-1941127c3ec3"))
	if err != nil {
		panic(err)
	}

	// Use page object
}

```

### Example 2: Retrieve block children

```go
package main

import (
	"context"

	"github.com/go-api-libs/notion/pkg/notion"
	"github.com/google/uuid"
)

func main() {
	c, err := notion.NewClient()
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	blocksList, err := c.GetBlocks(ctx, uuid.MustParse("96245c8f-1784-44a4-82ad-1941127c3ec3"), &notion.GetBlocksParams{
		PageSize:    100,
		StartCursor: uuid.MustParse("d05aa478-397e-4954-b694-b7c1c6c78956"),
	})
	if err != nil {
		panic(err)
	}

	// Use blocksList object
}

```

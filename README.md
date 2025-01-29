# Notion API
[![Go Reference](https://pkg.go.dev/badge/github.com/go-api-libs/notion.svg)](https://pkg.go.dev/github.com/go-api-libs/notion/pkg/notion)
[![OpenAPI](https://img.shields.io/badge/OpenAPI-3.1-blue)](/api/openapi.json)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-api-libs/notion)](https://goreportcard.com/report/github.com/go-api-libs/notion)
![Code Coverage](https://img.shields.io/badge/coverage-22%25-red)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](./LICENSE)

[Create an integration](https://www.notion.so/my-integrations) to retrieve an API token, add your database and page ID's as variables in the collection, and start making your requests!

For our full documentation, including sample integrations and guides, visit [developers.notion.com](developers.notion.com)

Need more help? Join our [developer community on Slack](https://join.slack.com/t/notiondevs/shared_invite/zt-lkrnk74h-YmPRroySRFGiqgjI193AqA/)

## Installation

To install the library, use the following command:

```shell
go get github.com/go-api-libs/notion/pkg/notion
```

## Usage

### Example: 

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

## Additional Information

- [**Go Reference**](https://pkg.go.dev/github.com/go-api-libs/notion/pkg/notion): The Go reference documentation for the client package.
- [**OpenAPI Specification**](./api/openapi.json): The OpenAPI 3.1.0 specification.
- [**Go Report Card**](https://goreportcard.com/report/github.com/go-api-libs/notion): Check the code quality report.

## Contributing

If you have any contributions to make, please submit a pull request or open an issue on the [GitHub repository](https://github.com/go-api-libs/notion).

## License

This project is licensed under the MIT License. See the [LICENSE](./LICENSE) file for details.

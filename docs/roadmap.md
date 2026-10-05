# Roadmap

## Test coverage of the example page

`testing/main.go` records the [example page](https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3),
its subpages and databases. Recording these too needs them added there first:

- Blocks: heading 4, tabs, AI meeting notes. The template block is deprecated and can no longer be created.
- Database properties: last edited time, ID (unique ID), button, place, and verification, which only a wiki has.
- Page icons: a custom emoji, and one of Notion's own icons.
- Inline mentions: a link pasted as a mention for a service such as GitHub or Figma (link preview), a custom emoji,
  and a template mention such as `@today` in a database template.
- Rollups: one whose result is a date.
- Comments: the integration answers every comments request with 403. It needs the "Read comments" capability, and
  the page a comment.

## Recordings the generated types cannot decode

Five of the recorded responses in `api/interactions.json` fail to decode into the generated types:

- `POST /v1/search`: its results mix pages and a data source, which enrich types as a three-item tuple next to the
  array, so `PageOrDataSourceList.Results` is a union instead of a slice.
- `POST /v1/data_sources/d58cc2b8-…/query`: `PageOrDataSourceListResults` matches none of its alternatives.
- `GET /v1/views/cb3265bc-…`: a board view matches no alternative of `DataSourceViewOrPartial`.
- `GET /v1/pages/24806928-…/properties/%3DuTh`: a formula property item matches no alternative.
- `GET /v1/blocks/7d4dc32d-…`: a paragraph that mentions a user matches no alternative of `BlockOrPartial`.

A test that decodes every recorded response with the generated types would catch the next one.

Enrich also still files a recorded `and` filter under the query filter's `or` branch, whose `or` it requires, and
infers a new item schema for it (`GroupFilterOperatorArrayItemOrAnd…`).

## Names in the specification

- 40 names end in a number to resolve a clash, such as `PageIcon2` or `QueryAgentsFilterStatus2`. Compress appends
  one when a shortened name is taken, and flatten when a name is (`AgentID6`); keeping the longer name would say more.
- Components that only make another nullable, such as `PageIcon` or `DateTimeZone`, need no name of their own: codegen
  makes them pointers. Flatten could leave them inline, and name their branch after them.
- 175 names hold a position (`…OneOf2…`, `…AllOf1…`), mostly union branches that no const, one-valued enum or single
  required member tells apart. Those need names picked by hand in `api/fetch.go`.
- The `ObjectResponse` and `Response` suffixes could go, and the names compress leaves in lower case, such as
  `numberSimplePropertyValue`, could be in Go's PascalCase like the rest.

## Library

- Iterators over paginated results: `pkg/notion/iterators.go` holds a commented-out one for blocks to build them from.
- The entries of a view, filtered and sorted as it shows them: read the view once, then query its data source with its
  filter and sorts. A view nests filters deeper than a query may, so an `or` within an `or`, or an `and` within an
  `and`, has to join the outer one. `testing/main.go` does this for its recording.

## Recordings

- `api/interactions.json` keeps the page that was once the first entry of the board's data source (`3ee88fcd-…`),
  which `testing/main.go` no longer records or refreshes.

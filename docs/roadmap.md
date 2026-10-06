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

One of the recorded responses in `api/interactions.json` fails to decode into the generated types:
`GET /v1/pages/24806928-…/properties/%3DuTh`, a formula property item, matches no alternative of
`RetrieveAPagePropertyOk`.

A test that decodes every recorded response with the generated types would catch the next one.

## Workarounds in `api/fetch.go` for the openapi-* libraries

Each can go once the library is fixed:

- `setFilterForms` and `setResults` set `Filter` and `PageOrDataSourceList.Results` again after enrich. Enrich files a
  recorded value under a union's first form even when it lacks that form's required members, and types an array
  whose items differ in shape as a tuple.
- `groupByConfigResponse` is made an `anyOf`, as codegen does not check a form's const members, so a `group_by`
  matched all its forms.
- Codegen does not check a form's required members either, which only `additionalProperties: false` makes up for when
  decoding strictly. Decoding into a caller's own type with `…WithResult` is lenient, so a `Filter` within one
  decodes `{"and": […]}` as an empty `or`.

## Names in the specification

- 40 names end in a number to resolve a clash, such as `PageIcon2` or `QueryAgentsFilterStatus2`. Compress appends
  one when a shortened name is taken, and flatten when a name is (`AgentID6`); keeping the longer name would say more.
- Components that only make another nullable, such as `PageIcon` or `DateTimeZone`, need no name of their own: codegen
  makes them pointers. Flatten could leave them inline, and name their branch after them.
- 175 names hold a position (`…OneOf2…`, `…AllOf1…`), mostly union branches that no const, one-valued enum or single
  required member tells apart. Those need names picked by hand in `api/fetch.go`.
- The `ObjectResponse` and `Response` suffixes could go, and the names compress leaves in lower case, such as
  `numberSimplePropertyValue`, could be in Go's PascalCase like the rest.

## Recordings

- `api/interactions.json` keeps the page that was once the first entry of the board's data source (`3ee88fcd-…`),
  which `testing/main.go` no longer records or refreshes.

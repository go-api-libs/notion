# Roadmap

## Test coverage of the example page

`testing/main.go` records the [example page](https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3),
its subpages, databases, their templates and the page's comments. It holds every kind of block, icon, mention and
property the API returns, but these:

- Database properties: button and location, which the API accepts but drops, and verification, which only a wiki has.
- A template mention such as `@today`, which only a database template may hold: none of the example page's databases
  has a template.

## Recordings the generated types cannot decode

One of the recorded responses in `api/interactions.json` fails to decode into the generated types:
`GET /v1/pages/24806928-…/properties/%3DuTh`, a formula of `{"type": "number", "number": null}`. Codegen takes a
required member that is `null` for a missing one, as its pointer stays nil.

A test that decodes every recorded response with the generated types would catch the next one.

## Workarounds in `api/fetch.go` for the openapi-* libraries

Each can go once the library is fixed:

- `setFilterForms` and `setResults` set `Filter` and `PageOrDataSourceList.Results` again after enrich. Enrich files a
  recorded value under a union's first form even when it lacks that form's required members, and types an array
  whose items differ in shape as a tuple.
- `allowDateTimes` and `setDateTimes` make a date's `start` and `end` a plain string for enrich, and a
  `DateOrDateTime` after it. Enrich infers no `date` format for a value such as `2026-10-05`, so it matches neither
  form of an `anyOf` of a date and a date-time.
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

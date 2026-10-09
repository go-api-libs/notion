# Roadmap

## Test coverage of the example page

`testing/main.go` records the [example page](https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3),
its subpages, databases, their templates and the page's comments. It holds every kind of block, icon, mention and
property the API returns, but these:

- Database properties: button and location, which the API accepts but drops, and verification, which only a wiki has.
- A template mention such as `@today`, which only a database template may hold: none of the example page's databases
  has a template.

## Workarounds in `api/fetch.go` for the openapi-* libraries

- `setFilterForms` sets `Filter` again after enrich, which does not look into a form that is itself a union, such as
  `PropertyFilter`: a recorded property filter lands under `FilterOr`, the first form. Enrich also drops the
  `maxItems` of an array it was given a value of. It can go once enrich matches a value to a form within such a form.

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

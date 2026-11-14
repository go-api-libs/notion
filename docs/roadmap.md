# Roadmap

## Use `github.com/MarkRosemaker/cassette`'s Framework

Use cassette's implementation fully instead of home-grown recorder.

## Test coverage of the example page

`testing/main.go` records the [example page](https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3),
its subpages, databases, their templates and the page's comments. It holds every kind of block, icon, mention and
property the API returns, but these:

- Database properties: button and location, which the API accepts but drops, and verification, which only a wiki has.
- A template mention such as `@today`, which only a database template may hold: none of the example page's databases
  has a template.

## Names in the specification

- 43 names say a word twice, as `UpdateABlockCodeCode` or `MarkdownEditInsertContentInsertContent`: what a block or
  an edit holds, named after it, where saying the word once would be the name of what holds it. Notion calls them by
  the word alone (`Code`), which the request, response and update forms of each cannot all have.

## Recordings

- `api/interactions.json` keeps the page that was once the first entry of the board's data source (`3ee88fcd-…`),
  which `testing/main.go` no longer records or refreshes.

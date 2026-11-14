# Roadmap

## Test coverage of the example page

`testing/main.go` records the [example page](https://app.notion.com/p/fae-tools/Example-Page-96245c8f178444a482ad1941127c3ec3),
its subpages, databases, their templates and the page's comments. It holds every kind of block, icon, mention and
property the API returns, but these:

- Database properties: button and location, which the API accepts but drops, and verification, which only a wiki has.
- A template mention such as `@today`, which only a database template may hold: none of the example page's databases
  has a template.

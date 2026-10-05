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

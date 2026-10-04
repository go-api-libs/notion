package notion

import (
	"encoding/json/v2"
	"testing"
)

// TestBlock decodes a block as Notion sent it, object first.
func TestBlock(t *testing.T) {
	b := &Block{}
	if err := json.Unmarshal([]byte(`{
		"object": "block",
		"id": "457d00dc-fc3a-4ddb-a42f-fe376530bee0",
		"parent": {"type": "page_id", "page_id": "11788fcd-daf4-80e7-a530-c2e25a86328e"},
		"created_time": "2024-03-19T08:28:00.000Z",
		"last_edited_time": "2026-10-03T15:37:00.000Z",
		"created_by": {"object": "user", "id": "af171d5d-c36f-45bc-a0a3-6086c0dafa45"},
		"last_edited_by": {"object": "user", "id": "af171d5d-c36f-45bc-a0a3-6086c0dafa45"},
		"has_children": false,
		"in_trash": false,
		"type": "child_database",
		"child_database": {"title": "Intentions"},
		"request_id": "466c960f-ce62-4648-9b60-5dce07df1dcb"
	}`), b, jsonOpts); err != nil {
		t.Fatal(err)
	}

	if b.Type != BlockTypeChildDatabase {
		t.Errorf("type is %q, want %q", b.Type, BlockTypeChildDatabase)
	}

	if got, want := b.ChildDatabase.Title, "Intentions"; got != want {
		t.Errorf("child database title is %q, want %q", got, want)
	}
}

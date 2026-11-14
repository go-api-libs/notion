package notion

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"
)

// TestCreateAComment sends each kind of comment Notion documents, as a client builds it: on a page or as a reply in
// a discussion, in rich text or in Markdown, signed with a name of its own and with a file attached. The client says
// the type of a parent and of rich text itself.
func TestCreateAComment(t *testing.T) {
	page := ID(uuid.MustParse("5c6a2821-6bb1-4a7e-b6e1-c50111515c3d"))
	discussion := ID(uuid.MustParse("59e3eb41-36a9-4e3a-8ea1-1a8b3b4d0a0e"))
	hello := []RichTextItem{{Text: &TextRichTextItemText{Content: "Hello world"}}}

	for _, tc := range []struct {
		name string
		body CreateAComment
		want string
	}{
		{
			"on a page",
			CreateAComment{NewComment: NewComment{PageComment: &PageComment{
				Parent: CreateACommentAllOfOneOfParent{PageID: &page}, RichText: hello,
			}}},
			`{"parent":{"type":"page_id","page_id":"5c6a2821-6bb1-4a7e-b6e1-c50111515c3d"},"rich_text":[{"type":"text","text":{"content":"Hello world"}}]}`,
		},
		{
			"on a page in Markdown",
			CreateAComment{MarkdownVersion: MarkdownVersionV2, NewComment: NewComment{PageMarkdownComment: &PageMarkdownComment{
				Parent: CreateACommentAllOfOneOfParent{PageID: &page}, Markdown: "**Hello** world",
			}}},
			`{"parent":{"type":"page_id","page_id":"5c6a2821-6bb1-4a7e-b6e1-c50111515c3d"},"markdown":"**Hello** world","markdown_version":"v2"}`,
		},
		{
			"a reply",
			CreateAComment{NewComment: NewComment{DiscussionComment: &DiscussionComment{
				DiscussionID: discussion, RichText: hello,
			}}},
			`{"discussion_id":"59e3eb41-36a9-4e3a-8ea1-1a8b3b4d0a0e","rich_text":[{"type":"text","text":{"content":"Hello world"}}]}`,
		},
		{
			"a reply in Markdown, signed and with a file",
			CreateAComment{
				DisplayName: CreateACommentDisplayName{Type: CreateACommentDisplayNameTypeCustom, Custom: &CustomDisplayName{Name: "Bot"}},
				Attachments: CreateACommentAttachments{{FileUploadID: "2e2cdb8b-9897-4a6c-a935-82922b1cfb87"}},
				NewComment: NewComment{DiscussionMarkdownComment: &DiscussionMarkdownComment{
					DiscussionID: discussion, Markdown: "Done",
				}},
			},
			`{"discussion_id":"59e3eb41-36a9-4e3a-8ea1-1a8b3b4d0a0e","markdown":"Done",` +
				`"display_name":{"type":"custom","custom":{"name":"Bot"}},` +
				`"attachments":[{"file_upload_id":"2e2cdb8b-9897-4a6c-a935-82922b1cfb87"}]}`,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/comments" {
					t.Errorf("sent %s %s, want POST /comments", r.Method, r.URL.Path)
				}

				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}

				if !sameJSON(t, tc.want, string(body)) {
					t.Errorf("sent %s\n want %s", body, tc.want)
				}

				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"object":"comment","id":"249911a1-3cbd-4ee2-9b0a-dc2bc4d2a8fb"}`))
			}))
			defer srv.Close()

			if _, err := testClient(t, srv).CreateAComment(t.Context(), tc.body); err != nil {
				t.Fatal(err)
			}
		})
	}
}

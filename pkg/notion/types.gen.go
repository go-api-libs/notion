// This file provides types for the API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.

package notion

import (
	"net/url"
	"time"

	"github.com/go-api-libs/types"
	"github.com/google/uuid"
)

type GetBlocksParams struct {
	PageSize int
}

// Style information which applies to the whole rich text object.
type Annotations struct {
	// Whether the text is **bolded**.
	Bold bool `json:"bold,omitzero"`
	// Whether the text is *italicized*.
	Italic bool `json:"italic,omitzero"`
	// Whether the text is struck through.
	Strikethrough bool `json:"strikethrough,omitzero"`
	// Whether the text is underlined.
	Underline bool `json:"underline,omitzero"`
	// Whether the text is code `style`.
	Code bool `json:"code,omitzero"`
	// The color of the block.
	Color Color `json:"color,omitzero"`
}

// A block object represents content within Notion. Blocks can be text, lists, media, and more. A page is a type of block, too!
//
// The optional fields are filled depending on the value of `type`.
type Block struct {
	// Always "block".
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
	// The `parent` property of a page or database contains these keys. Mandatory when creating, must be missing when updating.
	Parent Parent `json:"parent"`
	// Date and time when this block was created. Formatted as an ISO 8601 date time string.
	CreatedTime time.Time `json:"created_time,omitzero"`
	// Date and time when this block was last updated. Formatted as an ISO 8601 date time string.
	LastEditedTime   time.Time              `json:"last_edited_time,omitzero"`
	CreatedBy        UserReference          `json:"created_by"`
	LastEditedBy     UserReference          `json:"last_edited_by"`
	HasChildren      bool                   `json:"has_children,omitzero"`
	Archived         bool                   `json:"archived,omitzero"`
	InTrash          bool                   `json:"in_trash,omitzero"`
	Type             string                 `json:"type,omitzero"`
	Paragraph        BlockParagraph         `json:"paragraph"`
	Heading1         *BlockHeading1         `json:"heading_1,omitempty"`
	Heading2         *BlockHeading2         `json:"heading_2,omitempty"`
	Heading3         *BlockHeading3         `json:"heading_3,omitempty"`
	Callout          *BlockCallout          `json:"callout,omitempty"`
	Quote            *BlockQuote            `json:"quote,omitempty"`
	SyncedBlock      *BlockSyncedBlock      `json:"synced_block,omitempty"`
	NumberedListItem *BlockNumberedListItem `json:"numbered_list_item,omitempty"`
	BulletedListItem *BlockBulletedListItem `json:"bulleted_list_item,omitempty"`
	ToDo             *BlockToDo             `json:"to_do,omitempty"`
	Toggle           *BlockToggle           `json:"toggle,omitempty"`
	Code             *BlockCode             `json:"code,omitempty"`
	ChildPage        *BlockChildPage        `json:"child_page,omitempty"`
	ChildDatabase    *BlockChildDatabase    `json:"child_database,omitempty"`
	Embed            *BlockEmbed            `json:"embed,omitempty"`
	PDF              *BlockPdf              `json:"pdf,omitempty"`
	ColumnList       *struct{}              `json:"column_list,omitempty"`
	Video            *BlockVideo            `json:"video,omitempty"`
	File             *BlockFile             `json:"file,omitempty"`
	Bookmark         *BlockBookmark         `json:"bookmark,omitempty"`
	Equation         *BlockEquation         `json:"equation,omitempty"`
	Divider          *struct{}              `json:"divider,omitempty"`
	TableOfContents  *BlockTableOfContents  `json:"table_of_contents,omitempty"`
	Breadcrumb       *struct{}              `json:"breadcrumb,omitempty"`
	LinkPreview      *BlockLinkPreview      `json:"link_preview,omitempty"`
	Unsupported      *struct{}              `json:"unsupported,omitempty"`
	LinkToPage       *BlockLinkToPage       `json:"link_to_page,omitempty"`
	Table            *BlockTable            `json:"table,omitempty"`
	Audio            *BlockAudio            `json:"audio,omitempty"`
}

// BlockAudio defines a model
type BlockAudio struct {
	Caption  []struct{}         `json:"caption,omitempty"`
	Type     string             `json:"type,omitzero"`
	External BlockAudioExternal `json:"external"`
	File     *BlockAudioFile    `json:"file,omitempty"`
}

// BlockAudioExternal defines a model
type BlockAudioExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// BlockAudioFile defines a model
type BlockAudioFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// BlockBookmark defines a model
type BlockBookmark struct {
	Caption BlockBookmarkCaption `json:"caption,omitempty"`
	URL     url.URL              `json:"url,omitzero"`
}

// BlockBookmarkCaption defines a model
type BlockBookmarkCaption []BlockBookmarkCaptionItems

// BlockBookmarkCaptionItems defines a model
type BlockBookmarkCaptionItems struct {
	Type        string                               `json:"type,omitzero"`
	Text        BlockBookmarkCaptionItemsText        `json:"text"`
	Annotations BlockBookmarkCaptionItemsAnnotations `json:"annotations"`
	PlainText   string                               `json:"plain_text,omitzero"`
	Href        struct{}                             `json:"href"`
}

// BlockBookmarkCaptionItemsAnnotations defines a model
type BlockBookmarkCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockBookmarkCaptionItemsText defines a model
type BlockBookmarkCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockBulletedListItem defines a model
type BlockBulletedListItem struct {
	RichText BlockBulletedListItemRichText `json:"rich_text,omitempty"`
	Color    string                        `json:"color,omitzero"`
}

// BlockBulletedListItemRichText defines a model
type BlockBulletedListItemRichText []BlockBulletedListItemRichTextItems

// BlockBulletedListItemRichTextItems defines a model
type BlockBulletedListItemRichTextItems struct {
	Type        string                                        `json:"type,omitzero"`
	Text        BlockBulletedListItemRichTextItemsText        `json:"text"`
	Annotations BlockBulletedListItemRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                                        `json:"plain_text,omitzero"`
	Href        struct{}                                      `json:"href"`
}

// BlockBulletedListItemRichTextItemsAnnotations defines a model
type BlockBulletedListItemRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockBulletedListItemRichTextItemsText defines a model
type BlockBulletedListItemRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockCallout defines a model
type BlockCallout struct {
	RichText BlockCalloutRichText `json:"rich_text,omitempty"`
	Icon     BlockCalloutIcon     `json:"icon"`
	Color    string               `json:"color,omitzero"`
}

// BlockCalloutIcon defines a model
type BlockCalloutIcon struct {
	Type     string                    `json:"type,omitzero"`
	Emoji    string                    `json:"emoji,omitzero"`
	External *BlockCalloutIconExternal `json:"external,omitempty"`
	File     *BlockCalloutIconFile     `json:"file,omitempty"`
}

// BlockCalloutIconExternal defines a model
type BlockCalloutIconExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// BlockCalloutIconFile defines a model
type BlockCalloutIconFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// BlockCalloutRichText defines a model
type BlockCalloutRichText []BlockCalloutRichTextItems

// BlockCalloutRichTextItems defines a model
type BlockCalloutRichTextItems struct {
	Type        string                               `json:"type,omitzero"`
	Text        BlockCalloutRichTextItemsText        `json:"text"`
	Annotations BlockCalloutRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                               `json:"plain_text,omitzero"`
	Href        struct{}                             `json:"href"`
}

// BlockCalloutRichTextItemsAnnotations defines a model
type BlockCalloutRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockCalloutRichTextItemsText defines a model
type BlockCalloutRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockChildDatabase defines a model
type BlockChildDatabase struct {
	Title string `json:"title,omitzero"`
}

// BlockChildPage defines a model
type BlockChildPage struct {
	Title string `json:"title,omitzero"`
}

// BlockCode defines a model
type BlockCode struct {
	Caption  BlockCodeCaption  `json:"caption,omitempty"`
	RichText BlockCodeRichText `json:"rich_text,omitempty"`
	Language string            `json:"language,omitzero"`
}

// BlockCodeCaption defines a model
type BlockCodeCaption []BlockCodeCaptionItems

// BlockCodeCaptionItems defines a model
type BlockCodeCaptionItems struct {
	Type        string                           `json:"type,omitzero"`
	Text        BlockCodeCaptionItemsText        `json:"text"`
	Annotations BlockCodeCaptionItemsAnnotations `json:"annotations"`
	PlainText   string                           `json:"plain_text,omitzero"`
	Href        struct{}                         `json:"href"`
}

// BlockCodeCaptionItemsAnnotations defines a model
type BlockCodeCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockCodeCaptionItemsText defines a model
type BlockCodeCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockCodeRichText defines a model
type BlockCodeRichText []BlockCodeRichTextItems

// BlockCodeRichTextItems defines a model
type BlockCodeRichTextItems struct {
	Type        string                            `json:"type,omitzero"`
	Text        BlockCodeRichTextItemsText        `json:"text"`
	Annotations BlockCodeRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                            `json:"plain_text,omitzero"`
	Href        url.URL                           `json:"href,omitzero"`
}

// BlockCodeRichTextItemsAnnotations defines a model
type BlockCodeRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockCodeRichTextItemsText defines a model
type BlockCodeRichTextItemsText struct {
	Content string                         `json:"content,omitzero"`
	Link    BlockCodeRichTextItemsTextLink `json:"link"`
}

// BlockCodeRichTextItemsTextLink defines a model
type BlockCodeRichTextItemsTextLink struct {
	URL *url.URL `json:"url,omitempty"`
}

// BlockEmbed defines a model
type BlockEmbed struct {
	Caption BlockEmbedCaption `json:"caption,omitempty"`
	URL     url.URL           `json:"url,omitzero"`
}

// BlockEmbedCaption defines a model
type BlockEmbedCaption []BlockEmbedCaptionItems

// BlockEmbedCaptionItems defines a model
type BlockEmbedCaptionItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        *BlockEmbedCaptionItemsText        `json:"text,omitempty"`
	Annotations *BlockEmbedCaptionItemsAnnotations `json:"annotations,omitempty"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        *struct{}                          `json:"href,omitempty"`
}

// BlockEmbedCaptionItemsAnnotations defines a model
type BlockEmbedCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockEmbedCaptionItemsText defines a model
type BlockEmbedCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockEquation defines a model
type BlockEquation struct {
	Expression string `json:"expression,omitzero"`
}

// BlockFile defines a model
type BlockFile struct {
	Caption []struct{}    `json:"caption,omitempty"`
	Type    string        `json:"type,omitzero"`
	File    BlockFileFile `json:"file"`
	Name    string        `json:"name,omitzero"`
}

// BlockFileFile defines a model
type BlockFileFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// BlockHeading1 defines a model
type BlockHeading1 struct {
	RichText     BlockHeading1RichText `json:"rich_text,omitempty"`
	IsToggleable bool                  `json:"is_toggleable,omitzero"`
	Color        string                `json:"color,omitzero"`
}

// BlockHeading1RichText defines a model
type BlockHeading1RichText []BlockHeading1RichTextItems

// BlockHeading1RichTextItems defines a model
type BlockHeading1RichTextItems struct {
	Type        string                                `json:"type,omitzero"`
	Text        BlockHeading1RichTextItemsText        `json:"text"`
	Annotations BlockHeading1RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                                `json:"plain_text,omitzero"`
	Href        struct{}                              `json:"href"`
}

// BlockHeading1RichTextItemsAnnotations defines a model
type BlockHeading1RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockHeading1RichTextItemsText defines a model
type BlockHeading1RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockHeading2 defines a model
type BlockHeading2 struct {
	RichText     BlockHeading2RichText `json:"rich_text,omitempty"`
	IsToggleable bool                  `json:"is_toggleable,omitzero"`
	Color        string                `json:"color,omitzero"`
}

// BlockHeading2RichText defines a model
type BlockHeading2RichText []BlockHeading2RichTextItems

// BlockHeading2RichTextItems defines a model
type BlockHeading2RichTextItems struct {
	Type        string                                `json:"type,omitzero"`
	Text        BlockHeading2RichTextItemsText        `json:"text"`
	Annotations BlockHeading2RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                                `json:"plain_text,omitzero"`
	Href        struct{}                              `json:"href"`
}

// BlockHeading2RichTextItemsAnnotations defines a model
type BlockHeading2RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockHeading2RichTextItemsText defines a model
type BlockHeading2RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockHeading3 defines a model
type BlockHeading3 struct {
	RichText     BlockHeading3RichText `json:"rich_text,omitempty"`
	IsToggleable bool                  `json:"is_toggleable,omitzero"`
	Color        string                `json:"color,omitzero"`
}

// BlockHeading3RichText defines a model
type BlockHeading3RichText []BlockHeading3RichTextItems

// BlockHeading3RichTextItems defines a model
type BlockHeading3RichTextItems struct {
	Type        string                                `json:"type,omitzero"`
	Text        BlockHeading3RichTextItemsText        `json:"text"`
	Annotations BlockHeading3RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                                `json:"plain_text,omitzero"`
	Href        struct{}                              `json:"href"`
}

// BlockHeading3RichTextItemsAnnotations defines a model
type BlockHeading3RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockHeading3RichTextItemsText defines a model
type BlockHeading3RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockLinkPreview defines a model
type BlockLinkPreview struct {
	URL url.URL `json:"url,omitzero"`
}

// BlockLinkToPage defines a model
type BlockLinkToPage struct {
	Type   string    `json:"type,omitzero"`
	PageID uuid.UUID `json:"page_id,omitzero"`
}

// BlockNumberedListItem defines a model
type BlockNumberedListItem struct {
	RichText BlockNumberedListItemRichText `json:"rich_text,omitempty"`
	Color    string                        `json:"color,omitzero"`
}

// BlockNumberedListItemRichText defines a model
type BlockNumberedListItemRichText []BlockNumberedListItemRichTextItems

// BlockNumberedListItemRichTextItems defines a model
type BlockNumberedListItemRichTextItems struct {
	Type        string                                        `json:"type,omitzero"`
	Text        BlockNumberedListItemRichTextItemsText        `json:"text"`
	Annotations BlockNumberedListItemRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                                        `json:"plain_text,omitzero"`
	Href        struct{}                                      `json:"href"`
}

// BlockNumberedListItemRichTextItemsAnnotations defines a model
type BlockNumberedListItemRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockNumberedListItemRichTextItemsText defines a model
type BlockNumberedListItemRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockParagraph defines a model
type BlockParagraph struct {
	RichText BlockParagraphRichText `json:"rich_text,omitempty"`
	Color    string                 `json:"color,omitzero"`
}

// BlockParagraphRichText defines a model
type BlockParagraphRichText []BlockParagraphRichTextItems

// BlockParagraphRichTextItems defines a model
type BlockParagraphRichTextItems struct {
	Type        string                                  `json:"type,omitzero"`
	Text        *BlockParagraphRichTextItemsText        `json:"text,omitempty"`
	Annotations *BlockParagraphRichTextItemsAnnotations `json:"annotations,omitempty"`
	PlainText   string                                  `json:"plain_text,omitzero"`
	Href        *url.URL                                `json:"href,omitempty"`
	Mention     *BlockParagraphRichTextItemsMention     `json:"mention,omitempty"`
	Equation    *BlockParagraphRichTextItemsEquation    `json:"equation,omitempty"`
}

// BlockParagraphRichTextItemsAnnotations defines a model
type BlockParagraphRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockParagraphRichTextItemsEquation defines a model
type BlockParagraphRichTextItemsEquation struct {
	Expression string `json:"expression,omitzero"`
}

// BlockParagraphRichTextItemsMention defines a model
type BlockParagraphRichTextItemsMention struct {
	Type        string                                        `json:"type,omitzero"`
	LinkMention BlockParagraphRichTextItemsMentionLinkMention `json:"link_mention"`
	Database    *BlockParagraphRichTextItemsMentionDatabase   `json:"database,omitempty"`
	User        *BlockParagraphRichTextItemsMentionUser       `json:"user,omitempty"`
	Date        *BlockParagraphRichTextItemsMentionDate       `json:"date,omitempty"`
	Page        *BlockParagraphRichTextItemsMentionPage       `json:"page,omitempty"`
}

// BlockParagraphRichTextItemsMentionDatabase defines a model
type BlockParagraphRichTextItemsMentionDatabase struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// BlockParagraphRichTextItemsMentionDate defines a model
type BlockParagraphRichTextItemsMentionDate struct {
	Start    string   `json:"start,omitzero"`
	End      struct{} `json:"end"`
	TimeZone struct{} `json:"time_zone"`
}

// BlockParagraphRichTextItemsMentionLinkMention defines a model
type BlockParagraphRichTextItemsMentionLinkMention struct {
	Href        url.URL `json:"href,omitzero"`
	Title       string  `json:"title,omitzero"`
	Description string  `json:"description,omitzero"`
}

// BlockParagraphRichTextItemsMentionPage defines a model
type BlockParagraphRichTextItemsMentionPage struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// BlockParagraphRichTextItemsMentionUser defines a model
type BlockParagraphRichTextItemsMentionUser struct {
	Object    string                                       `json:"object,omitzero"`
	ID        uuid.UUID                                    `json:"id,omitzero"`
	Name      string                                       `json:"name,omitzero"`
	AvatarURL url.URL                                      `json:"avatar_url,omitzero"`
	Type      string                                       `json:"type,omitzero"`
	Person    BlockParagraphRichTextItemsMentionUserPerson `json:"person"`
}

// BlockParagraphRichTextItemsMentionUserPerson defines a model
type BlockParagraphRichTextItemsMentionUserPerson struct {
	Email types.Email `json:"email,omitzero"`
}

// BlockParagraphRichTextItemsText defines a model
type BlockParagraphRichTextItemsText struct {
	Content string                              `json:"content,omitzero"`
	Link    BlockParagraphRichTextItemsTextLink `json:"link"`
}

// BlockParagraphRichTextItemsTextLink defines a model
type BlockParagraphRichTextItemsTextLink struct {
	URL *url.URL `json:"url,omitempty"`
}

// BlockPdf defines a model
type BlockPdf struct {
	Caption  BlockPdfCaption  `json:"caption,omitempty"`
	Type     string           `json:"type,omitzero"`
	External BlockPdfExternal `json:"external"`
	File     *BlockPdfFile    `json:"file,omitempty"`
}

// BlockPdfCaption defines a model
type BlockPdfCaption []BlockPdfCaptionItems

// BlockPdfCaptionItems defines a model
type BlockPdfCaptionItems struct {
	Type        string                          `json:"type,omitzero"`
	Text        BlockPdfCaptionItemsText        `json:"text"`
	Annotations BlockPdfCaptionItemsAnnotations `json:"annotations"`
	PlainText   string                          `json:"plain_text,omitzero"`
	Href        struct{}                        `json:"href"`
}

// BlockPdfCaptionItemsAnnotations defines a model
type BlockPdfCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockPdfCaptionItemsText defines a model
type BlockPdfCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockPdfExternal defines a model
type BlockPdfExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// BlockPdfFile defines a model
type BlockPdfFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// BlockQuote defines a model
type BlockQuote struct {
	RichText BlockQuoteRichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// BlockQuoteRichText defines a model
type BlockQuoteRichText []BlockQuoteRichTextItems

// BlockQuoteRichTextItems defines a model
type BlockQuoteRichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        BlockQuoteRichTextItemsText        `json:"text"`
	Annotations BlockQuoteRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// BlockQuoteRichTextItemsAnnotations defines a model
type BlockQuoteRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockQuoteRichTextItemsText defines a model
type BlockQuoteRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockSyncedBlock defines a model
type BlockSyncedBlock struct {
	SyncedFrom BlockSyncedBlockSyncedFrom `json:"synced_from"`
}

// BlockSyncedBlockSyncedFrom defines a model
type BlockSyncedBlockSyncedFrom struct {
	Type    string    `json:"type,omitzero"`
	BlockID uuid.UUID `json:"block_id,omitzero"`
}

// BlockTable defines a model
type BlockTable struct {
	TableWidth      int  `json:"table_width,omitzero"`
	HasColumnHeader bool `json:"has_column_header,omitzero"`
	HasRowHeader    bool `json:"has_row_header,omitzero"`
}

// BlockTableOfContents defines a model
type BlockTableOfContents struct {
	Color string `json:"color,omitzero"`
}

// BlockToDo defines a model
type BlockToDo struct {
	RichText BlockToDoRichText `json:"rich_text,omitempty"`
	Checked  bool              `json:"checked,omitzero"`
	Color    string            `json:"color,omitzero"`
}

// BlockToDoRichText defines a model
type BlockToDoRichText []BlockToDoRichTextItems

// BlockToDoRichTextItems defines a model
type BlockToDoRichTextItems struct {
	Type        string                            `json:"type,omitzero"`
	Text        BlockToDoRichTextItemsText        `json:"text"`
	Annotations BlockToDoRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                            `json:"plain_text,omitzero"`
	Href        struct{}                          `json:"href"`
}

// BlockToDoRichTextItemsAnnotations defines a model
type BlockToDoRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockToDoRichTextItemsText defines a model
type BlockToDoRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockToggle defines a model
type BlockToggle struct {
	RichText BlockToggleRichText `json:"rich_text,omitempty"`
	Color    string              `json:"color,omitzero"`
}

// BlockToggleRichText defines a model
type BlockToggleRichText []BlockToggleRichTextItems

// BlockToggleRichTextItems defines a model
type BlockToggleRichTextItems struct {
	Type        string                              `json:"type,omitzero"`
	Text        BlockToggleRichTextItemsText        `json:"text"`
	Annotations BlockToggleRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                              `json:"plain_text,omitzero"`
	Href        struct{}                            `json:"href"`
}

// BlockToggleRichTextItemsAnnotations defines a model
type BlockToggleRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// BlockToggleRichTextItemsText defines a model
type BlockToggleRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// BlockVideo defines a model
type BlockVideo struct {
	Caption  []struct{}          `json:"caption,omitempty"`
	Type     string              `json:"type,omitzero"`
	File     BlockVideoFile      `json:"file"`
	External *BlockVideoExternal `json:"external,omitempty"`
}

// BlockVideoExternal defines a model
type BlockVideoExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// BlockVideoFile defines a model
type BlockVideoFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// Blocks defines a model
type Blocks []Block

// BlocksList defines a model
type BlocksList struct {
	Object  string `json:"object,omitzero"`
	Results Blocks `json:"results,omitempty"`
	// A unique identifier for a page, block, database, or user.
	NextCursor struct{}  `json:"next_cursor"`
	HasMore    bool      `json:"has_more,omitzero"`
	Type       string    `json:"type,omitzero"`
	Block      struct{}  `json:"block"`
	RequestID  uuid.UUID `json:"request_id,omitzero"`
}

// The color of the block.
type Color string

const (
	ColorDefault          Color = "default"
	ColorGray             Color = "gray"
	ColorBrown            Color = "brown"
	ColorOrange           Color = "orange"
	ColorYellow           Color = "yellow"
	ColorGreen            Color = "green"
	ColorBlue             Color = "blue"
	ColorPurple           Color = "purple"
	ColorPink             Color = "pink"
	ColorRed              Color = "red"
	ColorGrayBackground   Color = "gray_background"
	ColorBrownBackground  Color = "brown_background"
	ColorOrangeBackground Color = "orange_background"
	ColorYellowBackground Color = "yellow_background"
	ColorGreenBackground  Color = "green_background"
	ColorBlueBackground   Color = "blue_background"
	ColorPurpleBackground Color = "purple_background"
	ColorPinkBackground   Color = "pink_background"
	ColorRedBackground    Color = "red_background"
)

// An external file is any URL that isn't hosted by Notion.
type ExternalFile struct {
	// Link to the externally hosted content.
	URL url.URL `json:"url,omitzero"`
}

// File objects contain data about files uploaded to Notion as well as external files linked in Notion.
type File struct {
	// Type of this file object.
	Type FileType `json:"type,omitzero"`
	// An external file is any URL that isn't hosted by Notion.
	External *ExternalFile `json:"external,omitempty"`
}

// Type of this file object.
type FileType string

const (
	FileTypeFile     FileType = "file"
	FileTypeExternal FileType = "external"
)

// Page or database icon. It is either an emoji or a file.
type Icon struct {
	// Type of icon.
	Type IconType `json:"type,omitzero"`
	// Emoji character.
	Emoji string `json:"emoji,omitzero"`
}

// Type of icon.
type IconType string

const (
	IconTypeEmoji    IconType = "emoji"
	IconTypeFile     IconType = "file"
	IconTypeExternal IconType = "external"
)

// An inline link in a text.
type Link struct{}

// The Page object contains the [property values](https://developers.notion.com/reference/property-value-object) of a single Notion page.
//
// All pages have a parent. If the parent is a [database](https://developers.notion.com/reference/database), the property values conform to the schema laid out database's [properties](https://developers.notion.com/reference/property-object). Otherwise, the only property value is the `title`.
//
// Page content is available as [blocks](https://developers.notion.com/reference/block). The content can be read using [retrieve block children](https://developers.notion.com/reference/get-block-children) and appended using [append block children](https://developers.notion.com/reference/patch-block-children).
type Page struct {
	// Always "page".
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
	// Date and time when this page was created. Formatted as an ISO 8601 date time string.
	CreatedTime time.Time `json:"created_time,omitzero"`
	// Date and time when this page was updated. Formatted as an ISO 8601 date time string.
	LastEditedTime time.Time     `json:"last_edited_time,omitzero"`
	CreatedBy      UserReference `json:"created_by"`
	LastEditedBy   UserReference `json:"last_edited_by"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion.
	Cover File `json:"cover"`
	// Page or database icon. It is either an emoji or a file.
	Icon Icon `json:"icon"`
	// The `parent` property of a page or database contains these keys. Mandatory when creating, must be missing when updating.
	Parent Parent `json:"parent"`
	// The archived status of the page.
	Archived bool `json:"archived,omitzero"`
	InTrash  bool `json:"in_trash,omitzero"`
	// Properties of a page or database.
	Properties PropertyValues `json:"properties"`
	// The URL of the Notion page.
	URL       url.URL   `json:"url,omitzero"`
	PublicURL url.URL   `json:"public_url,omitzero"`
	RequestID uuid.UUID `json:"request_id,omitzero"`
}

// The `parent` property of a page or database contains these keys. Mandatory when creating, must be missing when updating.
type Parent struct {
	// The type of the parent.
	Type   ParentType `json:"type,omitzero"`
	PageID *uuid.UUID `json:"page_id,omitempty"`
}

// The type of the parent.
type ParentType string

const (
	ParentTypePageID     ParentType = "page_id"
	ParentTypeWorkspace  ParentType = "workspace"
	ParentTypeBlockID    ParentType = "block_id"
	ParentTypeDatabaseID ParentType = "database_id"
)

// Type of the property.
type PropertyType string

const (
	PropertyTypeRichText       PropertyType = "rich_text"
	PropertyTypeNumber         PropertyType = "number"
	PropertyTypeSelect         PropertyType = "select"
	PropertyTypeStatus         PropertyType = "status"
	PropertyTypeMultiSelect    PropertyType = "multi_select"
	PropertyTypeDate           PropertyType = "date"
	PropertyTypeFormula        PropertyType = "formula"
	PropertyTypeRelation       PropertyType = "relation"
	PropertyTypeRollup         PropertyType = "rollup"
	PropertyTypeTitle          PropertyType = "title"
	PropertyTypePeople         PropertyType = "people"
	PropertyTypeFiles          PropertyType = "files"
	PropertyTypeCheckbox       PropertyType = "checkbox"
	PropertyTypeURL            PropertyType = "url"
	PropertyTypeEmail          PropertyType = "email"
	PropertyTypePhoneNumber    PropertyType = "phone_number"
	PropertyTypeCreatedTime    PropertyType = "created_time"
	PropertyTypeCreatedBy      PropertyType = "created_by"
	PropertyTypeLastEditedTime PropertyType = "last_edited_time"
	PropertyTypeLastEditedBy   PropertyType = "last_edited_by"
	PropertyTypeButton         PropertyType = "button"
)

// A property value defines the identifier, type, and value of a page property in a page object. It's used when retrieving and updating pages ex: Create and Update pages.
type PropertyValue struct {
	/*
	   Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.

	   The id may be used in place of name when creating or updating pages.
	*/
	ID string `json:"id,omitzero"`
	// Type of the property.
	Type  PropertyType `json:"type,omitzero"`
	Title RichTexts    `json:"title,omitempty"`
}

// Properties of a page or database.
type PropertyValues map[string]PropertyValue

// Rich text objects contain data for displaying formatted text, mentions, and equations. A rich text object also contains annotations for style information. Arrays of rich text objects are used [within property objects](https://developers.notion.com/reference/database-property) and [property value objects](https://developers.notion.com/reference/page-property-value) to create what a user sees as a single text value in Notion.
type RichText struct {
	// Type of this rich text object.
	Type RichTextType `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	// The plain text without annotations.
	PlainText string `json:"plain_text,omitzero"`
	// The URL of any link or internal Notion mention in this text, if any.
	Href *url.URL `json:"href,omitempty"`
}

// Type of this rich text object.
type RichTextType string

const (
	RichTextTypeText     RichTextType = "text"
	RichTextTypeMention  RichTextType = "mention"
	RichTextTypeEquation RichTextType = "equation"
)

// RichTexts defines a model
type RichTexts []RichText

// Text objects contain this information within the `text` property of a RichText object.
type Text struct {
	// Text content. This field contains the actual content of your text and is probably the field you'll use most often.
	Content string `json:"content,omitzero"`
	// An inline link in a text.
	Link *Link `json:"link,omitempty"`
}

// UserReference defines a model
type UserReference struct {
	// Always "user"
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

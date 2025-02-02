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

// Annotations2 defines a model
type Annotations2 struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
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
	LastEditedTime time.Time     `json:"last_edited_time,omitzero"`
	CreatedBy      UserReference `json:"created_by"`
	LastEditedBy   UserReference `json:"last_edited_by"`
	// Whether or not the block has children blocks nested within it.
	HasChildren bool `json:"has_children,omitzero"`
	// The archived status of the block.
	Archived bool `json:"archived,omitzero"`
	InTrash  bool `json:"in_trash,omitzero"`
	// Type of block.
	Type             BlockType   `json:"type,omitzero"`
	Paragraph        *Paragraph  `json:"paragraph,omitempty"`
	Heading1         *Heading1   `json:"heading_1,omitempty"`
	Heading2         *Heading2   `json:"heading_2,omitempty"`
	Heading3         *Heading3   `json:"heading_3,omitempty"`
	BulletedListItem *Paragraph2 `json:"bulleted_list_item,omitempty"`
	NumberedListItem *Paragraph3 `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo       `json:"to_do,omitempty"`
	Toggle           *Paragraph4 `json:"toggle,omitempty"`
	Code             *Code       `json:"code,omitempty"`
	ChildPage        *Child      `json:"child_page,omitempty"`
	ChildDatabase    *Child      `json:"child_database,omitempty"`
	Embed            *Embed      `json:"embed,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Video *FileWithCaption `json:"video,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Audio    *FileWithCaption  `json:"audio,omitempty"`
	File     *FileWithCaption4 `json:"file,omitempty"`
	PDF      *FileWithCaption5 `json:"pdf,omitempty"`
	Bookmark *Embed3           `json:"bookmark,omitempty"`
	// Callout block objects contain the following information within the callout field.
	Callout  *Callout    `json:"callout,omitempty"`
	Quote    *Paragraph8 `json:"quote,omitempty"`
	Equation *Equation   `json:"equation,omitempty"`
	// Divider block objects do not contain any information within the divider property
	Divider         *struct{}        `json:"divider,omitempty"`
	TableOfContents *TableOfContents `json:"table_of_contents,omitempty"`
	// Column Lists are parent blocks for column children. They do not contain any information within the column_list property and can only contain children of type column.
	ColumnList  *struct{}    `json:"column_list,omitempty"`
	LinkPreview *LinkPreview `json:"link_preview,omitempty"`
	SyncedBlock *SyncedBlock `json:"synced_block,omitempty"`
	LinkToPage  *LinkToPage  `json:"link_to_page,omitempty"`
	Table       *Table       `json:"table,omitempty"`
	// Breadcrumb block objects do not contain any information within the breadcrumb property
	Breadcrumb *struct{} `json:"breadcrumb,omitempty"`
	// Some block types aren't available yet
	Unsupported *struct{} `json:"unsupported,omitempty"`
}

// Type of block.
type BlockType string

const (
	BlockTypeParagraph        BlockType = "paragraph"
	BlockTypeHeading1         BlockType = "heading_1"
	BlockTypeHeading2         BlockType = "heading_2"
	BlockTypeHeading3         BlockType = "heading_3"
	BlockTypeBulletedListItem BlockType = "bulleted_list_item"
	BlockTypeNumberedListItem BlockType = "numbered_list_item"
	BlockTypeToDo             BlockType = "to_do"
	BlockTypeToggle           BlockType = "toggle"
	BlockTypeChildPage        BlockType = "child_page"
	BlockTypeChildDatabase    BlockType = "child_database"
	BlockTypeEmbed            BlockType = "embed"
	BlockTypeImage            BlockType = "image"
	BlockTypeVideo            BlockType = "video"
	BlockTypeFile             BlockType = "file"
	BlockTypePdf              BlockType = "pdf"
	BlockTypeBookmark         BlockType = "bookmark"
	BlockTypeCallout          BlockType = "callout"
	BlockTypeQuote            BlockType = "quote"
	BlockTypeEquation         BlockType = "equation"
	BlockTypeDivider          BlockType = "divider"
	BlockTypeTableOfContents  BlockType = "table_of_contents"
	BlockTypeColumn           BlockType = "column"
	BlockTypeColumnList       BlockType = "column_list"
	BlockTypeLinkPreview      BlockType = "link_preview"
	BlockTypeSyncedBlock      BlockType = "synced_block"
	BlockTypeTemplate         BlockType = "template"
	BlockTypeLinkToPage       BlockType = "link_to_page"
	BlockTypeTable            BlockType = "table"
	BlockTypeTableRow         BlockType = "table_row"
	BlockTypeCode             BlockType = "code"
	BlockTypeBreadcrumb       BlockType = "breadcrumb"
	BlockTypeAudio            BlockType = "audio"
	BlockTypeUnsupported      BlockType = "unsupported"
)

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

// Callout block objects contain the following information within the callout field.
type Callout struct {
	RichText CalloutRichText `json:"rich_text,omitempty"`
	Icon     CalloutIcon     `json:"icon"`
	Color    string          `json:"color,omitzero"`
}

// CalloutIcon defines a model
type CalloutIcon struct {
	Type     string               `json:"type,omitzero"`
	Emoji    string               `json:"emoji,omitzero"`
	External *CalloutIconExternal `json:"external,omitempty"`
	File     *CalloutIconFile     `json:"file,omitempty"`
}

// CalloutIconExternal defines a model
type CalloutIconExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// CalloutIconFile defines a model
type CalloutIconFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

// CalloutRichText defines a model
type CalloutRichText []CalloutRichTextItems

// CalloutRichTextItems defines a model
type CalloutRichTextItems struct {
	Type        string                          `json:"type,omitzero"`
	Text        CalloutRichTextItemsText        `json:"text"`
	Annotations CalloutRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                          `json:"plain_text,omitzero"`
	Href        struct{}                        `json:"href"`
}

// CalloutRichTextItemsAnnotations defines a model
type CalloutRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// CalloutRichTextItemsText defines a model
type CalloutRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Child defines a model
type Child struct {
	Title string `json:"title,omitzero"`
}

// Code defines a model
type Code struct {
	Caption  CodeCaption  `json:"caption,omitempty"`
	RichText CodeRichText `json:"rich_text,omitempty"`
	Language string       `json:"language,omitzero"`
}

// CodeCaption defines a model
type CodeCaption []CodeCaptionItems

// CodeCaptionItems defines a model
type CodeCaptionItems struct {
	Type        string                      `json:"type,omitzero"`
	Text        CodeCaptionItemsText        `json:"text"`
	Annotations CodeCaptionItemsAnnotations `json:"annotations"`
	PlainText   string                      `json:"plain_text,omitzero"`
	Href        struct{}                    `json:"href"`
}

// CodeCaptionItemsAnnotations defines a model
type CodeCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// CodeCaptionItemsText defines a model
type CodeCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// CodeRichText defines a model
type CodeRichText []CodeRichTextItems

// CodeRichTextItems defines a model
type CodeRichTextItems struct {
	Type        string                       `json:"type,omitzero"`
	Text        CodeRichTextItemsText        `json:"text"`
	Annotations CodeRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                       `json:"plain_text,omitzero"`
	Href        url.URL                      `json:"href,omitzero"`
}

// CodeRichTextItemsAnnotations defines a model
type CodeRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// CodeRichTextItemsText defines a model
type CodeRichTextItemsText struct {
	Content string                    `json:"content,omitzero"`
	Link    CodeRichTextItemsTextLink `json:"link"`
}

// CodeRichTextItemsTextLink defines a model
type CodeRichTextItemsTextLink struct {
	URL *url.URL `json:"url,omitempty"`
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

// Embed defines a model
type Embed struct {
	Caption EmbedCaption `json:"caption,omitempty"`
	URL     url.URL      `json:"url,omitzero"`
}

// Embed3 defines a model
type Embed3 struct {
	Caption Embed3Caption `json:"caption,omitempty"`
	URL     url.URL       `json:"url,omitzero"`
}

// Embed3Caption defines a model
type Embed3Caption []Embed3CaptionItems

// Embed3CaptionItems defines a model
type Embed3CaptionItems struct {
	Type        string                        `json:"type,omitzero"`
	Text        Embed3CaptionItemsText        `json:"text"`
	Annotations Embed3CaptionItemsAnnotations `json:"annotations"`
	PlainText   string                        `json:"plain_text,omitzero"`
	Href        struct{}                      `json:"href"`
}

// Embed3CaptionItemsAnnotations defines a model
type Embed3CaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Embed3CaptionItemsText defines a model
type Embed3CaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// EmbedCaption defines a model
type EmbedCaption []EmbedCaptionItems

// EmbedCaptionItems defines a model
type EmbedCaptionItems struct {
	Type        string                        `json:"type,omitzero"`
	Text        *EmbedCaptionItemsText        `json:"text,omitempty"`
	Annotations *EmbedCaptionItemsAnnotations `json:"annotations,omitempty"`
	PlainText   string                        `json:"plain_text,omitzero"`
	Href        *struct{}                     `json:"href,omitempty"`
}

// EmbedCaptionItemsAnnotations defines a model
type EmbedCaptionItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// EmbedCaptionItemsText defines a model
type EmbedCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Equation defines a model
type Equation struct {
	Expression string `json:"expression,omitzero"`
}

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

// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
type FileWithCaption struct {
	Caption RichTexts `json:"caption,omitempty"`
	// Type of this file object.
	Type FileType    `json:"type,omitzero"`
	File *NotionFile `json:"file,omitempty"`
	// An external file is any URL that isn't hosted by Notion.
	External *ExternalFile `json:"external,omitempty"`
	Name     string        `json:"name,omitzero"`
}

// FileWithCaption4 defines a model
type FileWithCaption4 struct {
	Caption RichTexts  `json:"caption,omitempty"`
	Type    string     `json:"type,omitzero"`
	File    NotionFile `json:"file"`
	Name    string     `json:"name,omitzero"`
}

// FileWithCaption5 defines a model
type FileWithCaption5 struct {
	Caption RichTexts2 `json:"caption,omitempty"`
	Type    string     `json:"type,omitzero"`
	// An external file is any URL that isn't hosted by Notion.
	External ExternalFile `json:"external"`
	File     *NotionFile  `json:"file,omitempty"`
}

// Heading1 defines a model
type Heading1 struct {
	RichText     Heading1RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading1RichText defines a model
type Heading1RichText []Heading1RichTextItems

// Heading1RichTextItems defines a model
type Heading1RichTextItems struct {
	Type        string                           `json:"type,omitzero"`
	Text        Heading1RichTextItemsText        `json:"text"`
	Annotations Heading1RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                           `json:"plain_text,omitzero"`
	Href        struct{}                         `json:"href"`
}

// Heading1RichTextItemsAnnotations defines a model
type Heading1RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Heading1RichTextItemsText defines a model
type Heading1RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Heading2 defines a model
type Heading2 struct {
	RichText     Heading2RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading2RichText defines a model
type Heading2RichText []Heading2RichTextItems

// Heading2RichTextItems defines a model
type Heading2RichTextItems struct {
	Type        string                           `json:"type,omitzero"`
	Text        Heading2RichTextItemsText        `json:"text"`
	Annotations Heading2RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                           `json:"plain_text,omitzero"`
	Href        struct{}                         `json:"href"`
}

// Heading2RichTextItemsAnnotations defines a model
type Heading2RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Heading2RichTextItemsText defines a model
type Heading2RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Heading3 defines a model
type Heading3 struct {
	RichText     Heading3RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading3RichText defines a model
type Heading3RichText []Heading3RichTextItems

// Heading3RichTextItems defines a model
type Heading3RichTextItems struct {
	Type        string                           `json:"type,omitzero"`
	Text        Heading3RichTextItemsText        `json:"text"`
	Annotations Heading3RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                           `json:"plain_text,omitzero"`
	Href        struct{}                         `json:"href"`
}

// Heading3RichTextItemsAnnotations defines a model
type Heading3RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Heading3RichTextItemsText defines a model
type Heading3RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

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

// LinkPreview defines a model
type LinkPreview struct {
	URL url.URL `json:"url,omitzero"`
}

// LinkToPage defines a model
type LinkToPage struct {
	Type   string    `json:"type,omitzero"`
	PageID uuid.UUID `json:"page_id,omitzero"`
}

// NotionFile defines a model
type NotionFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
}

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

// Paragraph defines a model
type Paragraph struct {
	RichText ParagraphRichText `json:"rich_text,omitempty"`
	Color    string            `json:"color,omitzero"`
}

// Paragraph2 defines a model
type Paragraph2 struct {
	RichText Paragraph2RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph2RichText defines a model
type Paragraph2RichText []Paragraph2RichTextItems

// Paragraph2RichTextItems defines a model
type Paragraph2RichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        Paragraph2RichTextItemsText        `json:"text"`
	Annotations Paragraph2RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// Paragraph2RichTextItemsAnnotations defines a model
type Paragraph2RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Paragraph2RichTextItemsText defines a model
type Paragraph2RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Paragraph3 defines a model
type Paragraph3 struct {
	RichText Paragraph3RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph3RichText defines a model
type Paragraph3RichText []Paragraph3RichTextItems

// Paragraph3RichTextItems defines a model
type Paragraph3RichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        Paragraph3RichTextItemsText        `json:"text"`
	Annotations Paragraph3RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// Paragraph3RichTextItemsAnnotations defines a model
type Paragraph3RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Paragraph3RichTextItemsText defines a model
type Paragraph3RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Paragraph4 defines a model
type Paragraph4 struct {
	RichText Paragraph4RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph4RichText defines a model
type Paragraph4RichText []Paragraph4RichTextItems

// Paragraph4RichTextItems defines a model
type Paragraph4RichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        Paragraph4RichTextItemsText        `json:"text"`
	Annotations Paragraph4RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// Paragraph4RichTextItemsAnnotations defines a model
type Paragraph4RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Paragraph4RichTextItemsText defines a model
type Paragraph4RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Paragraph8 defines a model
type Paragraph8 struct {
	RichText Paragraph8RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph8RichText defines a model
type Paragraph8RichText []Paragraph8RichTextItems

// Paragraph8RichTextItems defines a model
type Paragraph8RichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        Paragraph8RichTextItemsText        `json:"text"`
	Annotations Paragraph8RichTextItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// Paragraph8RichTextItemsAnnotations defines a model
type Paragraph8RichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// Paragraph8RichTextItemsText defines a model
type Paragraph8RichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// ParagraphRichText defines a model
type ParagraphRichText []ParagraphRichTextItems

// ParagraphRichTextItems defines a model
type ParagraphRichTextItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        *ParagraphRichTextItemsText        `json:"text,omitempty"`
	Annotations *ParagraphRichTextItemsAnnotations `json:"annotations,omitempty"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        *url.URL                           `json:"href,omitempty"`
	Mention     *ParagraphRichTextItemsMention     `json:"mention,omitempty"`
	Equation    *ParagraphRichTextItemsEquation    `json:"equation,omitempty"`
}

// ParagraphRichTextItemsAnnotations defines a model
type ParagraphRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// ParagraphRichTextItemsEquation defines a model
type ParagraphRichTextItemsEquation struct {
	Expression string `json:"expression,omitzero"`
}

// ParagraphRichTextItemsMention defines a model
type ParagraphRichTextItemsMention struct {
	Type        string                                   `json:"type,omitzero"`
	LinkMention ParagraphRichTextItemsMentionLinkMention `json:"link_mention"`
	Database    *ParagraphRichTextItemsMentionDatabase   `json:"database,omitempty"`
	User        *ParagraphRichTextItemsMentionUser       `json:"user,omitempty"`
	Date        *ParagraphRichTextItemsMentionDate       `json:"date,omitempty"`
	Page        *ParagraphRichTextItemsMentionPage       `json:"page,omitempty"`
}

// ParagraphRichTextItemsMentionDatabase defines a model
type ParagraphRichTextItemsMentionDatabase struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// ParagraphRichTextItemsMentionDate defines a model
type ParagraphRichTextItemsMentionDate struct {
	Start    string   `json:"start,omitzero"`
	End      struct{} `json:"end"`
	TimeZone struct{} `json:"time_zone"`
}

// ParagraphRichTextItemsMentionLinkMention defines a model
type ParagraphRichTextItemsMentionLinkMention struct {
	Href        url.URL `json:"href,omitzero"`
	Title       string  `json:"title,omitzero"`
	Description string  `json:"description,omitzero"`
}

// ParagraphRichTextItemsMentionPage defines a model
type ParagraphRichTextItemsMentionPage struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// ParagraphRichTextItemsMentionUser defines a model
type ParagraphRichTextItemsMentionUser struct {
	Object    string                                  `json:"object,omitzero"`
	ID        uuid.UUID                               `json:"id,omitzero"`
	Name      string                                  `json:"name,omitzero"`
	AvatarURL url.URL                                 `json:"avatar_url,omitzero"`
	Type      string                                  `json:"type,omitzero"`
	Person    ParagraphRichTextItemsMentionUserPerson `json:"person"`
}

// ParagraphRichTextItemsMentionUserPerson defines a model
type ParagraphRichTextItemsMentionUserPerson struct {
	Email types.Email `json:"email,omitzero"`
}

// ParagraphRichTextItemsText defines a model
type ParagraphRichTextItemsText struct {
	Content string                         `json:"content,omitzero"`
	Link    ParagraphRichTextItemsTextLink `json:"link"`
}

// ParagraphRichTextItemsTextLink defines a model
type ParagraphRichTextItemsTextLink struct {
	URL *url.URL `json:"url,omitempty"`
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

// RichText2 defines a model
type RichText2 struct {
	// Type of this rich text object.
	Type        RichTextType `json:"type,omitzero"`
	Text        Text2        `json:"text"`
	Annotations Annotations2 `json:"annotations"`
	PlainText   string       `json:"plain_text,omitzero"`
	Href        struct{}     `json:"href"`
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

// RichTexts2 defines a model
type RichTexts2 []RichText2

// SyncedBlock defines a model
type SyncedBlock struct {
	SyncedFrom SyncedBlockSyncedFrom `json:"synced_from"`
}

// SyncedBlockSyncedFrom defines a model
type SyncedBlockSyncedFrom struct {
	Type    string    `json:"type,omitzero"`
	BlockID uuid.UUID `json:"block_id,omitzero"`
}

// Table defines a model
type Table struct {
	TableWidth      int  `json:"table_width,omitzero"`
	HasColumnHeader bool `json:"has_column_header,omitzero"`
	HasRowHeader    bool `json:"has_row_header,omitzero"`
}

// TableOfContents defines a model
type TableOfContents struct {
	Color string `json:"color,omitzero"`
}

// Text objects contain this information within the `text` property of a RichText object.
type Text struct {
	// Text content. This field contains the actual content of your text and is probably the field you'll use most often.
	Content string `json:"content,omitzero"`
	// An inline link in a text.
	Link *Link `json:"link,omitempty"`
}

// Text2 defines a model
type Text2 struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// ToDo defines a model
type ToDo struct {
	RichText ToDoRichText `json:"rich_text,omitempty"`
	Checked  bool         `json:"checked,omitzero"`
	Color    string       `json:"color,omitzero"`
}

// ToDoRichText defines a model
type ToDoRichText []ToDoRichTextItems

// ToDoRichTextItems defines a model
type ToDoRichTextItems struct {
	Type        string                       `json:"type,omitzero"`
	Text        ToDoRichTextItemsText        `json:"text"`
	Annotations ToDoRichTextItemsAnnotations `json:"annotations"`
	PlainText   string                       `json:"plain_text,omitzero"`
	Href        struct{}                     `json:"href"`
}

// ToDoRichTextItemsAnnotations defines a model
type ToDoRichTextItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// ToDoRichTextItemsText defines a model
type ToDoRichTextItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// UserReference defines a model
type UserReference struct {
	// Always "user"
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

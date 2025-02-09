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
	LastEditedTime time.Time     `json:"last_edited_time,omitzero"`
	CreatedBy      UserReference `json:"created_by"`
	LastEditedBy   UserReference `json:"last_edited_by"`
	// Whether or not the block has children blocks nested within it.
	HasChildren bool `json:"has_children,omitzero"`
	// The archived status of the block.
	Archived bool `json:"archived,omitzero"`
	InTrash  bool `json:"in_trash,omitzero"`
	// Type of block.
	Type BlockType `json:"type,omitzero"`
	// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
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
	// Embed blocks include block types that allow displaying another website within Notion.
	Embed *Embed `json:"embed,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Video *FileWithCaption `json:"video,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Audio *FileWithCaption `json:"audio,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	File *FileWithCaption `json:"file,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	PDF      *FileWithCaption `json:"pdf,omitempty"`
	Bookmark *Embed2          `json:"bookmark,omitempty"`
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
type CalloutRichText []RichText4

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
	Type string               `json:"type,omitzero"`
	Text CodeCaptionItemsText `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// CodeCaptionItemsText defines a model
type CodeCaptionItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// CodeRichText defines a model
type CodeRichText []RichText5

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

// Embed blocks include block types that allow displaying another website within Notion.
type Embed struct {
	Caption RichTexts2 `json:"caption,omitempty"`
	// Embedded link.
	URL url.URL `json:"url,omitzero"`
}

// Embed2 defines a model
type Embed2 struct {
	Caption RichTexts22 `json:"caption,omitempty"`
	URL     url.URL     `json:"url,omitzero"`
}

// Equation defines a model
type Equation struct {
	Expression string `json:"expression,omitzero"`
}

// Equation2 defines a model
type Equation2 struct {
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

// Heading1 defines a model
type Heading1 struct {
	RichText     Heading1RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading1RichText defines a model
type Heading1RichText []RichText6

// Heading2 defines a model
type Heading2 struct {
	RichText     Heading2RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading2RichText defines a model
type Heading2RichText []RichText7

// Heading3 defines a model
type Heading3 struct {
	RichText     Heading3RichText `json:"rich_text,omitempty"`
	IsToggleable bool             `json:"is_toggleable,omitzero"`
	Color        string           `json:"color,omitzero"`
}

// Heading3RichText defines a model
type Heading3RichText []RichText8

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
type Link struct {
	URL *url.URL `json:"url,omitempty"`
}

// LinkMention defines a model
type LinkMention struct {
	Href        url.URL `json:"href,omitzero"`
	Title       string  `json:"title,omitzero"`
	Description string  `json:"description,omitzero"`
}

// LinkPreview defines a model
type LinkPreview struct {
	URL url.URL `json:"url,omitzero"`
}

// LinkToPage defines a model
type LinkToPage struct {
	Type   string    `json:"type,omitzero"`
	PageID uuid.UUID `json:"page_id,omitzero"`
}

// Mention defines a model
type Mention struct {
	// Type of the inline mention.
	Type        MentionType      `json:"type,omitzero"`
	LinkMention *LinkMention     `json:"link_mention,omitempty"`
	Database    *MentionDatabase `json:"database,omitempty"`
	User        *MentionUser     `json:"user,omitempty"`
	Date        *MentionDate     `json:"date,omitempty"`
	Page        *MentionPage     `json:"page,omitempty"`
}

// MentionDatabase defines a model
type MentionDatabase struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// MentionDate defines a model
type MentionDate struct {
	Start    string   `json:"start,omitzero"`
	End      struct{} `json:"end"`
	TimeZone struct{} `json:"time_zone"`
}

// MentionPage defines a model
type MentionPage struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// Type of the inline mention.
type MentionType string

const (
	MentionTypeLinkMention MentionType = "link_mention"
)

// MentionUser defines a model
type MentionUser struct {
	Object    string            `json:"object,omitzero"`
	ID        uuid.UUID         `json:"id,omitzero"`
	Name      string            `json:"name,omitzero"`
	AvatarURL url.URL           `json:"avatar_url,omitzero"`
	Type      string            `json:"type,omitzero"`
	Person    MentionUserPerson `json:"person"`
}

// MentionUserPerson defines a model
type MentionUserPerson struct {
	Email types.Email `json:"email,omitzero"`
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

// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
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
type Paragraph2RichText []RichText9

// Paragraph3 defines a model
type Paragraph3 struct {
	RichText Paragraph3RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph3RichText defines a model
type Paragraph3RichText []RichText

// Paragraph4 defines a model
type Paragraph4 struct {
	RichText Paragraph4RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph4RichText defines a model
type Paragraph4RichText []RichText

// Paragraph8 defines a model
type Paragraph8 struct {
	RichText Paragraph8RichText `json:"rich_text,omitempty"`
	Color    string             `json:"color,omitzero"`
}

// Paragraph8RichText defines a model
type Paragraph8RichText []RichText

// ParagraphRichText defines a model
type ParagraphRichText []RichText

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
	Text     Text       `json:"text"`
	Mention  *Mention   `json:"mention,omitempty"`
	Equation *Equation2 `json:"equation,omitempty"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	// The plain text without annotations.
	PlainText string `json:"plain_text,omitzero"`
	// The URL of any link or internal Notion mention in this text, if any.
	Href *url.URL `json:"href,omitempty"`
}

// RichText2 defines a model
type RichText2 struct {
	Type string         `json:"type,omitzero"`
	Text *RichText2Text `json:"text,omitempty"`
	// Style information which applies to the whole rich text object.
	Annotations *Annotations `json:"annotations,omitempty"`
	PlainText   string       `json:"plain_text,omitzero"`
	Href        *struct{}    `json:"href,omitempty"`
}

// RichText2Text defines a model
type RichText2Text struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// RichText3 defines a model
type RichText3 struct {
	Type string        `json:"type,omitzero"`
	Text RichText3Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// RichText3Text defines a model
type RichText3Text struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// RichText4 defines a model
type RichText4 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// RichText5 defines a model
type RichText5 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        url.URL     `json:"href,omitzero"`
}

// RichText6 defines a model
type RichText6 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// RichText7 defines a model
type RichText7 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// RichText8 defines a model
type RichText8 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
}

// RichText9 defines a model
type RichText9 struct {
	Type string `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text Text `json:"text"`
	// Style information which applies to the whole rich text object.
	Annotations Annotations `json:"annotations"`
	PlainText   string      `json:"plain_text,omitzero"`
	Href        struct{}    `json:"href"`
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

// RichTexts22 defines a model
type RichTexts22 []RichText3

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

// ToDo defines a model
type ToDo struct {
	RichText ToDoRichText `json:"rich_text,omitempty"`
	Checked  bool         `json:"checked,omitzero"`
	Color    string       `json:"color,omitzero"`
}

// ToDoRichText defines a model
type ToDoRichText []RichText

// UserReference defines a model
type UserReference struct {
	// Always "user"
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

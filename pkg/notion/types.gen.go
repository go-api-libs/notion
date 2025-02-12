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
	Paragraph *Paragraph `json:"paragraph,omitempty"`
	// Heading block objects contain this information within their respective property.
	Heading1 *Heading `json:"heading_1,omitempty"`
	// Heading block objects contain this information within their respective property.
	Heading2 *Heading `json:"heading_2,omitempty"`
	// Heading block objects contain this information within their respective property.
	Heading3 *Heading `json:"heading_3,omitempty"`
	// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
	BulletedListItem *Paragraph `json:"bulleted_list_item,omitempty"`
	// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
	NumberedListItem *Paragraph `json:"numbered_list_item,omitempty"`
	ToDo             *ToDo      `json:"to_do,omitempty"`
	// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
	Toggle *Paragraph `json:"toggle,omitempty"`
	// Code block objects contain this information within the `code` property.
	Code          *Code  `json:"code,omitempty"`
	ChildPage     *Child `json:"child_page,omitempty"`
	ChildDatabase *Child `json:"child_database,omitempty"`
	// Embed blocks include block types that allow displaying another website within Notion.
	Embed *Embed `json:"embed,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Video *FileWithCaption `json:"video,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	Audio *FileWithCaption `json:"audio,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	File *FileWithCaption `json:"file,omitempty"`
	// File objects contain data about files uploaded to Notion as well as external files linked in Notion. A PDF can also have a caption.
	PDF *FileWithCaption `json:"pdf,omitempty"`
	// Embed blocks include block types that allow displaying another website within Notion.
	Bookmark *Embed `json:"bookmark,omitempty"`
	// Callout block objects contain the following information within the callout field.
	Callout *Callout `json:"callout,omitempty"`
	// Paragraph, quote, toggle and list item block objects contain this information within their respective property.
	Quote *Paragraph `json:"quote,omitempty"`
	// Equation block objects contain this information within the `equation` property
	Equation *Equation `json:"equation,omitempty"`
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
	RichText RichTexts `json:"rich_text,omitempty"`
	// Page or database icon. It is either an emoji or a file.
	Icon Icon `json:"icon"`
	// The color of the block.
	Color Color `json:"color,omitzero"`
}

// Child defines a model
type Child struct {
	Title string `json:"title,omitzero"`
}

// Code block objects contain this information within the `code` property.
type Code struct {
	Caption  RichTexts `json:"caption,omitempty"`
	RichText RichTexts `json:"rich_text,omitempty"`
	// Coding language in code block
	Language CodeLanguage `json:"language,omitzero"`
}

// Coding language in code block
type CodeLanguage string

const (
	CodeLanguageAbap                 CodeLanguage = "abap"
	CodeLanguageArduino              CodeLanguage = "arduino"
	CodeLanguageBash                 CodeLanguage = "bash"
	CodeLanguageBasic                CodeLanguage = "basic"
	CodeLanguageC                    CodeLanguage = "c"
	CodeLanguageClojure              CodeLanguage = "clojure"
	CodeLanguageCoffeescript         CodeLanguage = "coffeescript"
	CodeLanguageCPlusPlus            CodeLanguage = "c++"
	CodeLanguageCSharp               CodeLanguage = "c#"
	CodeLanguageCSS                  CodeLanguage = "css"
	CodeLanguageDart                 CodeLanguage = "dart"
	CodeLanguageDiff                 CodeLanguage = "diff"
	CodeLanguageDocker               CodeLanguage = "docker"
	CodeLanguageElixir               CodeLanguage = "elixir"
	CodeLanguageElm                  CodeLanguage = "elm"
	CodeLanguageErlang               CodeLanguage = "erlang"
	CodeLanguageFlow                 CodeLanguage = "flow"
	CodeLanguageFortran              CodeLanguage = "fortran"
	CodeLanguageFSharp               CodeLanguage = "f#"
	CodeLanguageGherkin              CodeLanguage = "gherkin"
	CodeLanguageGlsl                 CodeLanguage = "glsl"
	CodeLanguageGo                   CodeLanguage = "go"
	CodeLanguageGraphql              CodeLanguage = "graphql"
	CodeLanguageGroovy               CodeLanguage = "groovy"
	CodeLanguageHaskell              CodeLanguage = "haskell"
	CodeLanguageHTML                 CodeLanguage = "html"
	CodeLanguageJava                 CodeLanguage = "java"
	CodeLanguageJavascript           CodeLanguage = "javascript"
	CodeLanguageJSON                 CodeLanguage = "json"
	CodeLanguageJulia                CodeLanguage = "julia"
	CodeLanguageKotlin               CodeLanguage = "kotlin"
	CodeLanguageLatex                CodeLanguage = "latex"
	CodeLanguageLess                 CodeLanguage = "less"
	CodeLanguageLisp                 CodeLanguage = "lisp"
	CodeLanguageLivescript           CodeLanguage = "livescript"
	CodeLanguageLua                  CodeLanguage = "lua"
	CodeLanguageMakefile             CodeLanguage = "makefile"
	CodeLanguageMarkdown             CodeLanguage = "markdown"
	CodeLanguageMarkup               CodeLanguage = "markup"
	CodeLanguageMatlab               CodeLanguage = "matlab"
	CodeLanguageMermaid              CodeLanguage = "mermaid"
	CodeLanguageNix                  CodeLanguage = "nix"
	CodeLanguageObjectiveC           CodeLanguage = "objective-c"
	CodeLanguageOcaml                CodeLanguage = "ocaml"
	CodeLanguagePascal               CodeLanguage = "pascal"
	CodeLanguagePerl                 CodeLanguage = "perl"
	CodeLanguagePhp                  CodeLanguage = "php"
	CodeLanguagePlainText            CodeLanguage = "plain text"
	CodeLanguagePowershell           CodeLanguage = "powershell"
	CodeLanguageProlog               CodeLanguage = "prolog"
	CodeLanguageProtobuf             CodeLanguage = "protobuf"
	CodeLanguagePython               CodeLanguage = "python"
	CodeLanguageR                    CodeLanguage = "r"
	CodeLanguageReason               CodeLanguage = "reason"
	CodeLanguageRuby                 CodeLanguage = "ruby"
	CodeLanguageRust                 CodeLanguage = "rust"
	CodeLanguageSass                 CodeLanguage = "sass"
	CodeLanguageScala                CodeLanguage = "scala"
	CodeLanguageScheme               CodeLanguage = "scheme"
	CodeLanguageScss                 CodeLanguage = "scss"
	CodeLanguageShell                CodeLanguage = "shell"
	CodeLanguageSQL                  CodeLanguage = "sql"
	CodeLanguageSwift                CodeLanguage = "swift"
	CodeLanguageTypescript           CodeLanguage = "typescript"
	CodeLanguageVbNet                CodeLanguage = "vb.net"
	CodeLanguageVerilog              CodeLanguage = "verilog"
	CodeLanguageVhdl                 CodeLanguage = "vhdl"
	CodeLanguageVisualBasic          CodeLanguage = "visual basic"
	CodeLanguageWebassembly          CodeLanguage = "webassembly"
	CodeLanguageXML                  CodeLanguage = "xml"
	CodeLanguageYaml                 CodeLanguage = "yaml"
	CodeLanguageJavaCCPlusPlusCSharp CodeLanguage = "java/c/c++/c#"
)

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

// Date defines a model
type Date struct {
	// An ISO 8601 format date, with optional time.
	Start string `json:"start,omitzero"`
	// An ISO 8601 formatted date, with optional time. Represents the end of a date range.
	//
	// If `null`, this property's date value is not a range.
	End *struct{} `json:"end,omitempty"`
	// Time zone information for start and end. Possible values are extracted from the IANA database and they are based on the time zones from Moment.js.
	//
	// When time zone is provided, start and end should not have any UTC offset. In addition, when time zone is provided, start and end cannot be dates without time information.
	//
	// If null, time zone information will be contained in UTC offsets in start and end.
	TimeZone *struct{} `json:"time_zone,omitempty"`
}

// Embed blocks include block types that allow displaying another website within Notion.
type Embed struct {
	Caption RichTexts `json:"caption,omitempty"`
	// Embedded link.
	URL url.URL `json:"url,omitzero"`
}

// Equation block objects contain this information within the `equation` property
type Equation struct {
	// A KaTeX compatible string
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

// Heading block objects contain this information within their respective property.
type Heading struct {
	RichText     RichTexts `json:"rich_text,omitempty"`
	IsToggleable bool      `json:"is_toggleable,omitzero"`
	// The color of the block.
	Color Color `json:"color,omitzero"`
}

// Page or database icon. It is either an emoji or a file.
type Icon struct {
	// Type of icon.
	Type IconType `json:"type,omitzero"`
	// Emoji character.
	Emoji    string        `json:"emoji,omitzero"`
	External *IconExternal `json:"external,omitempty"`
	File     *IconFile     `json:"file,omitempty"`
}

// IconExternal defines a model
type IconExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// IconFile defines a model
type IconFile struct {
	URL        url.URL   `json:"url,omitzero"`
	ExpiryTime time.Time `json:"expiry_time,omitzero"`
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
	Type        MentionType  `json:"type,omitzero"`
	LinkMention *LinkMention `json:"link_mention,omitempty"`
	User        *User2       `json:"user,omitempty"`
	Page        *Reference   `json:"page,omitempty"`
	Database    *Reference   `json:"database,omitempty"`
	Date        *Date        `json:"date,omitempty"`
}

// Type of the inline mention.
type MentionType string

const (
	MentionTypeLinkMention MentionType = "link_mention"
)

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
	// Underlying identifier for the property. This identifier is guaranteed to remain constant when the property name changes. It may be a UUID, but is often a short random string.
	//
	// The id may be used in place of name when creating or updating pages.
	ID string `json:"id,omitzero"`
	// Type of the property.
	Type  PropertyType `json:"type,omitzero"`
	Title RichTexts    `json:"title,omitempty"`
}

// Properties of a page or database.
type PropertyValues map[string]PropertyValue

// Reference defines a model
type Reference struct {
	ID uuid.UUID `json:"id,omitzero"`
}

// Rich text objects contain data for displaying formatted text, mentions, and equations. A rich text object also contains annotations for style information. Arrays of rich text objects are used [within property objects](https://developers.notion.com/reference/database-property) and [property value objects](https://developers.notion.com/reference/page-property-value) to create what a user sees as a single text value in Notion.
type RichText struct {
	// Type of this rich text object.
	Type RichTextType `json:"type,omitzero"`
	// Text objects contain this information within the `text` property of a RichText object.
	Text    Text     `json:"text"`
	Mention *Mention `json:"mention,omitempty"`
	// Equation block objects contain this information within the `equation` property
	Equation *Equation `json:"equation,omitempty"`
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

// User2 defines a model
type User2 struct {
	Object    string      `json:"object,omitzero"`
	ID        uuid.UUID   `json:"id,omitzero"`
	Name      string      `json:"name,omitzero"`
	AvatarURL url.URL     `json:"avatar_url,omitzero"`
	Type      string      `json:"type,omitzero"`
	Person    User2Person `json:"person"`
}

// User2Person defines a model
type User2Person struct {
	Email types.Email `json:"email,omitzero"`
}

// UserReference defines a model
type UserReference struct {
	// Always "user"
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

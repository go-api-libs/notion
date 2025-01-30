// This file provides types for the API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.

package notion

import (
	"net/url"
	"time"

	"github.com/google/uuid"
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
	LastEditedTime time.Time      `json:"last_edited_time,omitzero"`
	CreatedBy      UserReference  `json:"created_by"`
	LastEditedBy   UserReference2 `json:"last_edited_by"`
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

// PropertyValue defines a model
type PropertyValue struct {
	ID    string             `json:"id,omitzero"`
	Type  string             `json:"type,omitzero"`
	Title PropertyValueTitle `json:"title,omitempty"`
}

// PropertyValueTitle defines a model
type PropertyValueTitle []PropertyValueTitleItems

// PropertyValueTitleItems defines a model
type PropertyValueTitleItems struct {
	Type        string                             `json:"type,omitzero"`
	Text        PropertyValueTitleItemsText        `json:"text"`
	Annotations PropertyValueTitleItemsAnnotations `json:"annotations"`
	PlainText   string                             `json:"plain_text,omitzero"`
	Href        struct{}                           `json:"href"`
}

// PropertyValueTitleItemsAnnotations defines a model
type PropertyValueTitleItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// PropertyValueTitleItemsText defines a model
type PropertyValueTitleItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

// Properties of a page or database.
type PropertyValues map[string]PropertyValue

// UserReference defines a model
type UserReference struct {
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

// UserReference2 defines a model
type UserReference2 struct {
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

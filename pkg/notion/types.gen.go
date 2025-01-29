// This file provides types for the API.
//
// Code generated by github.com/MarkRosemaker DO NOT EDIT.

package notion

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

// Page defines a model
type Page struct {
	Object         string           `json:"object,omitzero"`
	ID             uuid.UUID        `json:"id,omitzero"`
	CreatedTime    time.Time        `json:"created_time,omitzero"`
	LastEditedTime time.Time        `json:"last_edited_time,omitzero"`
	CreatedBy      PageCreatedBy    `json:"created_by"`
	LastEditedBy   PageLastEditedBy `json:"last_edited_by"`
	Cover          PageCover        `json:"cover"`
	Icon           PageIcon         `json:"icon"`
	Parent         PageParent       `json:"parent"`
	Archived       bool             `json:"archived,omitzero"`
	InTrash        bool             `json:"in_trash,omitzero"`
	Properties     PageProperties   `json:"properties"`
	URL            url.URL          `json:"url,omitzero"`
	PublicURL      url.URL          `json:"public_url,omitzero"`
	RequestID      uuid.UUID        `json:"request_id,omitzero"`
}

// PageCover defines a model
type PageCover struct {
	Type     string            `json:"type,omitzero"`
	External PageCoverExternal `json:"external"`
}

// PageCoverExternal defines a model
type PageCoverExternal struct {
	URL url.URL `json:"url,omitzero"`
}

// PageCreatedBy defines a model
type PageCreatedBy struct {
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

// PageIcon defines a model
type PageIcon struct {
	Type  string `json:"type,omitzero"`
	Emoji string `json:"emoji,omitzero"`
}

// PageLastEditedBy defines a model
type PageLastEditedBy struct {
	Object string    `json:"object,omitzero"`
	ID     uuid.UUID `json:"id,omitzero"`
}

// PageParent defines a model
type PageParent struct {
	Type   string    `json:"type,omitzero"`
	PageID uuid.UUID `json:"page_id,omitzero"`
}

// PageProperties defines a model
type PageProperties struct {
	Title PagePropertiesTitle `json:"title"`
}

// PagePropertiesTitle defines a model
type PagePropertiesTitle struct {
	ID    string                   `json:"id,omitzero"`
	Type  string                   `json:"type,omitzero"`
	Title PagePropertiesTitleTitle `json:"title,omitempty"`
}

// PagePropertiesTitleTitle defines a model
type PagePropertiesTitleTitle []PagePropertiesTitleTitleItems

// PagePropertiesTitleTitleItems defines a model
type PagePropertiesTitleTitleItems struct {
	Type        string                                   `json:"type,omitzero"`
	Text        PagePropertiesTitleTitleItemsText        `json:"text"`
	Annotations PagePropertiesTitleTitleItemsAnnotations `json:"annotations"`
	PlainText   string                                   `json:"plain_text,omitzero"`
	Href        struct{}                                 `json:"href"`
}

// PagePropertiesTitleTitleItemsAnnotations defines a model
type PagePropertiesTitleTitleItemsAnnotations struct {
	Bold          bool   `json:"bold,omitzero"`
	Italic        bool   `json:"italic,omitzero"`
	Strikethrough bool   `json:"strikethrough,omitzero"`
	Underline     bool   `json:"underline,omitzero"`
	Code          bool   `json:"code,omitzero"`
	Color         string `json:"color,omitzero"`
}

// PagePropertiesTitleTitleItemsText defines a model
type PagePropertiesTitleTitleItemsText struct {
	Content string   `json:"content,omitzero"`
	Link    struct{} `json:"link"`
}

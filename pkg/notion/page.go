package notion

// Title returns the plain text of the page's title property.
func (p *Page) Title() string {
	for _, prop := range p.Properties {
		if t := prop.title(); t != nil {
			return t.PlainText()
		}
	}

	return ""
}

// title returns the property's title, or nil if it is not a title property.
func (v *PropertyValue) title() RichTexts {
	sa := v.PropertyValueAllOf1.SimpleOrArrayPropertyValueResponse
	if sa == nil || sa.ArrayBasedPropertyValueResponse == nil ||
		sa.ArrayBasedPropertyValueResponse.TitleArrayBasedPropertyValueResponse == nil {
		return nil
	}

	return sa.ArrayBasedPropertyValueResponse.TitleArrayBasedPropertyValueResponse.Title
}

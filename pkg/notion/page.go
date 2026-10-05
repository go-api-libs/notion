package notion

// Title returns the plain text of the page's title property.
func (p *Page) Title() string {
	for _, prop := range p.Properties {
		if t := prop.PropertyValueAllOf1.TitleArrayBasedPropertyValueResponse; t != nil {
			return t.Title.PlainText()
		}
	}

	return ""
}

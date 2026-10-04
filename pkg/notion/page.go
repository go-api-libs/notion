package notion

// Title returns the plain text of the page's title property.
func (p *Page) Title() string {
	for _, prop := range p.Properties {
		if prop.Type == PropertyValueTypeTitle {
			return prop.Title.PlainText()
		}
	}

	return ""
}

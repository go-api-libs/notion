package notion

// Title returns the plain text of the page's title property.
func (p *Page) Title() string {
	for _, prop := range p.Properties {
		if t := prop.Title; t != nil {
			return t.PlainText()
		}
	}

	return ""
}

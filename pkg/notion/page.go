package notion

// Title returns the page title.
func (p Page) Title() string {
	return p.Properties.title()
}

// Title returns the title of the page.
func (props PropertyValues) title() string {
	for _, prop := range props {
		if prop.Type == PropertyTypeTitle {
			return prop.Title.plainText()
		}
	}

	return ""
}

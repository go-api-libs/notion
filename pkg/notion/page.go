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

// Page returns the result if it is a page, in full, or else nil.
func (r PageOrDataSource) Page() *Page {
	if r.PageOrPartial == nil {
		return nil
	}

	return r.PageOrPartial.Page
}

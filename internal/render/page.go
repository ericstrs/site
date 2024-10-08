package render

import "os"

type Page struct {
	Title   string
	Content []byte
}

// LoadPage loads a page
func LoadPage(title, path string) (*Page, error) {
	md, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	html, err := markdownToHTML(md)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Content: html}, nil
}

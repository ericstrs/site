package render

import (
	"bytes"
	"html/template"
)

var (
	tmplDir   = "public/templates"
	templates = template.Must(template.ParseFS(content, tmplDir+"/index.html", tmplDir+"/notfound.html"))
)

// RenderTemplate renders the specified HTML template and returns it
func RenderTemplate(tmpl string, p *Page) ([]byte, error) {
	data := struct {
		Title   string
		Content string
	}{
		Title:   p.Title,
		Content: string(p.Content),
	}

	buff := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buff, tmpl+".html", data); err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

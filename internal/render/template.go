package render

import (
	"bytes"
	"html/template"
)

var (
	tmplDir   = "public/templates"
	homePath  = tmplDir + "/index.html" // path to the home html page
	nfPath    = tmplDir + "/notfound.html"
	aboutPath = tmplDir + "/about.html"
	notesPath = tmplDir + "/notes.html"
	notePath  = tmplDir + "/note.html"
	templates = template.Must(template.ParseFS(content, homePath, nfPath, aboutPath, notesPath, notePath))
)

// Template renders the specified HTML template and returns it
func Template(tmpl string, p *Page) ([]byte, error) {
	data := struct {
		Title   string
		Content template.HTML
	}{
		Title:   p.Title,
		Content: template.HTML(string(p.Content)),
	}

	buff := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buff, tmpl+".html", data); err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

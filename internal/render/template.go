package render

import (
	"bytes"
	"html/template"
)

var (
	tmplDir    = "public/templates"
	headPath   = tmplDir + "/head.tmpl"
	headerPath = tmplDir + "/header.tmpl"
	footerPath = tmplDir + "/footer.tmpl"
	homePath   = tmplDir + "/home.html" // path to the home html page
	nfPath     = tmplDir + "/notfound.html"
	aboutPath  = tmplDir + "/about.html"
	notesPath  = tmplDir + "/notes.html"
	notePath   = tmplDir + "/note.html"
	blogsPath  = tmplDir + "/blogs.html"
	blogPath   = tmplDir + "/blog.html"
	templates  = template.Must(template.ParseFS(Public, headPath,
		headerPath, footerPath, homePath, nfPath, aboutPath, notesPath,
		notePath, blogsPath, blogPath))
)

// Template renders the specified HTML template and returns it
func Template(tmpl string, data any) ([]byte, error) {
	buff := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buff, tmpl+".html", data); err != nil {
		return []byte{}, err
	}
	return buff.Bytes(), nil
}

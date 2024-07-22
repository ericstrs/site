package render

import (
	"bufio"
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var (
	tmplDir   = "public/templates"
	homePath  = tmplDir + "/home.html" // path to the home html page
	nfPath    = tmplDir + "/notfound.html"
	aboutPath = tmplDir + "/about.html"
	notesPath = tmplDir + "/notes.html"
	notePath  = tmplDir + "/note.html"
	blogsPath = tmplDir + "/blogs.html"
	blogPath  = tmplDir + "/blog.html"
	templates = template.Must(template.ParseFS(content, homePath, nfPath,
		aboutPath, notesPath, notePath, blogsPath, blogPath))
)

type Content struct {
	Title     string
	Id        string
	UpdatedAt time.Time
}

// Template renders the specified HTML template and returns it
func Template(tmpl string, p *Page) ([]byte, error) {
	recentBlogs, err := recentContent("blogs", 5)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recent blogs: %v", err)
	}

	recentNotes, err := recentContent("notes", 5)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve recent notes: %v", err)
	}

	data := struct {
		Title       string
		Content     template.HTML
		RecentBlogs []Content
		RecentNotes []Content
	}{
		Title:       p.Title,
		Content:     template.HTML(string(p.Content)),
		RecentBlogs: recentBlogs,
		RecentNotes: recentNotes,
	}

	buff := new(bytes.Buffer)
	if err := templates.ExecuteTemplate(buff, tmpl+".html", data); err != nil {
		return []byte{}, err
	}

	return buff.Bytes(), nil
}

// recentContent returns n recent content for the given content type
func recentContent(contentType string, n int) ([]Content, error) {
	var items []Content
	baseDir := "docs/" + contentType

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the root directory's README
		if path == filepath.Join(baseDir, "README.md") {
			return nil
		}

		if !info.IsDir() && filepath.Ext(path) == ".md" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			title := mdTitle(content)
			id := idFromPath(path)
			items = append(items, Content{
				Title:     title,
				Id:        id,
				UpdatedAt: info.ModTime(),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].UpdatedAt.After(items[j].UpdatedAt)
	})

	if len(items) > n {
		items = items[:n]
	}

	return items, nil
}

// mdTitle returns the markdown title for given markdown content.
func mdTitle(content []byte) string {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	if scanner.Scan() {
		title := strings.TrimSpace(scanner.Text())
		title = strings.TrimLeft(title, "# ")
		return title
	}
	return "Untitled"
}

// idFromPath returns the page id for the given path.
func idFromPath(path string) string {
	// Get parent directory name
	dir := filepath.Dir(path)
	// Extract the last part of the path, which should be the ID
	return filepath.Base(dir)
}

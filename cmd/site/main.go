package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"text/template"
	"time"

	chromahtml "github.com/alecthomas/chroma/v2/formatters/html"
	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting/v2"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// content contains static web server content
//
//go:embed public docs
var content embed.FS

var (
	tmplDir   = "public/templates"
	templates = template.Must(template.ParseFS(content, tmplDir+"/index.html", tmplDir+"/notfound.html"))
)

type Page struct {
	Title   string
	Content []byte
}

func main() {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", Logging(http.HandlerFunc(indexHandler)))

	handler := PanicRecovery(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

// indexHandler handles the index end point
func indexHandler(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	p, err := loadPage("docs/README.md")
	if err != nil {
		logger.Error("failed to load index markdown file", "err", err)
		http.Error(w, "error: something went wrong", http.StatusInternalServerError)
		return
	}

	renderTemplate(w, "index", p)
}

// loadPage loads a page
func loadPage(path string) (*Page, error) {
	md, err := content.ReadFile(path)
	if err != nil {
		return nil, err
	}
	title := "ericstrs"
	html, err := markdownToHTML(md)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Content: html}, nil
}

// markdownToHTML converts the given markdown into its HTML
// representation
func markdownToHTML(content []byte) ([]byte, error) {
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
			highlighting.NewHighlighting(
				highlighting.WithStyle("gruvbox"),
				highlighting.WithFormatOptions(
					chromahtml.WithLineNumbers(true),
				),
			),
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
	var buf bytes.Buffer
	if err := md.Convert(content, &buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// renderTemplate renders the specified html template
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	data := struct {
		Title   string
		Content string
	}{
		Title:   p.Title,
		Content: string(p.Content),
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	if err := templates.ExecuteTemplate(w, tmpl+".html", data); err != nil {
		logger.Error("failed to execute html template", "err", err)
		http.Error(w, "error: something went wrong", http.StatusInternalServerError)
	}
}

// Logging is a middleware function for logging requests
func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		start := time.Now()

		next.ServeHTTP(w, r)

		var (
			took    = formatDuration(time.Since(start))
			method  = r.Method
			referer = r.Referer()
			addr    = r.RemoteAddr
			uri     = r.RequestURI
		)

		logger.Info("Request", "method", method, "took", took, "referer", referer, "remote_addr", addr, "uri", uri)
	})
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		us := float64(d.Nanoseconds()) / float64(time.Microsecond)
		return fmt.Sprintf("%.3fµs", us)
	}
	if d < time.Second {
		ms := float64(d.Nanoseconds()) / float64(time.Millisecond)
		return fmt.Sprintf("%.3fms", ms)
	}
	if d < time.Minute {
		s := d.Seconds()
		return fmt.Sprintf("%.3fs", s)
	}
	m := d.Minutes()
	return fmt.Sprintf("%.2fm", m)
}

// PanicRecovery is a middleware function for recovering on panic
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("Server failed", "err", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

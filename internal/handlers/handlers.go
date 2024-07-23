package handlers

import (
	"html/template"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/render"
)

// Home handles the home endpoint
func Home(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			filePath = filepath.Join(cfg.DocsPath, "README.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Error("failed to load home markdown file", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		recentBlogs, err := render.RecentContent("blogs", 5)
		if err != nil {
			slog.Error("failed to retrieve recent blogs", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}
		recentNotes, err := render.RecentContent("notes", 5)
		if err != nil {
			slog.Error("failed to retrieve recent notes", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
			RecentBlogs []render.Content
			RecentNotes []render.Content
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
			RecentBlogs: recentBlogs,
			RecentNotes: recentNotes,
		}

		output, err := render.Template("home", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

// About handles the about endpoint
func About(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			filePath = filepath.Join(cfg.DocsPath, "about.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Error("failed to load about markdown file", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
		}

		output, err := render.Template("about", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

// Notes handles the notes endpoint
func Notes(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			filePath = filepath.Join(cfg.DocsPath, "notes", "README.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Error("failed to load about markdown file", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		notes, err := render.AllContent("notes")
		if err != nil {
			slog.Error("failed to retrieve notes", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
			Notes       []render.Content
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
			Notes:       notes,
		}

		output, err := render.Template("notes", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

// Note handles the note endpoint
func Note(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			idStr    = r.PathValue("id")
			filePath = filepath.Join(cfg.DocsPath, "notes", idStr, "README.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Warn("markdown file not found", "err", err,
				"method", method, "uri", uri,
			)
			http.NotFound(w, r)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
		}

		output, err := render.Template("note", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

// Blogs handles the blogs endpoint
func Blogs(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			filePath = filepath.Join(cfg.DocsPath, "blogs", "README.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Error("failed to load about markdown file", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		blogs, err := render.AllContent("blogs")
		if err != nil {
			slog.Error("failed to retrieve blogs", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
			Blogs       []render.Content
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
			Blogs:       blogs,
		}

		output, err := render.Template("blogs", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

// Blog handles the blog endpoint
func Blog(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			title    = cfg.Title
			idStr    = r.PathValue("id")
			filePath = filepath.Join(cfg.DocsPath, "blogs", idStr, "README.md")

			method = r.Method
			uri    = r.URL.RequestURI()
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			slog.Warn("markdown file not found", "err", err,
				"method", method, "uri", uri,
			)
			http.NotFound(w, r)
			return
		}

		data := struct {
			Nav         []config.NavItem
			Title       string
			Description string
			Content     template.HTML
		}{
			Nav:         cfg.Nav,
			Title:       p.Title,
			Description: cfg.Description,
			Content:     template.HTML(string(p.Content)),
		}

		output, err := render.Template("note", data)
		if err != nil {
			slog.Error("failed to execute html template", "err", err,
				"method", method, "uri", uri,
			)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

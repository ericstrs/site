package handlers

import (
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

		output, err := render.Template("index", p)
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

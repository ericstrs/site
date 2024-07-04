package handlers

import (
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/render"
)

// Home handles the home end point
func Home(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		var (
			title    = cfg.Title
			filePath = filepath.Join(cfg.DocsPath, "README.md")
		)

		p, err := render.LoadPage(title, filePath)
		if err != nil {
			logger.Error("failed to load home markdown file", "err", err)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		output, err := render.Template("index", p)
		if err != nil {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			logger.Error("failed to execute html template", "err", err)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

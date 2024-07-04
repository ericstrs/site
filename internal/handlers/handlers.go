package handlers

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/render"
)

// Home handles the home end point
func Home(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

		var (
			title = cfg.Title
			path  = "docs/README.md"
		)

		p, err := render.LoadPage(title, path)
		if err != nil {
			logger.Error("failed to load index markdown file", "err", err)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		output, err := render.RenderTemplate("index", p)
		if err != nil {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			logger.Error("failed to execute html template", "err", err)
			http.Error(w, "error: something went wrong", http.StatusInternalServerError)
			return
		}

		w.Write(output)
	}
}

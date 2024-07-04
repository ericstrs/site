package server

import (
	"log"
	"net/http"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/handlers"
	"github.com/ericstrs/site/internal/middleware"
)

func Serve() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()

	mux.Handle("GET /{$}", middleware.Logging(handlers.Home(cfg)))

	handler := middleware.PanicRecovery(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

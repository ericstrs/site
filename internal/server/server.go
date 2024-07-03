package server

import (
	"log"
	"net/http"

	"github.com/ericstrs/site/internal/handlers"
	"github.com/ericstrs/site/internal/middleware"
)

func Serve() {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", middleware.Logging(http.HandlerFunc(handlers.Home)))

	handler := middleware.PanicRecovery(mux)

	log.Fatal(http.ListenAndServe(":8080", handler))
}

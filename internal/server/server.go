package server

import (
	"context"
	"crypto/tls"
	"io/fs"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"syscall"
	"time"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/handlers"
	"github.com/ericstrs/site/internal/middleware"
	"github.com/ericstrs/site/internal/render"
)

func Serve() {
	var trace = string(debug.Stack())
	var logLevel = new(slog.LevelVar)

	opts := &slog.HandlerOptions{
		AddSource: true,
		Level:     logLevel,
	}
	h := slog.NewJSONHandler(os.Stdout, opts)
	logger := slog.New(h)
	slog.SetDefault(logger)

	logLevel.Set(slog.LevelInfo)

	cfg, err := config.LoadConfig()
	if err != nil {
		slog.Error("Server failed", "err", err, "trace", trace)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", middleware.LogRequest(handlers.Home(cfg)))
	mux.Handle("GET /about", middleware.LogRequest(handlers.About(cfg)))
	mux.Handle("GET /notes", middleware.LogRequest(handlers.Notes(cfg)))
	mux.Handle("GET /notes/{id}", middleware.LogRequest(handlers.Note(cfg)))
	mux.Handle("GET /blogs", middleware.LogRequest(handlers.Blogs(cfg)))
	mux.Handle("GET /blogs/{id}", middleware.LogRequest(handlers.Blog(cfg)))

	publicFS, err := fs.Sub(render.Public, "public")
	if err != nil {
		log.Fatal(err)
	}
	mux.Handle("/", http.FileServer(http.FS(publicFS)))

	handler := middleware.SecurityHeaders(mux)
	handler = middleware.PanicRecovery(handler)

	portStr := strconv.Itoa(cfg.Port)
	addr := cfg.Host + ":" + portStr
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
		TLSConfig: &tls.Config{
			MinVersion:               tls.VersionTLS13,
			PreferServerCipherSuites: true,
		},
	}

	go func() {
		logger.Info("Server is starting...", "addr", server.Addr)
		if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to serve", "err", err, "trace", trace)
			os.Exit(1)
		}
	}()

	// Channel to listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Blocking until signal received
	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "err", err, "trace", trace)
		os.Exit(1)
	}
	logger.Info("Server gracefully stopped")
}

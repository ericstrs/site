package server

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ericstrs/site/internal/config"
	"github.com/ericstrs/site/internal/handlers"
	"github.com/ericstrs/site/internal/middleware"
)

func Serve() {
	h := slog.NewJSONHandler(os.Stdout, nil)
	logger := slog.New(h)

	cfg, err := config.LoadConfig()
	if err != nil {
		logger.Error("Server failed", "err", err)
		os.Exit(1)
	}

	mux := http.NewServeMux()
	mux.Handle("GET /{$}", middleware.Logging(logger, handlers.Home(logger, cfg)))

	handler := middleware.PanicRecovery(logger, mux)

	portStr := strconv.Itoa(cfg.Port)
	addr := cfg.Host + ":" + portStr
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		logger.Info("Server is starting...", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("Server failed to serve", "err", err)
			os.Exit(1)
		}
	}()

	// Channel to listen for signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGQUIT)

	// Blocking until signal received
	<-stop
	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Error("Server shutdown failed", "err", err)
		os.Exit(1)
	}
	logger.Info("Server gracefully stopped")
}

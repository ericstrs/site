package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Logging middleware function for logging requests
func Logging(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

// PanicRecovery is middleware for recovering from panics in `next` and
// returning a StatusInternalServerError to the client.
func PanicRecovery(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				logger.Error("Server failed", "err", err)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

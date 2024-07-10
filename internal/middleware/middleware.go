package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// LogRequest middleware function for logging requests
func LogRequest(next http.Handler) http.Handler {
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

		slog.Info("Request", "method", method, "took", took, "referer", referer, "remote_addr", addr, "uri", uri)
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
func PanicRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				slog.Error("Server failed", "err", err, "trace", string(debug.Stack()))
			}
		}()
		next.ServeHTTP(w, r)
	})
}

// SecurityHeaders middleware function for logging requests
func SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		csp := []string{
			"default-src 'self'",
			"form-action 'self'",
			"object-src 'none'",
			"frame-ancestors 'none'",
			"upgrade-insecure-requests",
			"block-all-mixed-content",
		}

		w.Header().Set("Content-Security-Policy", strings.Join(csp, "; "))
		w.Header().Set("Strict-Transport-Security", "max-age=604800; includeSubDomains")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-XSS-Protection", "0")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Cross-Origin-Embedder-Policy", "require-corp")
		w.Header().Set("Cross-Origin-Opener-Policy", "same-origin")
		w.Header().Set("Cross-Origin-Resource-Policy", "same-origin")
		w.Header().Set("Cache-Control", "no-store, max-age=0")

		w.Header().Set("Server", "")

		next.ServeHTTP(w, r)
	})
}

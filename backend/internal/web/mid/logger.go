package mid

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

var requestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "request_duration_seconds",
	Help:    "Time (in seconds) spent serving HTTP requests.",
	Buckets: prometheus.DefBuckets,
}, []string{"method", "route", "status_code"})

func init() { // nolint:gochecknoinits
	prometheus.MustRegister(requestDuration)
}

type spy struct {
	http.ResponseWriter
	status int
}

func (s *spy) WriteHeader(status int) {
	s.status = status
	s.ResponseWriter.WriteHeader(status)
}

const frontendVersionHeader = "X-Frontend-Version"

func extractFrontendVersionHeader(r *http.Request) string {
	if v := r.Header.Get(frontendVersionHeader); v != "" {
		return v
	}
	return ""
}

// Logger is a middleware that logs the request and response.
//
// - General request information (Method, Path, Request ID, Duration, Status)
// - Warning when the frontend version does not match the backend version.
// - Request duration is recorded in Prometheus with labels (Method, Route, Status Code).
func Logger(l zerolog.Logger, version string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			now := time.Now()

			frontendVersionHeaderValue := extractFrontendVersionHeader(r)
			if version != "" && frontendVersionHeaderValue != "" && frontendVersionHeaderValue != version {
				l.Warn().
					Ctx(r.Context()).
					Str("frontend_version", frontendVersionHeaderValue).
					Str("backend_version", version).
					Msg("frontend version mismatch")
			}

			l.Info().
				Ctx(r.Context()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Msg("->")

			s := &spy{ResponseWriter: w}
			h.ServeHTTP(s, r)

			// routePattern only available after middleware execution
			routePattern := chi.RouteContext(r.Context()).RoutePattern()
			requestDuration.
				WithLabelValues(r.Method, routePattern, strconv.Itoa(s.status)).
				Observe(time.Since(now).Seconds())

			l.Info().Ctx(r.Context()).
				Str("method", r.Method).
				Str("path", r.URL.Path).
				Int64("dur_ms", time.Since(now).Milliseconds()).
				Int("status", s.status).Msg("<-")
		})
	}
}

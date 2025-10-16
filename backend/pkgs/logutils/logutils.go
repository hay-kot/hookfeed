// Package logutils provides some small utilities for logging.
// specifically around setting up a logger with a given log level.
package logutils

import (
	"io"
	"net/http"
	"os"

	"github.com/rs/zerolog"
)

// Factory returns a new logger with the given level and style.
// This is used to simplify the creation of a logger in the main
// function of an application and keep logs consistent across the
// applications.
//
// The level parameter can be one of the following:
//   - debug
//   - info
//   - warn
//   - error
//   - fatal
//
// The style parameter can be one of the following:
//   - console
//   - json
func Factory(level string, style string) (zerolog.Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return zerolog.Logger{}, err
	}

	var logWriter io.Writer
	switch style {
	case "console":
		logWriter = zerolog.ConsoleWriter{Out: os.Stdout}
	default:
		logWriter = os.Stdout
	}

	l := zerolog.New(logWriter).
		With().
		Caller().    // adds the file and line number of the caller
		Timestamp(). // adds a timestamp to each log line
		Logger().
		Level(lvl)

	return l, nil
}

// Handle consumes the error provided and logs it to the logger with a stack trace for
// debugging purposes.
func Handle(log zerolog.Logger, err error, args ...any) {
	if err != nil {
		log.Warn().Err(err).Stack().Msg("error handled")
	}
}

type RoundTripper struct {
	ref  string
	l    zerolog.Logger
	next http.RoundTripper
}

func NewRoundTripper(ref string, l zerolog.Logger, next http.RoundTripper) RoundTripper {
	return RoundTripper{
		ref:  ref,
		l:    l,
		next: next,
	}
}

// RoundTrip implements http.RoundTripper.
func (rt RoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	ctx := r.Context()

	rt.l.Debug().
		Ctx(ctx).
		Str("ref", rt.ref).
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Msg("request")

	resp, err := rt.next.RoundTrip(r)
	if err != nil {
		rt.l.Error().
			Ctx(ctx).
			Str("ref", rt.ref).
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Err(err).
			Msg("request failed")
		return nil, err
	}

	rt.l.Debug().
		Ctx(ctx).
		Str("ref", rt.ref).
		Str("method", r.Method).
		Str("url", r.URL.String()).
		Int("status", resp.StatusCode).
		Msg("request completed")

	return resp, nil
}

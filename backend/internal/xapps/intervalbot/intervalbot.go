package intervalbot

import (
	"context"

	"github.com/rs/zerolog"
)

type IntervalBot struct {
	l zerolog.Logger
}

func New(l zerolog.Logger) *IntervalBot {
	return &IntervalBot{
		l: l.With().Str("service", "interval_bot").Logger(),
	}
}

func (ib *IntervalBot) Start(ctx context.Context) error {
	ib.l.Info().Msg("starting service")
	<-ctx.Done()
	ib.l.Info().Msg("stopping service")

	return nil
}

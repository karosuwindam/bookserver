package config

import (
	"context"
	"log/slog"
	"os"

	"github.com/m-mizutani/clog"
	"github.com/newrelic/go-agent/v3/integrations/logcontext-v2/nrslog"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func Init_logConfig() {
	if NewRelic.App != nil {
		handler := nrslog.TextHandler(
			NewRelic.App, os.Stdout, &slog.HandlerOptions{
				AddSource: true,
			},
		)
		logger := slog.New(handler)
		slog.SetDefault(logger)
		txnCtx := newrelic.NewContext(context.Background(), nil)
		slog.InfoContext(txnCtx, "Log configuration initialized")
	} else {
		handler := clog.New(
			clog.WithColor(true),
			clog.WithSource(true),
		)
		logger := slog.New(handler)
		slog.SetDefault(logger)

	}

}

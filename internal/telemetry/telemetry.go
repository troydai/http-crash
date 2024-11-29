package telemetry

import (
	"log/slog"
	"os"

	"go.uber.org/fx"
)

var Module = fx.Provide(ProvideTelemetry)

type Result struct {
	fx.Out

	Logger *slog.Logger
}

func ProvideTelemetry() Result {
	return Result{
		Logger: slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}
}

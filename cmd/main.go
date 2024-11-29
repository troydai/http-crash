package main

import (
	"go.uber.org/fx"

	"github.com/troydai/http-crash/internal/http"
	"github.com/troydai/http-crash/internal/settings"
	"github.com/troydai/http-crash/internal/telemetry"
)

func main() {
	fx.New(
		http.Module,
		telemetry.Module,
		settings.Module,
	).Run()
}

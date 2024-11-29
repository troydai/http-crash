package settings

import (
	"fmt"

	"github.com/caarlos0/env/v11"
	"go.uber.org/fx"
)

var Module = fx.Provide(ProvideSettings)

type Values struct {
	CrashFrequency uint64 `env:"HTTP_CRASH_FREQUENCY" envDefault:"10"` // Set to zero to disable
	LatencyMean    uint64 `env:"HTTP_LATENCY_MEAN" envDefault:"50"`
	LatencyStddev  uint64 `env:"HTTP_LATENCY_STDDEV" envDefault:"25"`
}

func ProvideSettings() (*Values, error) {
	var s Values
	if err := env.Parse(&s); err != nil {
		return nil, fmt.Errorf("error parsing environment variables: %w", err)
	}

	return &s, nil
}

package http

import (
	"math/rand"
	"time"
)

type normalDistribution struct {
	mean   float64
	stddev float64
}

func newNormalDistribution(mean, stddev uint64) *normalDistribution {
	return &normalDistribution{mean: float64(mean), stddev: float64(stddev)}
}

func (n *normalDistribution) latency() time.Duration {
	ms := n.mean + rand.NormFloat64()*n.stddev
	if ms < 0 {
		return 0
	}

	return time.Duration(ms) * time.Millisecond
}

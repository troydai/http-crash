package http

import (
	"context"
	"log/slog"
	"net/http"
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	"github.com/troydai/http-crash/internal/settings"
)

var Module = fx.Options(
	fx.Provide(ProvideServer),
	fx.Invoke(StartHTTPServer),
)

func ProvideServer(s *settings.Values, l *slog.Logger) *Server {
	return &Server{
		counter:    &atomic.Uint64{},
		frequency:  s.CrashFrequency,
		logger:     l,
		latencyGen: newNormalDistribution(s.LatencyMean, s.LatencyStddev),
	}
}

type Server struct {
	counter    *atomic.Uint64
	frequency  uint64
	logger     *slog.Logger
	latencyGen *normalDistribution
}

func (s *Server) HandleHTTP(w http.ResponseWriter, r *http.Request) {
	newCount := s.counter.Add(1)
	latency := s.latencyGen.latency()

	s.logger.Info("inbound",
		"method", r.Method,
		"url", r.URL.Path,
		"host", r.Host,
		"proto", r.Proto,
		"agent", r.UserAgent(),
		"peer", r.RemoteAddr,
		"counter", newCount,
		"latency", latency,
	)

	time.Sleep(s.latencyGen.latency())

	if s.frequency != 0 && newCount%s.frequency == 0 {
		s.logger.Warn("crash", "counter", newCount)
		panic("crash as requested")
	}

	s.logger.Info("outbound", "status", 200)
	w.WriteHeader(200)
}

func (s *Server) HandleHealthCheck(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("healh-check", "method", r.Method, "url", r.URL.Path, "host", r.Host, "proto", r.Proto, "agent", r.UserAgent(), "peer", r.RemoteAddr, "counter", s.counter.Load())
	w.WriteHeader(200)
}

func StartHTTPServer(s *Server, lifecycle fx.Lifecycle, l *slog.Logger) {
	mux := http.NewServeMux()
	mux.HandleFunc("/_health", s.HandleHealthCheck)
	mux.HandleFunc("/", s.HandleHTTP)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	lifecycle.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			l.Info("starting HTTP server", "addr", server.Addr)
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					l.Error("server error", "error", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			l.Info("stopping HTTP server")
			return server.Shutdown(ctx)
		},
	})
}

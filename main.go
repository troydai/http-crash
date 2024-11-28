package main

import (
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	s.logger.Info("inbound", "method", r.Method, "url", r.URL.Path, "host", r.Host, "proto", r.Proto, "agent", r.UserAgent(), "peer", r.RemoteAddr, "counter", s.counter.Load())

	if newValue := s.counter.Add(1); newValue%s.frequency == 0 {
		s.logger.Warn("crash", "counter", newValue)
		panic("crash after 10 requests")
	}

	s.logger.Info("outbound", "status", 200)
	w.WriteHeader(200)
}

type server struct {
	counter   *atomic.Uint64
	frequency uint64
	logger    *slog.Logger
}

func main() {
	ev := os.Getenv("HTTP_CRASH_FREQUENCY")
	frequency, err := strconv.ParseUint(ev, 10, 64)
	if err != nil {
		frequency = 10
	}

	s := &server{
		counter:   &atomic.Uint64{},
		frequency: frequency,
		logger:    slog.New(slog.NewTextHandler(os.Stdout, nil)),
	}

	http.HandleFunc("/", s.handler)
	http.ListenAndServe(":8080", nil)
}

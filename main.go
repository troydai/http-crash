package main

import (
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	if newValue := s.counter.Add(1); newValue%s.frequency == 0 {
		panic("crash after 10 requests")
	}

	w.WriteHeader(200)
}

type server struct {
	counter   *atomic.Uint64
	frequency uint64
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
	}

	http.HandleFunc("/", s.handler)
	http.ListenAndServe(":8080", nil)
}

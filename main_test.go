package main

import (
	"io"
	"log/slog"
	"net/http"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	s := &server{
		counter:   &atomic.Uint64{},
		frequency: 10,
		logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
	s.handler(&fakeResponseWriter{}, randomRequest())
	assert.Equal(t, uint64(1), s.counter.Load())
}

func randomRequest() *http.Request {
	req, err := http.NewRequest("GET", "http://example.com", nil)
	if err != nil {
		panic(err)
	}
	return req
}

type fakeResponseWriter struct{}

func (frw *fakeResponseWriter) Header() http.Header {
	return http.Header{}
}

func (frw *fakeResponseWriter) Write(data []byte) (int, error) {
	return len(data), nil
}

func (frw *fakeResponseWriter) WriteHeader(statusCode int) {}

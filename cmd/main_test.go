package main

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/gabrielpgava/stress-test-fullcycle/internal/storage"
)

func TestValidateFlags(t *testing.T) {
	cases := []struct {
		name        string
		url         string
		requests    int
		concurrency int
		wantErr     bool
	}{
		{name: "empty url", url: "", requests: 1, concurrency: 1, wantErr: true},
		{name: "requests zero", url: "http://example.com", requests: 0, concurrency: 1, wantErr: true},
		{name: "concurrency zero", url: "http://example.com", requests: 1, concurrency: 0, wantErr: true},
		{name: "concurrency greater", url: "http://example.com", requests: 1, concurrency: 2, wantErr: true},
		{name: "valid", url: "http://example.com", requests: 5, concurrency: 2, wantErr: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateFlags(tc.url, tc.requests, tc.concurrency)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
		})
	}
}

func TestRunLoadTotals(t *testing.T) {
	storage.Reset()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := server.Client()
	client.Timeout = 2 * time.Second

	runLoad(server.URL, 10, 3, client)

	counts := storage.GetStatusCounts()
	if counts[200] != 10 {
		t.Fatalf("expected 10 status 200, got %d", counts[200])
	}
	if storage.GetErrorCount() != 0 {
		t.Fatalf("expected 0 errors, got %d", storage.GetErrorCount())
	}
}

func TestRunLoadStatusDistribution(t *testing.T) {
	storage.Reset()

	var counter int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt64(&counter, 1)
		switch {
		case n <= 5:
			w.WriteHeader(http.StatusOK)
		case n <= 9:
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	defer server.Close()

	client := server.Client()
	client.Timeout = 2 * time.Second

	runLoad(server.URL, 12, 5, client)

	counts := storage.GetStatusCounts()
	if counts[200] != 5 {
		t.Fatalf("expected 5 status 200, got %d", counts[200])
	}
	if counts[404] != 4 {
		t.Fatalf("expected 4 status 404, got %d", counts[404])
	}
	if counts[500] != 3 {
		t.Fatalf("expected 3 status 500, got %d", counts[500])
	}
	if storage.GetErrorCount() != 0 {
		t.Fatalf("expected 0 errors, got %d", storage.GetErrorCount())
	}
}

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kagirinay/MetricsCollector.git/internal/handlers"
	"github.com/kagirinay/MetricsCollector.git/internal/store"
)

func TastUpdateHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantStatus    int
		expectGauge   bool
		expectCounter bool
	}{
		{
			name:        "valid gauge",
			url:         "/update/gauge/Alloc/123.5",
			wantStatus:  http.StatusOK,
			expectGauge: true,
		},
		{
			name:          "valid counter",
			url:           "/update/counter/Requests/10",
			wantStatus:    http.StatusOK,
			expectCounter: true,
		},
		{
			name:       "unknown type",
			url:        "/update/unknown/a/5",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "no name",
			url:        "/update/gauge//10",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "bad value",
			url:        "/update/counter/Bad/notInt",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := store.NewMemStorage()
			h := handlers.Update(mem)
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			h(w, req)
			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			if tt.expectGauge {
				if v, ok := mem.GetGauge("Alloc"); !ok || v != 123.5 {
					t.Fatalf("gauge not store correctly, got %v ok=%v", v, ok)
				}
			}
		})
	}
}

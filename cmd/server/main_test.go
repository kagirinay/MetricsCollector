package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
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

func TestGetMetricHandler(t *testing.T) {
	tests := []struct {
		name       string
		url        string
		wantStatus int
		wantBody   string
	}{
		{
			name:       "existing gauge",
			url:        "/value/gauge/Alloc",
			wantStatus: http.StatusOK,
			wantBody:   "123,5",
		},
		{
			name:       "non-existing gauge",
			url:        "/value/gauge/NonExisting",
			wantStatus: http.StatusNotFound,
		},
		{
			name:       "invalid type",
			url:        "/value/invalid/Alloc",
			wantStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := store.NewMemStorage()
			mem.UpdateGauge("Alloc", 123.5)
			req := httptest.NewRequest(http.MethodGet, tt.url, nil)
			w := httptest.NewRecorder()
			// Используем CHI роутер
			r := chi.NewRouter()
			r.Get("/value/{type}/{name}", handlers.GetMetric(mem))
			r.ServeHTTP(w, req)
			resp := w.Result()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("Ожидаемый статус: %d, получаем %q", tt.wantStatus, resp.StatusCode)
			}
			if tt.wantBody != "" {
				body, _ := io.ReadAll(resp.Body)
				if string(body) != tt.wantBody {
					t.Errorf("Желаемое значение тела: %q, получаем %q", tt.wantBody, string(body))
				}
			}
		})
	}
}

func TestHomeHandlers(t *testing.T) {
	mem := store.NewMemStorage()
	mem.UpdateGauge("Alloc", 123.5)
	mem.UpdateCounter("Requests", 10)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	r := chi.NewRouter()
	r.Get("/", handlers.Home(mem))
	r.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый статус: 200, получаемый статус: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	if !strings.Contains(string(body), "Alloc") || !strings.Contains(string(body), "Requests") {
		t.Errorf("Запрашиваемое тело на содержит ожидаемые метрики")
	}
}

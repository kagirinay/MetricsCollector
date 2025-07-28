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
	"github.com/kagirinay/MetricsCollector.git/models"
)

func TestUpdateHandler(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		wantStatus    int
		expectGauge   bool
		expectCounter bool
		metricName    string
		metricValue   float64
	}{
		{
			name:        "valid gauge",
			url:         "/update/gauge/Alloc/123.5",
			wantStatus:  http.StatusOK,
			expectGauge: true,
			metricName:  "Alloc",
			metricValue: 123.5,
		},
		{
			name:          "valid counter",
			url:           "/update/counter/Requests/10",
			wantStatus:    http.StatusOK,
			expectCounter: true,
			metricName:    "Requests",
			metricValue:   10,
		},
		{
			name:       "unknown type",
			url:        "/update/unknown/TestMetric/5",
			wantStatus: http.StatusNotImplemented,
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
		{
			name:        "very long metric name",
			url:         "/update/gauge/" + strings.Repeat("a", 1000) + "/42.42",
			wantStatus:  http.StatusOK,
			expectGauge: true,
			metricName:  strings.Repeat("a", 1000),
			metricValue: 42.42,
		},
		{
			name:        "special characters in name",
			url:         "/update/gauge/Test%21%40%23%24%25%5E%26%2A%28%29_%2B/42.42",
			wantStatus:  http.StatusOK,
			expectGauge: true,
			metricName:  "Test!@#$%^&*()_+",
			metricValue: 42.42,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mem := store.NewMemStorage()
			// Используем chi роутер
			r := chi.NewRouter()
			r.Post("/update/{type}/{name}/{value}", handlers.Update(mem))
			req := httptest.NewRequest(http.MethodPost, tt.url, nil)
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)
			resp := w.Result()
			defer resp.Body.Close()
			if resp.StatusCode != tt.wantStatus {
				t.Errorf("expected %d, got %d", tt.wantStatus, resp.StatusCode)
			}
			// Проверка сохранения значений
			if tt.expectGauge {
				if v, ok := mem.GetGauge(tt.metricName); !ok || v != models.Gauge(tt.metricValue) {
					t.Errorf("Значения корректно не сохраняются, got %v ok=%v", v, ok)
				}
			}
			if tt.expectCounter {
				if v, ok := mem.GetCounter(tt.metricName); !ok || v != models.Counter(int64(tt.metricValue)) {
					t.Errorf("Counter не сохраняется корректно, got %v ok=%v", v, ok)
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
			wantBody:   "123.5",
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
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ожидаемый статус: 200, получаемый статус: %d", resp.StatusCode)
	}
	body, _ := io.ReadAll(resp.Body)
	bodyStr := string(body)
	// Проверяем наличие ожидаемых метрик в HTML
	if !strings.Contains(bodyStr, "Alloc") {
		t.Errorf("Тело ответа не содержит метрику 'Alloc'")
	}
	if !strings.Contains(bodyStr, "Requests") {
		t.Errorf("Тело ответа не содержит метрику 'Requests'")
	}
	// Дополнительные проверки структуры HTML
	if !strings.Contains(bodyStr, "<table>") {
		t.Errorf("Тело ответа не содержит HTML таблицу")
	}
	if !strings.Contains(bodyStr, "Gauges") || !strings.Contains(bodyStr, "Counters") {
		t.Errorf("Тело ответа не содержит разделы метрик")
	}
}

package handlers

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kagirinay/MetricsCollector.git/internal/store"
	"github.com/kagirinay/MetricsCollector.git/models"
)

// Update возвращает http.HandlerFunc, замыкающий хранилище.
func Update(s store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Извлекаем параметры из URL
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		rawValue := chi.URLParam(r, "value")
		// Декодируем имя метрики из URL-формата
		name, err := url.PathUnescape(metricName)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		// Валидация
		if name == "" || metricName == "" || rawValue == "" {
			w.WriteHeader(http.StatusNotFound)

			return
		}
		switch models.MetricType(metricType) {
		case models.GaugeType:
			err = handleGauge(s, name, rawValue)
		case models.CounterType:
			err = handleCounter(s, name, rawValue)
		default:
			w.WriteHeader(http.StatusNotImplemented)

			return
		}
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handleGauge(s store.Storage, name, raw string) error {
	f, err := strconv.ParseFloat(raw, 64)
	if err != nil {
		return err
	}
	s.UpdateGauge(name, models.Gauge(f))
	return nil
}

func handleCounter(s store.Storage, name, raw string) error {
	i, err := strconv.ParseInt(raw, 10, 64)
	if err != nil {
		return err
	}
	s.UpdateCounter(name, models.Counter(i))
	return nil
}

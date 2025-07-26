package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kagirinay/MetricsCollector.git/internal/store"
	"github.com/kagirinay/MetricsCollector.git/models"
)

// GetMetric возвращает значение метрики
func GetMetric(s store.Storage) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Извлекаем параметры из URL
		metricType := chi.URLParam(r, "type")
		metricName := chi.URLParam(r, "name")
		switch models.MetricType(metricType) {
		case models.GaugeType:
			if value, ok := s.GetGauge(metricName); ok {
				w.Write([]byte(strconv.FormatFloat(float64(value), 'f', -1, 64)))

				return
			}
		case models.CounterType:
			if value, ok := s.GetCounter(metricName); ok {
				w.Write([]byte(strconv.FormatInt(int64(value), 10)))

				return
			}
		default:
			http.Error(w, "Неверный тип метрик", http.StatusBadRequest)

			return
		}
		http.Error(w, "Метрика не найдена", http.StatusNotFound)
	}
}

package handlers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/kagirinay/MetricsCollector.git/internal/store"
	"github.com/kagirinay/MetricsCollector.git/models"
)

// Update возвращает http.HandlerFunc, замыкающий хранилище.
func Update(s store.Storage) http.HandlerFunc {
	const prefix = "/update/"
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		if !strings.HasPrefix(r.URL.Path, prefix) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, prefix), "/")
		if len(parts) != 3 || parts[1] == "" {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		typ, name, val := parts[0], parts[1], parts[2]
		var err error
		switch models.MetricType(typ) {
		case models.GaugeType:
			err = handleGauge(s, name, val)
		case models.CounterType:
			err = handleCounter(s, name, val)
		default:
			err = errors.New("Неизвестный тип")
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

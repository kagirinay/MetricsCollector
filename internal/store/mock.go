package store

import "github.com/kagirinay/MetricsCollector.git/models"

type Mock struct {
	UpdateGaugeFn   func(name string, value models.Gauge)
	UpdateCounterFn func(name string, delta models.Counter)
}

func (m *Mock) UpdateGauge(name string, value models.Gauge) {
	if m.UpdateGaugeFn != nil {
		m.UpdateGaugeFn(name, value)
	}
}
func (m *Mock) UpdateCounter(name string, delta models.Counter) {
	if m.UpdateCounterFn != nil {
		m.UpdateCounterFn(name, delta)
	}
}

// Методы необходимы только хранилищу. Хэндлеру не нужны. Пустые.
func (m *Mock) GetGauge(string) (models.Gauge, bool)     { return 0, false }
func (m *Mock) GetCounter(string) (models.Counter, bool) { return 0, false }

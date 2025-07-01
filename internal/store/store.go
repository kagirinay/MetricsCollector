package store

import (
	"sync"

	"github.com/kagirinay/MetricsCollector.git/models"
)

// Store - контракт, который должны удовлетворять любые хранилища.
// Handlers и агенты знают только об этом интерфейсе.
type Storage interface {
	UpdateGauge(name string, value models.Gauge)
	UpdateCounter(name string, delta models.Counter)
	GetGauge(name string) (models.Gauge, bool)
	GetCounter(name string) (models.Counter, bool)
}

type MemStorage struct {
	gauges   map[string]models.Gauge
	counters map[string]models.Counter
	mu       sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauges:   make(map[string]models.Gauge),
		counters: make(map[string]models.Counter),
	}
}

func (m *MemStorage) UpdateGauge(name string, value models.Gauge) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauges[name] = value
}

func (m *MemStorage) UpdateCounter(name string, delta models.Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += delta
}

func (m *MemStorage) GetGauge(name string) (models.Gauge, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.gauges[name]
	return v, ok
}

func (m *MemStorage) GetCounter(name string) (models.Counter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.counters[name]
	return v, ok
}

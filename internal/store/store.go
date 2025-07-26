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
	GetAllGauge() map[string]models.Gauge
	GetAllCounter() map[string]models.Counter
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

func (m *MemStorage) GetAllGauge() map[string]models.Gauge {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Необходимо создать копию для безопасности и во избежание гонок.
	copy := make(map[string]models.Gauge, len(m.gauges))
	for k, v := range m.gauges {
		copy[k] = v
	}

	return copy
}

func (m *MemStorage) GetAllCounter() map[string]models.Counter {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// Необходимо создать копию для безопасности и во избежание гонок.
	copy := make(map[string]models.Counter, len(m.counters))
	for k, v := range m.counters {
		copy[k] = v
	}

	return copy
}

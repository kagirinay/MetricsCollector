package main

import "sync"

// MetricType описывает допустимые строковые идентификаторы типов метрик.
type MetricType string

const (
	GaugeType   MetricType = "gauge"   // Плавающая метрика со значением - float64
	CounterType MetricType = "counter" // Счётчик со значением - int64
)

// Gauge и Counter определены отдельными алиасами
// только для лучшей читаемости кода.
type (
	Gauge   float64
	Counter int64
)

// Storage задаёт набор операций.
type Storage interface {
	UpdateGauge(name string, value Gauge)     // Сохраняет или перезаписывает Gauge
	UpdateCounter(name string, delta Counter) // Увеличивает Counter на delta
	GetGauge(name string) (Gauge, bool)       // Получить gauge (значение, флаг существования)
	GetCounter(name string) (Counter, bool)   // Получить Counter (значение, флаг существования)
}

// MemStorage - конкретный тип, удовлетворяющий интерфейсу Storage
type MemStorage struct {
	gauge    map[string]Gauge   // Карта : имя - значение gauge
	counters map[string]Counter // Карта: имя - значение counter
	mu       sync.RWMutex       // Обеспечивает конкурентную безопасностьы
}

// NewMemStorage конструирует корректно инициализированное хранилище
func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:    make(map[string]Gauge),
		counters: make(map[string]Counter),
	}
}

// UpdateGauge перезаписывает значение gauge
func (m *MemStorage) UpdateGauge(name string, value Gauge) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.gauge[name] = value
}

// UpdateCounter увеличивает существующий counter на delta
// Если счётчик видим впервые, он считается равным 0
func (m *MemStorage) UpdateCounter(name string, delta Counter) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.counters[name] += delta
}

// GetGauge и GetCounter пока не нужны серверу, но пригодятся для тестирования
func (m *MemStorage) GetGauge(name string) (Gauge, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.gauge[name]
	return v, ok
}

func (m *MemStorage) GetCounter(name string) (Counter, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	v, ok := m.counters[name]
	return v, ok
}

func main() {}

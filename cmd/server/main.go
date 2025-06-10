package main

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

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
	mu       sync.RWMutex       // Обеспечивает конкурентную безопасность
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

// updateHandler генерирует http.HandlerFunc, "замыкая" в себе ссылку на хранилище.
// Это упрощает тестирование.
func updateHandler(storage Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Валидируем метод.
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		// Парсим путь вида /update/<type>/<name>/<value>
		const prefix = "/update"
		if !strings.HasPrefix(r.URL.Path, prefix) {
			// Если по каким-то причинам нас повесили не на тот префикс -
			// для надёжности возвращаем 404
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Обрезаем "/update/" и разбиваем по "/"
		parts := strings.Split(strings.TrimPrefix(r.URL.Path, prefix), "/")
		if len(parts) != 3 {
			// если недостаёт компонентов маршрута
			w.WriteHeader(http.StatusNotFound)
			return
		}
		metricType, name, valueStr := parts[0], parts[1], parts[2]
		if name == "" {
			// если метрика пустая
			w.WriteHeader(http.StatusNotFound)
			return
		}
		// Пробуем сконвертировать value и обновить хранилище
		switch MetricType(metricType) {
		case GaugeType:
			if err := handleGauge(storage, name, valueStr); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
		case CounterType:
			if err := handleCounter(storage, name, valueStr); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			// если приходит неизвестный тип метрики
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

// handleGauge конвертирует строку в float64 и кладёт в сторедж.
func handleGauge(storage Storage, name, val string) error {
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return errors.New("invalid gauge value")
	}
	storage.UpdateGauge(name, Gauge(f))
	return nil
}

// handleCounter конвертирует строку в int64 и добавляет к счётчику.
func handleCounter(storage Storage, name, val string) error {
	i, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return errors.New("invalid counter value")
	}
	storage.UpdateCounter(name, Counter(i))
	return nil
}

func main() {}

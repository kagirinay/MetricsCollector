package main

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

func main() {}

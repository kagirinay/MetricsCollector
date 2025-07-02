package agent

import (
	"math/rand"
	"runtime"
	"sync"
	"time"

	"github.com/kagirinay/MetricsCollector.git/models"
)

const (
	defaultPollInterval   = 2 * time.Second
	defaultReportInterval = 10 * time.Second
)

// StatsCollector хранит текущие значения и ведёт PollCount
type StatsCollector struct {
	rnd        *rand.Rand
	PollCount  models.Counter
	gauges     map[string]models.Gauge
	pollTicker *time.Ticker
	mu         sync.RWMutex
}

func NewStatsCollector() *StatsCollector {
	return &StatsCollector{
		rnd:    rand.New(rand.NewSource(time.Now().UnixNano())),
		gauges: make(map[string]models.Gauge),
	}
}

// StartPolling запускает отдельную горутину,
// обновляющую набор метрик каждые DefaultPollInterval
func (c *StatsCollector) StartPolling(done <-chan struct{}) {
	c.pollTicker = time.NewTicker(defaultPollInterval)
	go func() {
		for {
			select {
			case <-c.pollTicker.C:
				c.poll()
			case <-done:
				c.pollTicker.Stop()
				return
			}
		}
	}()
}

// poll заполняет c.gauges свежими значениями runtime.ReadMemStats
func (c *StatsCollector) poll() {
	c.mu.Lock()
	defer c.mu.Unlock()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	c.gauges["Alloc"] = models.Gauge(m.Alloc)
	c.gauges["BuckHashSys"] = models.Gauge(m.BuckHashSys)
	c.gauges["Frees"] = models.Gauge(m.Frees)
	c.gauges["GCCPUFraction"] = models.Gauge(m.GCCPUFraction)
	c.gauges["GCSys"] = models.Gauge(m.GCSys)
	c.gauges["HeapAlloc"] = models.Gauge(m.HeapAlloc)
	c.gauges["HeapIdle"] = models.Gauge(m.HeapIdle)
	c.gauges["HeapInuse"] = models.Gauge(m.HeapInuse)
	c.gauges["HeapObjects"] = models.Gauge(m.HeapObjects)
	c.gauges["HeapReleased"] = models.Gauge(m.HeapReleased)
	c.gauges["HeapSys"] = models.Gauge(m.HeapSys)
	c.gauges["LastGC"] = models.Gauge(m.LastGC)
	c.gauges["Lookups"] = models.Gauge(m.Lookups)
	c.gauges["MCacheInuse"] = models.Gauge(m.MCacheInuse)
	c.gauges["MCasheSys"] = models.Gauge(m.MCacheSys)
	c.gauges["MSpanInuse"] = models.Gauge(m.MSpanInuse)
	c.gauges["MSpanSys"] = models.Gauge(m.MSpanSys)
	c.gauges["Mallocs"] = models.Gauge(m.Mallocs)
	c.gauges["NextGC"] = models.Gauge(m.NextGC)
	c.gauges["NumForcedGC"] = models.Gauge(m.NumForcedGC)
	c.gauges["NumGC"] = models.Gauge(m.NumGC)
	c.gauges["OtherSys"] = models.Gauge(m.OtherSys)
	c.gauges["PauseTotalNs"] = models.Gauge(m.PauseTotalNs)
	c.gauges["StackInuse"] = models.Gauge(m.StackInuse)
	c.gauges["StackSys"] = models.Gauge(m.StackSys)
	c.gauges["Sys"] = models.Gauge(m.Sys)
	c.gauges["TotalAlloc"] = models.Gauge(m.TotalAlloc)
	c.gauges["RandomValue"] = models.Gauge(c.rnd.Float64())
	c.PollCount++
}

// Export формирует срез "имя/тип/значение",
// готовый к отправке на сервер.
func (c *StatsCollector) Export() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()
	res := make([]string, 0, len(c.gauges)+1)
	for name, v := range c.gauges {
		res = append(res, "/update/gauge"+name+"/"+formatFloat(float64(v)))
	}
	res = append(res, "/update/counter/PollCount"+formatInt(int64(c.PollCount)))
	return res
}

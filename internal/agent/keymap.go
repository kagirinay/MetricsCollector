package agent

import (
	"bytes"
	"log"
	"time"
)

// Agent объединяет сборщик и http-клиент, периодически шлёт данные на сервер.
type Agent struct {
	collector      *StatsCollector
	httpClient     Sender
	serverAddr     string
	reportInterval time.Duration
}

func New(httpClient Sender, serverAddr string) *Agent {
	return &Agent{
		collector:      NewStatsCollector(),
		httpClient:     httpClient,
		serverAddr:     serverAddr,
		reportInterval: defaultReportInterval,
	}
}

// Run блокирует текущую горутину до <-done:
func (a *Agent) Run(done <-chan struct{}) {
	a.collector.StartPolling(done)
	ticker := time.NewTicker(a.reportInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			a.report()
		case <-done:
			return
		}
	}
}

func (a *Agent) report() {
	for _, path := range a.collector.Export() {
		url := a.serverAddr + path
		resp, err := a.httpClient.Post(url, "text/plain", bytes.NewBuffer(nil))
		if err != nil {
			log.Print("Ошибка %s: %v", url, err)
			continue
		}
		_ = resp.Body.Close()
	}
}

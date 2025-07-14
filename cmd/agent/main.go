package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kagirinay/MetricsCollector.git/internal/agent"
)

func main() {
	a := agent.New(&http.Client{}, "http://localhost:8080")
	done := make(chan struct{})
	go a.Run(done)
	// Корректное завершение Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	close(done)
	log.Println("Агент остановлен")
}

package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/kagirinay/MetricsCollector.git/internal/handlers"
	"github.com/kagirinay/MetricsCollector.git/internal/store"
	"github.com/spf13/pflag"
)

func main() {
	// Парсинг флагов
	addr := pflag.StringP("address", "a", "localhost:8080", "HTTP server endpoint address")
	pflag.Parse()
	// Проверка на неизвестные флаги
	if pflag.NArg() > 0 {
		log.Fatalf("Неизвестный флаг или аргумент: %v", pflag.Args())
	}
	st := store.NewMemStorage()
	r := chi.NewRouter()
	// Регистрация обработчиков
	r.Post("/update/{type}/{name}/{value}", handlers.Update(st))
	r.Get("/value/{type}/{name}", handlers.GetMetric(st))
	r.Get("/", handlers.Home(st))
	server := &http.Server{
		Addr:    *addr,
		Handler: r,
	}
	// Безопасное завершение работы
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		log.Printf("Сервер запустился по адресу http://%s", *addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка сервера: %v", err)
		}
	}()
	<-done
	log.Println("Сервер остановлен...")
}

package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/kagirinay/MetricsCollector.git/internal/handlers"
	"github.com/kagirinay/MetricsCollector.git/internal/store"
)

func main() {
	st := store.NewMemStorage()
	mux := http.NewServeMux()
	mux.HandleFunc("/update", handlers.Update(st))
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	log.Println("Сервер запущен по адресу: 8080")
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Ошибка прослушивания: %v", err)
	}
}

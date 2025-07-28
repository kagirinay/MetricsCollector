package handlers

import (
	"embed"
	"html/template"
	"net/http"

	"github.com/kagirinay/MetricsCollector.git/internal/store"
)

//go:embed templates/*
var indexTemplate embed.FS

// Home возвращает HTML-страницу со всеми метриками.
func Home(s store.Storage) http.HandlerFunc {
	// Загружаем шаблон один раз при инициализации
	tmpl := template.Must(template.ParseFS(indexTemplate, "templates/index.html"))

	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем все метрики
		data := struct {
			Gauges   map[string]float64
			Counters map[string]int64
		}{
			Gauges:   make(map[string]float64),
			Counters: make(map[string]int64),
		}
		for name, value := range s.GetAllGauge() {
			data.Gauges[name] = float64(value)
		}
		for name, value := range s.GetAllCounter() {
			data.Counters[name] = int64(value)
		}
		// Рендерим шаблон
		w.Header().Set("Content-Type", "text/html")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Ошибка рендеринга шаблона: "+err.Error(), http.StatusInternalServerError)
		}
	}
}

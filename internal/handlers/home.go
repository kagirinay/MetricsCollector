package handlers

import (
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/kagirinay/MetricsCollector.git/internal/store"
)

// Home возвращает HTML-страницу со всеми метриками.
func Home(s store.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Загружаем шаблон
		tmplPath := filepath.Join("templates", "index.html")
		tmpl, err := template.ParseFiles(tmplPath)
		if err != nil {
			http.Error(w, "Ошибка при загрузке шаблона: "+err.Error(), http.StatusInternalServerError)

			return
		}
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

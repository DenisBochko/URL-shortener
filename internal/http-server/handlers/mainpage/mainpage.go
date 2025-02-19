package mainpage

import (
	"html/template"
	"log/slog"
	"net/http"
)

type PageData struct {
	Title  string
	Result string
}

func New(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// const operation = "handlers.mainpage.New"

		fs := http.FileServer(http.Dir("./static"))
		http.Handle("/static/", http.StripPrefix("/static/", fs))

		// log := log.With(
		// 	slog.String("operation", operation),
		// 	slog.String("request_id", middleware.GetReqID(r.Context())),
		// )
		tmpl := template.Must(template.ParseFiles("internal/frontend/templates/index.html"))
		data := PageData{Title: "Сервис сокращения ссылок"}
		tmpl.Execute(w, data)

	}
}

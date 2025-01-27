package logger

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

func New(log *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// создаём копию логгера, добавляя подсказку, что это компонент "middleware/logger"
		log := log.With(
			slog.String("component", "middleware/logger"),
		)
		// выводится один раз, при запуске
		log.Info("logger middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			// чать, которая будет выполняться при каждом запросе

			// выполняется до обработки запроса
			entry := log.With(
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("remote_addr", r.RemoteAddr),
				slog.String("user_agent", r.UserAgent()),
				slog.String("request_id", middleware.GetReqID(r.Context())), // каждому запросу пресваивается request_id
			)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			// будет вызвана после обрабоки запроса
			defer func() {
				entry.Info("request completed",
					slog.Int("status", ww.Status()),
					slog.Int("bytes", ww.BytesWritten()),
					slog.String("duration", time.Since(t1).String()),
				)
			}()
			// передаём значение следующему хендлеру в цепочке
			next.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}

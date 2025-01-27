package main

import (
	"os"
	"url-shortener/internal/config"
	mylogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

const (
	envLocal = "local"
	envDev   = "dev"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()

	// init logger: slog (import "log/slog") - библиотека для работы с различными логгерами
	log := setupLogger(cfg.Env)

	// Старт приложения, также выводится какое окружение используется
	log.Info("Starting url-shortener", slog.String("env", cfg.Env))
	log.Debug("debug messages are enabled")

	// init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	_ = storage

	// init router: chi (полностью совместим с net/http), render
	router := chi.NewRouter()
	// добавляем middleware
	router.Use(middleware.RequestID) // middleware, которая добавляет id к каждому запросу, чтобы легче было ориентироваться в будущем
	// router.Use(middleware.Logger) // логирование, но он не работает с slog, поэтому в проекте собственный middleware logger
	router.Use(mylogger.New(log)) // используем собственный middleware для логов
	router.Use(middleware.Recoverer) // если случается паника в одном из хендлеров, то приложение восстанавилвается, а не падает 
	router.Use(middleware.URLFormat) // "красивые" url (с id)



	// run server
}

// Функция конфигурации логгера, зависит от входного параметра env, локально мы хотим видеть текстовые логи, на сервере,
// т.е. в окружении dev или prod, мы хотим видеть логи в формате json
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	// если local, то выводим текстовые логи, если dev, то выводим json логи
	switch env {
	case envLocal:
		// создаём логгер, использую текстовые хендлер
		// Level обозначает минимальный уровень логов, который мы выводим
		log = setupPrettySlog()
	case envDev:
		// создаём логгер, использую json хендлер
		// Level обозначает минимальный уровень логов, который мы выводим
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

// функция для запуска красивого логгера
// сообщения выделяются разными цветами, что упрощает чтение логов
func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}

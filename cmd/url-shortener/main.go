package main

import (
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/lib/logger/sl"
	"url-shortener/internal/storage/sqlite"

	"log/slog"
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
		log.Error("failed to init storage", sl.Err(err)) // добавили в лог ошибку
		os.Exit(1)                                       // Падаем с ошибкой
	}

	_ = storage
	// init router: chi (полностью совместим с net/http), render

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
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		// создаём логгер, использую json хендлер
		// Level обозначает минимальный уровень логов, который мы выводим
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return log
}

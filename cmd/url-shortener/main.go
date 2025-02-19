package main

import (
	"net/http"
	"os"
	"url-shortener/internal/config"
	"url-shortener/internal/http-server/handlers/mainpage"
	"url-shortener/internal/http-server/handlers/redirect"
	"url-shortener/internal/http-server/handlers/url/delete"
	"url-shortener/internal/http-server/handlers/url/save"
	mylogger "url-shortener/internal/http-server/middleware/logger"
	"url-shortener/internal/lib/logger/handlers/slogpretty"
	"url-shortener/internal/lib/logger/sl"
	postgresql "url-shortener/internal/storage/postgreSQL"

	// "url-shortener/internal/storage/sqlite"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
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
	log.Info("Starting url-shortener",
		slog.String("env", cfg.Env),
	)
	log.Debug("debug messages are enabled")

	// init storage: sqlite
	// storage, err := sqlite.New(cfg.StoragePath)
	// if err != nil {
	// 	log.Error("failed to init storage", sl.Err(err))
	// 	os.Exit(1)
	// }

	// init storage: postgresql
	storage, err := postgresql.New(cfg.User, cfg.Password, cfg.DBname, cfg.SSLmode, cfg.Port)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// Настройка CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // Разрешить все источники или укажите конкретные
		AllowedMethods: []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders: []string{"Content-Type"},
	})

	// init router: chi (полностью совместим с net/http), render
	router := chi.NewRouter()
	// добавляем middleware
	router.Use(middleware.RequestID) // middleware, которая добавляет id к каждому запросу, чтобы легче было ориентироваться в будущем
	// router.Use(middleware.Logger) // логирование, но он не работает с slog, поэтому в проекте собственный middleware logger
	router.Use(mylogger.New(log))    // используем собственный middleware для логов
	router.Use(middleware.Recoverer) // если случается паника в одном из хендлеров, то приложение восстанавилвается, а не падает
	router.Use(middleware.URLFormat) // "красивые" url (с id)
	router.Use(c.Handler)

	// // группа url для модифицирующих операций
	// router.Route("/url", func(r chi.Router) {
	// 	// максимально простая авторизация, которая предполагает отправку логина и пароля в заголовке
	// 	r.Use(middleware.BasicAuth("url-shortener", map[string]string{
	// 		cfg.HTTPServer.User_auth: cfg.HTTPServer.Password_auth,
	// 	}))

	// 	// хендлеры в группе /url/ с авторизацией
	// 	// r.Post("/", save.New(log, storage))
	// 	// r.Delete("/{alias}", delete.New(log, storage))
	// })

	// главная страница
	router.Get("/", mainpage.New(log))

	// хендлер редиректа
	router.Get("/{alias}", redirect.New(log, storage))
	
	// хендлеры для работы с api
	router.Route("/api", func(r chi.Router) {
		r.Post("/url", save.New(log, storage))
		r.Delete("/url/{alias}", delete.New(log, storage))
	})

	// run server
	log.Info("starting server", slog.String("addres", cfg.Addres))
	srv := &http.Server{
		Addr:         cfg.Addres,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,     // время, чтобы успели прочитать запрос
		WriteTimeout: cfg.HTTPServer.Timeout,     // время, чтобы успели написать ответ на запрос
		IdleTimeout:  cfg.HTTPServer.IdleTimeout, // время жизни соединения клиента с сервером
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}
	// Если добрались до этого блока, то произошла ошибка, т.к. в обычной ситуации до этого блока кода никогда не доберёмся
	log.Error("server stopped")
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

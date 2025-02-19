package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	// struct tag yaml определяет, какое имя у соответствующего параметра будет в yaml файле
	// struct tag env определяет, какое имя будет у переменной окружения, если мы будем использовать её
	// struct tag env-default определяет дефолтное значение
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env:"storage_path" env-required:"true"`
	// Встраиваем структуру HTTPServer в общий конфиг
	HTTPServer `yaml:"http_server"`
	PostgreSQL `yaml:"postgres"`
}

type HTTPServer struct {
	Addres        string        `yaml:"address" env-default:"localhost:8080"`
	Timeout       time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout   time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User_auth     string        `yaml:"user_auth" env-required:"true"`
	Password_auth string        `yaml:"password_auth" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

type PostgreSQL struct {
	User     string `yaml:"user" env-default:"postgres"`
	Port     string `yaml:"port" env-default:"5432"`
	Password string `yaml:"password" env-default:"postgres"`
	DBname   string `yaml:"dbname" env-default:"url-shortener"`
	SSLmode  string `yaml:"sslmode" env-default:"disable"`
}

// Функция загрузки файла конфигурации.
// Функии с приставкой Must не возвращают ошибку, а паникуют, так делать можно в редких случаях, например, при запуске приложения.
func MustLoad() *Config {
	// Загружаем файл local.env
	if err := godotenv.Load("local.env"); err != nil {
		log.Print("No .env file found")
	}
	// Получаем путь до файла конфигураций из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	// Роняем приложение, если не получили путь к конфигу
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// проверяем, существует ли файл по данному пути
	if _, err := os.Stat(configPath); os.IsNotExist(err) { // IsNotExist возвращает true, если файл не существует
		log.Fatalf("config file does not exist: %s", configPath)
	}

	var cfg Config

	// Читаем конфиг-файл и заполняем нашу структуру
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatalf("error reading config file: %s", err)
	}

	return &cfg
}

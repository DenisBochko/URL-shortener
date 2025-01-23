package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3" // init sqlite3 driver
)

// В струкуре лежит коннект для базы
type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	// В этой константе храниться имя этой функции, которое мы вернём в случае ошибки
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}

	// Первый запрос на создание таблички и индекса, если их нет
	// alias - ссылка, по которой будет происходить редирект
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARYKEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)

	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}

	// Делаем запрос
	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}

	return &Storage{db: db}, nil
}

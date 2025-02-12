package postgresql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(user, password, dbname, sslmode string) (*Storage, error) {
	const operation = "storage.postgresql.New"

	connStr := createConnString(user, password, dbname, sslmode)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	// Выполняем создание таблицы и индекса сразу через Exec
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS url(
    	id SERIAL PRIMARY KEY,
    	alias TEXT NOT NULL UNIQUE,
    	url TEXT NOT NULL
	);
	
	CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &Storage{db: db}, nil
}

// Функция для генерации строки коннекта к бд
func createConnString(user, password, dbname, sslmode string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		user,
		password,
		dbname,
		sslmode,
	)
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) {
	const operation = "storage.postgresql.SaveURL"

	var id int64
	err := s.db.QueryRow("INSERT INTO url(url, alias) VALUES($1, $2) RETURNING id", urlToSave, alias).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("%s: failed to insert url: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.postgresql.GetURL"

	var resURL string
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = $1", alias).Scan(&resURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", fmt.Errorf("%s: url not found for alias %s", operation, alias)
		}
		return "", fmt.Errorf("%s: query error: %w", operation, err)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const operation = "storage.postgresql.DeleteURL"

	res, err := s.db.Exec("DELETE FROM url WHERE alias = $1", alias)
	if err != nil {
		return fmt.Errorf("%s: execute delete: %w", operation, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: getting rows affected: %w", operation, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows deleted, alias not found", operation)
	}

	return nil
}

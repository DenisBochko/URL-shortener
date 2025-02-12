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
		return nil, fmt.Errorf("%s: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}
    defer db.Close()

	// Первый запрос на создание таблички и индекса, если их нет
	// alias - ссылка, по которой будет происходить редирект
	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
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

// Функция для генерации строки коннекта к бд
func createConnString(user, password, dbname, sslmode string) string {
	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
			user,
			password,
			dbname,
			sslmode,
		)
}

func (s *Storage) SaveURL(urlToSave string, alias string) (int64, error) { // возвращаем индекс добавленного элемента и ошибку
	// В этой константе храниться имя этой функции, которое мы вернём в случае ошибки
	const operation = "storage.postgresql.SaveURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	defer stmt.Close()

	// отправляем запрос
	res, err := stmt.Exec(urlToSave, alias)
	
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}

	// Получаем id добавленной записи
	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", operation, err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.postgresql.GetURL"

	stmt, err := s.db.Prepare("SELECT url FROM url WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: prepare statement: %w", operation, err)
	}
	defer stmt.Close()
	var resURL string

	err = stmt.QueryRow(alias).Scan(&resURL) // выполняет подготовленный запрос
	if err != nil {
		return "", fmt.Errorf("%s: failed to get last insert id: %w", operation, err) // возвращаем имя функции, как ошибку в ней (оборачиваем ошибку)
	}

	return resURL, nil
}

func (s *Storage) DeleteURL(alias string) error {
	const operation = "storage.postgresql.DeleteURL"

	stmt, err := s.db.Prepare("DELETE FROM url WHERE alias = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", operation, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(alias)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", operation, err)
	}

	rowsAffected, err := res.RowsAffected() // возвращает количество строк, затронутых обновлением
	if err != nil {
		return fmt.Errorf("%s: getting rows affected: %w", operation, err)
	}
	// проверяем, действительно ли произошло удаление
	if rowsAffected == 0 {
		return fmt.Errorf("%s: no rows deleted, alias not found", operation)
	}

	return nil
}
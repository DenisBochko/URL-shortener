package storage

import "errors"

/*

Общая информация для реализации различных вариантов storage
Но интерфейс самого storage будет находиться по месту использования

*/

// Определям общие ошибки для storage
var (
	ErrURLNotFound = errors.New("url not found")
	ErrURLExists = errors.New("url exists")
)
package main

import (
	"fmt"
	"url-shortener/internal/config"
)

func main() {
	// init config: cleanenv
	cfg := config.MustLoad()

	fmt.Print(cfg)

	// init logger: slog (import "log/slog")

	// init storage: sqlite

	// init router: chi (полностью совместим с net/http), render

	// run server
}
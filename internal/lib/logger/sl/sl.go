package sl

import "log/slog"

// Для удобной обёртки ошибок 
func Err(err error) slog.Attr {
	return slog.Attr{
		Key: "error",
		Value: slog.StringValue(err.Error()),
	}
}
package log

import "log/slog"

type Config struct {
	Level slog.Level `default:"info"`
	JSON  bool       `default:"true"`
}

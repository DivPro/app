package log

import (
	"log/slog"
	"os"
)

func SetLogger(conf Config) {
	var logHandler slog.Handler
	if conf.JSON {
		logHandler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: conf.Level,
		})
	} else {
		logHandler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: conf.Level,
		})
	}
	slog.SetDefault(slog.New(logHandler))
}

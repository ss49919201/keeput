package appslog

import (
	"log/slog"
	"os"
	"strings"

	"github.com/samber/lo"
	"github.com/ss49919201/keeput/app/analyzer/internal/config"
)

func Init() {
	slog.SetDefault(
		slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: lo.Switch[string, slog.Level](strings.ToUpper(config.LogLevel())).
					Case("DEBUG", slog.LevelDebug).
					Case("INFO", slog.LevelInfo).
					Case("WARN", slog.LevelWarn).
					Case("ERROR", slog.LevelError).
					Default(slog.LevelWarn),
			}),
		),
	)
}

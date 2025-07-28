package init

import (
	"github.com/mayberries/shared/pkg/config"
	"github.com/mayberries/shared/pkg/log"
	"log/slog"
)

func Log(options *config.LoggerOptions) {
	handlerOptions := map[string]*slog.HandlerOptions{
		"debug":    {Level: slog.LevelDebug},
		"infoUser": {Level: slog.LevelInfo},
		"warn":     {Level: slog.LevelWarn},
		"error":    {Level: slog.LevelError},
	}

	switch options.Format {
	case "json":
		log.Configure(nil, handlerOptions[options.Level], true)
	case "text":
		log.Configure(nil, handlerOptions[options.Level], false)
	}
}

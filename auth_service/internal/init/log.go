package init

import (
	"auth_service/pkg/config"
	"auth_service/pkg/log"
	"log/slog"
)

func Log(options *config.LoggerOptions) {
	handlerOptions := map[string]*slog.HandlerOptions{
		"debug": {Level: slog.LevelDebug},
		"info":  {Level: slog.LevelInfo},
		"warn":  {Level: slog.LevelWarn},
		"error": {Level: slog.LevelError},
	}

	switch options.Format {
	case "json":
		log.Configure(nil, handlerOptions[options.Level], true)
	case "text":
		log.Configure(nil, handlerOptions[options.Level], false)
	}
}

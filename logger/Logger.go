package logger

import (
	"encoding/json"
	"fmt"
	"log"

	"go.uber.org/zap"
)

func Init(rawJSON []byte) (*zap.Logger, error) {
	var cfg zap.Config
	if err := json.Unmarshal(rawJSON, &cfg); err != nil {
		log.Println(err)
	}

	logger := zap.Must(cfg.Build())
	defer func(logger *zap.Logger) {
		err := logger.Sync()
		if err != nil {
			fmt.Println("Failed to initialize logger")
		}
	}(logger)

	return logger, nil
}

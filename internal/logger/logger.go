package logger

import (
	"fmt"
	"go.uber.org/zap"
)

var Logger *zap.Logger

func Init() error {
	var err error
	Logger, err = zap.NewProduction()
	if err != nil {
		return fmt.Errorf("error initializing logger: %w", err)
	}
	return nil
}

package framework

import (
	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() error {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		return err
	}
	defer Logger.Sync()
	return nil
}

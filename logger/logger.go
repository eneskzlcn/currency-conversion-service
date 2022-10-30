package logger

import (
	"go.uber.org/zap"
)

func NewZapLoggerForEnv(env string, callerSkip int) (*zap.SugaredLogger, error) {
	if env == "prod" {
		logger, err := zap.NewProduction(zap.AddCallerSkip(callerSkip), zap.AddStacktrace(zap.ErrorLevel))
		return logger.Sugar(), err
	}
	logger, err := zap.NewDevelopment(zap.AddCallerSkip(callerSkip))
	return logger.Sugar(), err
}

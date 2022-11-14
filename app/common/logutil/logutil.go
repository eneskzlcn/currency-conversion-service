package logutil

import "go.uber.org/zap"

func LogThenReturn(logger *zap.SugaredLogger, err error) error {
	logger.Error(err)
	return err
}

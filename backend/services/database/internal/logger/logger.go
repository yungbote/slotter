package logger

import (
  "go.uber.org/zap"
  "sync"
)

var (
  LogOnce sync.Once
  l         *zap.Logger
)

func GetLogger() *zap.Logger {
  logOnce.Do(func() {
    logger, err := zap.NewProduction()
    if err != nil {
      panic("failed to initialize zap logger: " + err.Error())
    }
    l = logger
  })
  return l
}

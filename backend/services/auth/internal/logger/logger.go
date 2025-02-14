package logger

import (
  "go.uber.org/zap"
  "sync"
)

var (
  once sync.Once
  l       *zap.Logger
)

func GetLogger() *zap.Logger {
  once.Do(func() {
    logger, err := zap.NewProduction()
    if err != nil {
      panic("Failed to init zap: " + err.Error())
    }
    l = logger
  })
  return l
}

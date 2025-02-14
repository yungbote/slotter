package main

import (
  "os"
  "github.com/yungbote/slotter/backend/services/auth/internal/api"
  "github.com/yungbote/slotter/backend/services/auth/internal/logger"
  "go.uber.org/zap"
)

func main() {
  log := logger.GetLogger()
  port := os.Getenv("AUTH_PORT")
  if port == "" {
    port = "8090"
    log.Warn("AUTH_PORT not set, using default", zap.String("port", port))
  }
  r := api.NewRouter()
  log.Info("Starting Auth service", zap.String("port", port))
  if err := r.Run(":" + port); err != nil {
    log.Fatal("Failed to start auth server", zap.Error(err))
  }
}

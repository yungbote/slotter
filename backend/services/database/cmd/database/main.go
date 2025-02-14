package main

import (
  "go.uber.org/zap"
  "os"
  "github.com/yungbote/slotter/backend/services/database/internal/api"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
  "github.com/yungbote/slotter/backend/services/database/internal/logger"
)

func main() {
  log := logger.GetLogger()
  dsn := os.Getenv("DATABASE_DSN")
  if dsn == "" {
    dsn = "postgres://bote:bote@db:5432/slotter?sslmode=disable"
    log.Warn("DATABASE_DSN not set, using default", zap.String("dsn", dsn))
  }
  migrationsPath := os.Getenv("MIGRATIONS_PATH")
  if migrationsPath == "" {
    migrationsPath = "migrations"
    log.Warn("MIGRATIONS_PATH not set, using default", zap.String("migrationsPath", migrationsPath))
  }
  db, err := database.InitDB(dsn)
  if err != nil {
    log.Fatal("Failed to init DB", zap.Error(err))
  }
  if err := database.MigrateDB(db, migrationsPath); err != nil {
    log.Fatal("Migrations failed", zap.Error(err))
  }
  router := api.NewRouter(db)
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Warn("PORT not set, using default", zap.String("port", port))
  }
  log.Info("Starting database service", zap.String("port", port))
  if err := router.Run(":" + port); err != nil {
    log.Fatal("Gin server shutdown", zap.Error(err))
  }
}

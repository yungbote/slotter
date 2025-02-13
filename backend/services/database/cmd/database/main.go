package main

import (
  "log"
  "os"
  "github.com/gin-gonic/gin"
  "github.com/yungbote/slotter/backend/services/database/internal/api"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
)

func main() {
  dsn := os.Getenv("DATABASE_DSN")
  if dsn == "" {
    dsn = "postgres://user:pass@localhost:5432/mydb?sslmode=disable"
    log.Printf("WARN: DATABASE_DSN not set, using default: %s", dsn)
  }
  migrationsPath := os.Getenv("MIGRATIONS_PATH")
  if migrationsPath == "" {
    migrationsPath = "migrations"
    log.Printf("WARN: MIGRATIONS_PATH not set, using default: %s", migrationsPath)
  }
  if err := database.InitDB(dsn); err != nil {
    log.Fatalf("ERROR: Failed to init DB: %v", err)
  }
  if err := database.MigrateDB(dsn, migrationsPath); err != nil {
    log.Fatalf("ERROR: Migrations failed: %v", err)
  }
  router := api.NewRouter()
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
  }
  log.Printf("INFO: Starting database service on port: %s", port)
  if err := router.Run(":" + port); err != nil {
    log.Fatalf("ERROR: Gin server shutdown: %v", err)
  }
}

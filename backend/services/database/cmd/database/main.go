package main

import (
  "log"
  "os"
  "github.com/yungbote/slotter/backend/services/database/internal/api"
  "github.com/yungbote/slotter/backend/services/database/internal/database"
)

func main() {
  dsn := os.Getenv("DATABASE_DSN")
  if dsn == "" {
    dsn = "postgres://bote:bote@db:5432/slotter?sslmode=disable"
    log.Printf("WARN: DATABASE_DSN not set, using default: %s", dsn)
  }
  migrationsPath := os.Getenv("MIGRATIONS_PATH")
  if migrationsPath == "" {
    migrationsPath = "migrations"
    log.Printf("WARN: MIGRATIONS_PATH not set, using default: %s", migrationsPath)
  }
  db, err := database.InitDB(dsn)
  if err != nil {
    log.Fatalf("ERROR: Failed to init DB: %v", err)
  }
  if err := database.MigrateDB(db, migrationsPath); err != nil {
    log.Fatalf("ERROR: Migrations failed: %v", err)
  }
  router := api.NewRouter(db)
  port := os.Getenv("PORT")
  if port == "" {
    port = "8080"
    log.Printf("WARN: PORT not set, using default: %s", port)
  }
  log.Printf("INFO: Starting database service on port: %s", port)
  if err := router.Run(":" + port); err != nil {
    log.Fatalf("ERROR: Gin server shutdown: %v", err)
  }
}

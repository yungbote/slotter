package database

import (
  "fmt"
  "go.uber.org/zap"
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/logger"
) func MigrateDB(db *gorm.DB, migrationsPath string) error {
  sqlDB, err := db.DB()
  if err != nil {
    return fmt.Errorf("ERROR: Failed to get *sql.DB: %w", err)
  }
  driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
  if err != nil {
    return fmt.Errorf("ERROR: Failed to create postgres driver: %w", err)
  }
  m, err := migrate.NewWithDatabaseInstance(
    "file://"+migrationsPath,
    "postgres",
    driver,
  )
  err = m.Up()
  if err != nil && err != migrate.ErrNoChange {
    return fmt.Errorf("ERROR: Migration up error: %w", err)
  }
  logger.GetLogger().Info("Migrations applied", zap.String("migrationsPath", migrationsPath))
  return nil
}

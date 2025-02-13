package database

import (
  "fmt"
  "github.com/golang-migrate/migrate/v4"
  "github.com/golang-migrate/migrate/v4/database/postgres"
  _ "github.com/golang-migrate/migrate/v4/source/file"
)

func MigrateDB(dsn, migrationsPath string) error {
  sqlDB, err := DB.DB()
  if  err != nil {
    return fmt.Errorf("ERROR: Failed to get *sql.DB from GORM: %w", err)
  }
  driver, err := postgres.WithInstance(sqlDB, &postgres.Config{})
  if err != nil {
    return fmt.Errorf("ERROR: Failed to create a postgres driver: %w", err)
  }
  m, err := migrate.NewWithDatabaseInstance(
    "file://"+migrationsPath,
    "postgres",
    driver,
  )
  if err != nil {
    return fmt.Errorf("ERROR: Migration init error: %w", err)
  }
  err = m.Up()
  if err != nil && err != migrate.ErrNoChange {
    return fmt.Errorf("ERROR: Migration up error: %w", err)
  }
  return nil
}

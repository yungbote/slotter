package database

import (
  "fmt"
  "go.uber.org/zap"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
  "github.com/yungbote/slotter/backend/services/database/internal/logger"
)

func InitDB(dsn string) (*gorm.DB, error) {
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{SingularTable: true},
  })
  if err != nil {
    return nil, fmt.Errorf("ERROR: Failed to connect to database: %w", err)
  }
  logger.GetLogger().Info("Connected to database", zap.String("dsn", dsn))
  return db, nil
}

package database

import (
  "fmt"
  "log"
  "gorm.io/driver/postgres"
  "gorm.io/gorm"
  "gorm.io/gorm/schema"
)

var DB *gorm.DB

func InitDB(dsn string) error {
  var err error
  DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
    NamingStrategy: schema.NamingStrategy{
      SingularTable: true,
    },
  })
  if err != nil {
    return fmt.Errorf("ERROR: Failed to connect database: %w", err)
  }
  log.Println("INFO: GORM connected to DB.")
  return nil
}

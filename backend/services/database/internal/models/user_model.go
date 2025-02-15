package models

import (
  "time"
)

type User struct {
  ID            uint      `gorm:"primaryKey"`
  Email         string    `gorm:"uniqueIndex;not null"`
  PasswordHash  string    `gorm:"not null"`
  FirstName     string    `gorm:"not null"`
  LastName      string    `gorm:"not null"`
  CompanyID     uint
  Company       Company
  CreatedAt     time.Time
}

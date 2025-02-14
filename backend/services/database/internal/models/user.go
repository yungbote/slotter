package models

import (
  "time"
)

type User struct {
  ID            uint      `gorm:"primaryKey"`
/*  RoleID        uint
  Role          Role*/
  Email         string    `gorm:"uniqueIndex;not null"`
  PasswordHash  string    `gorm:"not null"`
  FullName      string    `gorm:"not null"`
  Status        string    `gorm:"not null;default:'pending'"`
  CompanyID     uint
  Company       Company
  CreatedAt     time.Time
}

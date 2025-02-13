package models

type Role struct {
  ID          uint      `gorm:"primaryKey"`
  Name        string    `gorm:"not null"`
  Permissions string    `gorm:"type:jsonb;default:'[]'"`
  CompanyID   uint
  Company     Company
}

package models

type Company struct {
  ID      uint    `gorm:"primaryKey"`
  Name    string  `gorm:"uniqueIndex;not null"`
}


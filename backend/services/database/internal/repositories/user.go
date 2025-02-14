package repositories

import (
  "gorm.io/gorm"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type UserRepository interface {
  Create(user *models.User) error
  GetByID(id uint) (*models.User, error)
  GetByEmail(email string) (*models.User, error)
  Update(user *models.User) error
  Delete(user *models.User) error
}

type userRepository struct {
  db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
  return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
  return r.db.Create(user).Error
} 

func (r *userRepository) GetById(id uint) (*models.User, error) {
  var u models.User
  err := r.db.First(&u, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ErrNotFound
  }
  if err != nil {
    return nil, err
  }
  return &u, nil
}

func (r *userRepository) Update(user *models.User) error {
  return r.db.Save(user).Error
}

func(r *userRepository) Delete(user *models.User) error {
  return r.db.Delete(user).Error
}

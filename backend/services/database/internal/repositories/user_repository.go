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
  err := r.db.Create(user).Error
  if err != nil {
    return ParseDBError("UserRepository.Create", err)
  }
  return nil
}

func (r *userRepository) GetByID(id uint) (*models.User, error) {
  var u models.User
  err := r.db.First(&u, id).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {

    return nil, ParseDBError("UserRepository.GetByID", err)
  }
  if err != nil {
    return nil, ParseDBError("UserRepository.GetByID", err)
  }
  return &u, nil
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
  var u models.User
  err := r.db.Where("email = ?", email).First(&u).Error
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return nil, ParseDBError("UserRepository.GetByEmail", err)
  }
  if err != nil {
    return nil, ParseDBError("UserRepository.GetByEmail", err)
  }
  return &u, nil
}

func (r *userRepository) Update(user *models.User) error {
  err := r.db.Save(user).Error
  if err != nil {
    return ParseDBError("UserRepository.Update", err)
  }
  return nil
}

func (r *userRepository) Delete(user *models.User) error {
  err := r.db.Delete(user).Error
  if err != nil {
    return ParseDBError("UserRepository.Delete", err)
  }
  return nil
}

func (r *userRepository) ListByCompanyID(companyID uint) ([]*models.User, error) {
  const op = "UserRepository.ListByCompanyID"
  var users []*models.User
  err := r.db.Where("company_id = ?", companyID).Find(&users).Error
  if err != nil {
    return nil, ParseDBError(op, err)
  }
  return users, nil
}

func (r *userRepository) CountByCompanyID(companyID uint) (int64, error) {
  const op = "UserRepository.CountByCompanyID"
  var count int64
  err := r.db.Model(&models.User{}).Where("company_id = ?", companyID).Count(&count).Error
  if err != nil {
    return 0, ParseDBError(op, err)
  }
  return count, nil
}



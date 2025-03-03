package repos

import (
  "fmt"

  "gorm.io/gorm"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type UserFilter struct {
  CompanyID   uuid.UUID
  Role        string
  EmailLike   string
  SortField   string
  SortDir     string
}

type URepo interface {
  //GENERAL CRUD
  Create(user models.User) (*models.User, error)
  UpdateFirstName(userID uuid.UUID, newFirst string) error
  UpdateLastName(userID uuid.UUID, newLast string) error
  UpdateEmail(userID uuid.UUID, newEmail string) error
  UpdatePassword(userID uuid.UUID, newPass string) error
  UpdateRole(userID uuid.UUID, newRole string) error
  UpdateAvatarURL(userID uuid.UUID, newAvatarURL string) error
  GetByID(userID uuid.UUID) (*models.User, error)
  GetByEmail(userEmail string) (*models.User, error)
  Delete(userID uuid.UUID) error
  ListUsers(f UserFilter) ([]*models.User, error)

}
type uRepo struct {
  db *gorm.DB
}

func NewURepo(db *gorm.DB) URepo {
  return &uRepo{db: db}
}

func (r *uRepo) Create(user models.User) (*models.User, error) {
  if err := r.db.Create(&user).Error; err != nil {
    return nil, fmt.Errorf("Failed to create user: %w", err)
  }
  return &user, nil
}

func (r *uRepo) UpdateFirstName(userID uuid.UUID, newFirst string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("first_name", newFirst).Error
}

func (r *uRepo) UpdateLastName(userID uuid.UUID, newLast string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("last_name", newLast).Error
}

func (r *uRepo) UpdateEmail(userID uuid.UUID, newEmail string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("email", newEmail).Error
}

func (r *uRepo) UpdatePassword(userID uuid.UUID, newPass string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("password", newPass).Error
}

func (r *uRepo) UpdateRole(userID uuid.UUID, newRole string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("role", newRole).Error
}

func (r *uRepo) UpdateAvatarURL(userID uuid.UUID, newAvatarURL string) error {
  return r.db.Model(&models.User{}).
    Where("id = ?", userID).
    Update("avatar_url", newAvatarURL).Error
}

func (r *uRepo) GetByID(userID uuid.UUID) (*models.User, error) {
  var u models.User
  if err := r.db.First(&u, "id = ?", userID).Error; err != nil {
    return nil, fmt.Errorf("user not found: %w", err)
  }
  return &u, nil
}

func (r *uRepo) GetByEmail(userEmail string) (*models.User, error) {
  var u models.User
  if err := r.db.First(&u, "email = ?", userEmail).Error; err != nil {
    return nil, fmt.Errorf("user not found: %w", err)
  }
  return &u, nil
}

func (r *uRepo) Delete(userID uuid.UUID) error {
  u, err := r.GetByID(userID)
  if err != nil {
    return err
  }
  if err := r.db.Delete(u).Error; err != nil {
    return fmt.Errorf("Failed to delete user: %w", err)
  }
  return nil
}

func (r *uRepo) ListUsers(f UserFilter) ([]*models.User, error) {
  dbq := r.db.Model(&models.User{})
  if f.CompanyID != uuid.Nil {
    dbq = dbq.Where("company_id = ?", f.CompanyID)
  }
  if f.Role != "" {
    dbq = dbq.Where("role = ?", f.Role)
  }
  if f.EmailLike != "" {
    dbq = dbq.Where("email ILIKE ?", "%"+f.EmailLike+"%")
  }
  allowed := []string{"email", "first_name", "created_at"}
  dbq = applySorting(dbq, f.SortField, f.SortDir, allowed)
  var users []*models.User
  if err := dbq.Find(&users).Error; err != nil {
    return nil, err
  }
  return users, nil
}



package services

import (
  "fmt"
  "strings"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type UserService interface {
  CreateUser(input *models.User) (*models.User, error)
  GetUserByID(id uint) (*models.User, error)
  GetUserByEmail(email string) (*models.User, error)
  UpdateUser(existing *models.User, updates *models.User) (*models.User, error)
  DeleteUser(existing *models.User) error
}

type userService struct {
  userRepo            repositories.UserRepository
  companyRepo         repositories.CompanyRepository
}

func NewUserService(repo repositories.UserRepository, companyRepo repositories.CompanyRepository) UserService {
  return &userService{repo: repo, companyRepo: companyRepo}
}

func (s *userService) CreateUser(input *models.User) (*models.User, error) {
  if err := s.validateCreateUserInput(input); err != nil {
    return nil, err
  }
  if err := s.userRepo.Create(input); err != nil {
    return nil, err
  }
  return input, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
  const op = "UserService.GetUserByID"
  if id == 0 {
    return nil, repositories.NewDomainError(op, repositories.ErrCodeValidation, "user id must be > 0", nil)
  }
  user, err := s.userRepo.GetbyID(id)
  if err != nil {
    return nil, err
  }
  return user, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
  const op = "UserService.GetUserByEmail"
  email = strings.TrimSpace(email)
  if email = "" {
    return nil, repositories.NewDomainError(op, repositories.ErrCodeValidation, "email cannot be empty", nil)
  }
  user, err := s.userRepo.GetByEmail(email)
  if err != nil {
    return nil, err
  }
  return user, nil
}

func (s *userService) UpdateUser(existing *models.User, updates *models.User) (*models.User, error) {
  existing.Email = updates.Email
  existing.PasswordHash = updates.PasswordHash
  existing.FullName = updates.FullName
  existing.CompanyID = updates.CompanyID
  if err := s.userRepo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *userService) DeleteUser(existing *models.User) error {
  return s.userRepo.Delete(existing)
}

func (s *userService) validateCreateUserInput(u *models.User) error {
  const op = "UserService.validateCreateUserInput"
  email := strings.TrimSpace(u.Email)
  if email == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "email is required", nil)
  }
  if !strings.Contains(email, "@") {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "invalid email format (missing '@')", nil)
  }
  if strings.TrimSpace(u.PasswordHash) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "password is required", nil)
  }
  if strings.TrimSpace(u.FullName) == "" {
    return repositories.NewDomainError(op, repositories.ErrCodeValidation, "full name is required", nil)
  }
  existingUser, err := s.userRepo.GetByEmail(email)
  if err == nil && existingUser != nil {
    return repositories.NewDomainError(op, repositories.ErrCodeDuplicate, fmt.Sprintf("the email '%s' is already in use", email), nil)
  } else if dErr, ok := err.(*repositories.DomainError); ok {
    if dErr.Code != repositories.ErrCodeNotFound {
      return dErr
    }
  } else if err != nil {
    return repositories.NewDomainError(op, repositories.ErrCodeUnknown, err.Error(), err)
  }
  if u.CompanyID != 0 {
    _, err := s.companyRepo.GetByID(u.CompanyID)
    if err != nil {
      if dErr, ok := err.(*repositories.DomainError); ok {
        if dErr.Code == repositories.ErrCodeNotFound {
          return repositories.NewDomainError(op, repositories.ErrCodeForeignKey, fmt.Sprintf("company with id %d does not exist", u.CompanyID), dErr)
        }
        return dErr
      }
      return repositories.NewDomainError(op, repositories.ErrCodeUnknown, err.Error(), err)
    }
  }
  return nil
} 

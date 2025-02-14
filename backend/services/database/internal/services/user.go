package services

import (
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type UserService interface {
  CreateUser(input *models.User) (*models.User, error)
  GetUserByID(id uint) (*models.User, error)
  UpdateUser(existing *models.User, updates *models.User) (*models.User, error)
  DeleteUser(existing *models.User) error
}

type userService struct {
  repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
  return &userService{repo: repo}
}

func (s *userService) CreateUser(input *models.User) (*models.User, error) {
  if err := s.repo.Create(input); err != nil {
    return nil, err
  }
  return input, nil
}

func (s *userService) GetUserByID(id uint) (*models.User, error) {
  return s.repo.GetByID(id)
}

func (s *userService) UpdateUser(existing *models.User, updates *models.User) (*models.User, error) {
  existing.Email = updates.Email
  existing.PasswordHash = updates.PasswordHash
  existing.FullName = updates.FullName
  existing.CompanyID = updates.CompanyID
  if err := s.repo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *userService) DeleteUser(existing *models.User) error {
  return s.repo.Delete(existing)
}

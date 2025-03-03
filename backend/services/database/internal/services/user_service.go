package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
)

type USvc interface {
 //GENERAL CRUD
  CreateUser(user models.User) (*models.User, error)
  UpdateUserFirstName(userID uuid.UUID, newFirst string) error
  UpdateUserLastName(userID uuid.UUID, newLast string) error
  UpdateUserEmail(userID uuid.UUID, newEmail string) error
  UpdateUserPassword(userID uuid.UUID, newPass string) error
  UpdateUserRole(userID uuid.UUID, newRole string) error
  UpdateUserAvatarURL(userID uuid.UUID, newAvatarURL string) error
  GetUserByID(userID uuid.UUID) (*models.User, error)
  GetUserByEmail(email string) (*models.User, error)
  DeleteUser(userID uuid.UUID) error
  
  ListUsers(f repos.UserFilter) ([]*models.User, error)

}

type uSvc struct {
  repo        repos.UserRepo
  useractionRepo  repos.UserActionRepo
  pub             events.PubSubPublisher
}

func NewUSvc(
  repo        repos.URepo,
) USvc {
  return &uSvc{repo: repos.URepo}
}

func (s *uSvc) CreateUser(user models.User) (*models.User, error) {
  if user.Email == "" {
    return nil, fmt.Errorf("email is required")
  }
  if user.CompanyID == uuid.Nil {
    return nil, fmt.Errorf("company is required")
  }
  if user.FirstName == "" {
    return nil, fmt.Errorf("First Name is required")
  }
  if user.LastName == "" {
    return nil, fmt.Errorf("Last Name is required")
  }
  if user.AvatarURL == "" {
    return nil, fmt.Errorf("Avatar URL is required")
  }
  if user.Password == "" {
    return nil, fmt.Errorf("Password is required")
  }
  created, err := s.repo.Create(user)
  if err != nil {
    return nil, fmt.Errorf("Failed to create user: %w", err)
  }
  return created, nil
}

func (s *uSvc) UpdateUserFirstName(userID uuid.UUID, newFirst string) error {
  if userID == uuid.Nil || newFirst == "" {
    return fmt.Errorf("Invalid input for update first name")
  }
  return s.repo.UpdateFirstName(userID, newFirst)
}

func (s *uSvc) UpdateUserLastName(userID uuid.UUID, newLast string) error {
  if userID == uuid.Nil || newLast == "" {
    return fmt.Errorf("Invalid input for update last name")
  }
  return s.repo.UpdateLastName(userID, newLast)
}

func (s *uSvc) UpdateUserEmail(userID uuid.UUID, newEmail string) error {
  if userID == uuid.Nil || newEmail == "" {
    return fmt.Errorf("Invalid input for update email")
  }
  return s.repo.UpdateEmail(userID, newEmail)
}

func (s *uSvc) UpdateUserPassword(userID uuid.UUID, newPass string) error {
  if userID == uuid.Nil || newPass == "" {
    return fmt.Errorf("Invalid input for update password")
  }
  return s.repo.UpdatePassword(userID, newPass)
}

func (s *uSvc) UpdateUserRole(userID uuid.UUID, newRole string) error {
  if userID == uuid.Nil || newRole == "" {
    return fmt.Errorf("Invalid input for update role")
  }
  return s.repo.UpdateRole(userID, newRole)
}

func (s *uSvc) GetUserByID(userID uuid.UUID) (*models.User, error) {
  if userID == uuid.Nil {
    return nil, fmt.Errorf("Invalid UserID")
  }
  user, err := s.repo.GetByID(userID)
  if err != nil {
    return nil, fmt.Errorf("failed to get user by id: %w", err)
  }
  return user, nil
}

func (s *uSvc) GetUserByEmail(email string) (*models.User, error) {
  if email == "" {
    return nil, fmt.Errorf("email is empty")
  }
  return s.repo.GetByEmail(email)
}

func (s *uSvc) GetUserByEmail(userEmail string) (*models.User, error) {
  if userEmail == "" {
    return nil, fmt.Errorf("email is empty")
  }
  user, err := s.repo.GetByEmail(userEmail)
  if err != nil {
    return nil, fmt.Errorf("failed to get user by email: %w", err)
  }
  return user, nil
}

func (s *uSvc) ListUsers(f repos.UserFilter) ([]*models.User, error) {
  users, err := s.repo.ListUsers(f)
  if err != nil {
    return nil, fmt.Errorf("failed to list users: %w", err)
  }
  return users, nil
}


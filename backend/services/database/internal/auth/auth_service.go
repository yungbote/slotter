package auth

import (
  "context"
  "fmt"
  "strings"
  "errors"
  "github.com/google/uuid"
  "golang.org/x/crypto/bcrypt"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/services"
)

type AuthService interface {
  RegisterLocalUser(email, password, firstName, lastName string) (*models.User, string, error)
  LoginLocalUser(email, password string) (*models.User, string, error)
  LoginWithGoogle(code, state string) (*models.User, string, error)
}

type authService struct {
  userSvc       services.UserSvc
  tokenService  TokenService
  oauthService  OAuthService
}

func NewAuthService(userSvc services.UserSvc, tokenService TokenService, oauthService OAuthService) AuthService {
  return &authService{userSvc: userSvc, tokenService: tokenService, oauthService: oauthService}
}

func (a *authService) RegisterLocalUser(email, password, firstName, lastName string) (*models.User, string, error) {
  email = strings.TrimpSpace(email)
  if email == "" || password == "" || firstName == "" || lastName == "" {
    return nil, "", fmt.Errorf("missing required fields")
  }
  hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
  if err != nil {
    return nil, "", fmt.Errorf("failed to hash password: %w", err)
  }
  newUser := &models.User{
    Email:      email,
    Password:   string(hashed),
    FirstName:  firstName,
    LastName:   lastName,
  }
  createdUser, err := a.appSvc.CreateUser(newUser)
  if err != nil {
    return nil, "", err
  }
  token, err := a.tokenService.GenerateToken(createdUser.ID, createdUser.Email)
  if err != nil {
    return nil, "", fmt.Errorf("failed to generate token: %w", err)
  }
  return createdUser, token, nil
}

func (a *authService) LoginLocalUser(email, password string) (*models.User, string, error) {
  email = strings.TrimSpace(email)
  if email == "" || password == "" {
    return nil, "", errors.New("missing email or password")
  }
  user, err := a.userSvc.GetUserByEmail(email)
  if err != nil {
    return nil, "", errors.New("invalid credentials")
  }
  if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
    return nil, "", errors.New("invalid credentials")
  }
  token, err := a.tokenService.GenerateToken(user.ID, user.Email)
  if err != nil {
    return nil, "", err
  }
  return user, token, nil
}

func (a *authService) LoginWithGoogle(code, state string) (*models.User, string, error) {
  user, err := a.oauthService.HandleGoogleCallback(code, state)
  if err != nil {
    return nil, "", err
  }
  token, err := a.tokenService.GenerateToken(user.ID, user.Email)
  if err != nil {
    return nil, "", err
  }
  return user, token, nil
}

package auth

import (
  "context"
  "encoding/json"
  "fmt"
  "net/http"
  "golang.org/x/oauth2"
  "golang.org/x/oauth2/google"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/services"
  "github.com/google/uuid"
)

type OAuthConfig struct {
  TwitterClientID     string
  TwitterSecret       string
  AppleClientID       string
  AppleSecret         string
  GoogleClientID      string
  GoogleSecret        string
  RedirectURL         string
}

type OAuthService interface {
  GetGoogleLoginURL(state string) string
  HandleGoogleCallback(code, state string) (*models.User, error)
}

type oauthService struct {
  googleConfig *oauth2.Config
  userSvc       services.UserSvc
}

func NewOAuthService(cfg OAuthConfig, userSvc services.UserSvc) (OAuthService, error) {
  googleConf := &oauth2.Config{
    ClientID:       cfg.GoogleClientID,
    ClientSecret:   cfg.GoogleSecret,
    Endpoint:       google.Endpoint,
    Scopes:         []string{"openid", "profile", "email"},
    RedirectURL:    cfg.RedirectURL,
  }
  return &oauthService{
    googleConfig: googleConf,
    userSvc:      userSvc,
  }, nil
}

func (o *oauthService) GetGoogleLoginURL(state string) string {
  return o.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOnLine)
}

func (o *oauthService) HandleGoogleCallback(code, state string) (*models.User, error) {
  token, err := o.googleConfig.Exchange(context.Background(), code)
  if err != nil {
    return nil, fmt.Errorf("failed to exchange code: %w", err)
  }

  client := o.googleConfig.Client(context.Background(), token)
  resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
  if err != nil {
    return nil, fmt.Errorf("failed to fetch google userinfo: %w", err)
  }
  defer resp.Body.Close()

  if resp.StatusCode >= 300 {
    return nil, fmt.Errorf("userinfo returned status %d", resp.StatusCode)
  }

  var info struct {
    ID              string      `json:"id"`
    Email           string      `json:"email"`
    VerifiedEmail   bool        `json:"verified_email"`
    Name            string      `json:"name"`
    GivenName       string      `json:"given_name"`
    FamilyName      string      `json:"family_name"`
    Picture         string      `json:"picture"`
  }
  
  if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
    return nil, fmt.Errorf("failed to parse google userinfo: %w", err)
  }

  user, err := o.userSvc.GetUserByEmail(info.Email)
  if err != nil {
    newUser := &models.User{
      Email: info.Email,
      FirstName: info.GivenName,
      LastName:  info.FamilyName,
    }
    user, err := o.userSvc.CreateUser(newUser)
    if err != nil {
      return nil, fmt.Errorf("failed to create user: %w", err)
    }
  }
  return user, nil
}


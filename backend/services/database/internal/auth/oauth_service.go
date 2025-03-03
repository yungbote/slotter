package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/google/uuid"
	"github.com/yungbote/slotter/backend/services/database/internal/models"
	"github.com/yungbote/slotter/backend/services/database/internal/services"
)

type OAuthConfig struct {
	GoogleClientID string
	GoogleSecret   string
	RedirectURL    string
	// add Apple/Twitter if needed
}

type OAuthService interface {
	GetGoogleLoginURL(state string) string
	HandleGoogleCallback(code, state string) (*models.User, error)
}

type oauthService struct {
	googleConfig *oauth2.Config
	userSvc      services.USvc
}

func NewOAuthService(cfg OAuthConfig, userSvc services.USvc) (OAuthService, error) {
	if cfg.GoogleClientID == "" || cfg.GoogleSecret == "" {
		return nil, fmt.Errorf("missing Google OAuth creds")
	}
	gConf := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{"openid", "profile", "email"},
		RedirectURL:  cfg.RedirectURL,
	}
	return &oauthService{googleConfig: gConf, userSvc: userSvc}, nil
}

func (o *oauthService) GetGoogleLoginURL(state string) string {
	return o.googleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (o *oauthService) HandleGoogleCallback(code, state string) (*models.User, error) {
	token, err := o.googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange Google OAuth code: %w", err)
	}
	client := o.googleConfig.Client(context.Background(), token)

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get userinfo from Google: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Google userinfo returned status %d", resp.StatusCode)
	}

	var info struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		VerifiedEmail bool   `json:"verified_email"`
		Name          string `json:"name"`
		GivenName     string `json:"given_name"`
		FamilyName    string `json:"family_name"`
		Picture       string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("failed to parse google userinfo: %w", err)
	}

	// Check if user already exists by email
	user, err := o.userSvc.GetUserByEmail(info.Email)
	if err == nil && user != nil {
		// user found
		return user, nil
	}

	// Otherwise create new user stub
	newUser := models.User{
		ID:        uuid.New(),
		Email:     info.Email,
		FirstName: info.GivenName,
		LastName:  info.FamilyName,
		// for OAuth, you can store a dummy password or handle differently
		Password:  "", 
	}
	created, err := o.userSvc.CreateUser(newUser)
	if err != nil {
		return nil, fmt.Errorf("failed to create user from Google profile: %w", err)
	}
	return created, nil
}


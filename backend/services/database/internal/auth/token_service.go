package auth

import (
  "time"
  "fmt"
  "errors"
  "os"
  "github.com/golang-jwt/jwt/v4"
  "github.com/google/uuid"
)

type TokenService interface {
  GenerateToken(userID uuid.UUID, email string) (string, error)
  ValidateToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type tokenService struct {
  secretKey       string
  issuer          string
}

func NewTokenService() (TokenService, error) {
  secret := os.Getenv("JWT_SECRET")
  if secret == "" {
    return nil, fmt.Errorf("JWT_SECRET not set")
  }
  issuer := os.Getenv("JWT_ISSUER")
  if issuer == "" {
    issuer = "slotter"
  }
  return &tokenService{
    secretKey:    secret,
    issuer:       issuer,
  }, nil
}

func (t *tokenService) GenerateToken(userID uuid.UUID, email string) (string, error) {
  now := time.Now()
  claims := jwt.RegisteredClaims{
    Issuer:         t.issuer,
    Subject:        userID.String(),
    ExpiresAt:      jwt.NewNumericDate(now.Add(time.Hour * 24 * 7)),
    IssuedAt:       jwt.NewNumericDate(now),
  }
  token := jwt.NewWithClaims(jwt.SigningMethodHS246, claims)
  signed, err := token.SignedString([]byte(t.secretKey))
  if err != nil {
    return "", err
  }
  return signed, nil
}

func (t *tokenService) ValidateToken(tokenStr string) (*jwt.RegisteredClaims, error) {
  parsedToken, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
    return []byte(t.secretKey), nil
  })
  if err != nil || !parsedToken.Valid {
    return nil, errors.New("invalid token")
  }
  claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
  if !ok {
    return nil, errors.New("invalid token claims")
  }
  return claims, nil
}

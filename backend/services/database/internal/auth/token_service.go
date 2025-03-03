package services

import (
    "errors"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
)

// TokenService deals with short-lived JWT "access tokens" only.
// We'll handle refresh tokens in RefreshTokenService.
type TokenService interface {
    GenerateAccessToken(userID uuid.UUID, email string) (string, error)
    ValidateAccessToken(tokenStr string) (*jwt.RegisteredClaims, error)
}

type tokenService struct {
    secretKey string
    issuer    string
    expMin    int // how many minutes the access token is valid
}

// NewTokenService loads JWT_SECRET, JWT_ISSUER, and optionally ACCESS_TOKEN_EXP_MIN from environment.
func NewTokenService() (TokenService, error) {
    secret := os.Getenv("JWT_SECRET")
    if secret == "" {
        return nil, fmt.Errorf("JWT_SECRET not set")
    }
    issuer := os.Getenv("JWT_ISSUER")
    if issuer == "" {
        issuer = "slotter"
    }

    expStr := os.Getenv("ACCESS_TOKEN_EXP_MIN")
    if expStr == "" {
        expStr = "15" // default 15 minutes
    }
    expMin, err := strconv.Atoi(expStr)
    if err != nil {
        expMin = 15
    }

    return &tokenService{
        secretKey: secret,
        issuer:    issuer,
        expMin:    expMin,
    }, nil
}

func (t *tokenService) GenerateAccessToken(userID uuid.UUID, email string) (string, error) {
    now := time.Now()
    // Use RegisteredClaims for standard fields: subject, expiry, etc.
    claims := jwt.RegisteredClaims{
        Issuer:    t.issuer,
        Subject:   userID.String(),
        ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(t.expMin) * time.Minute)),
        IssuedAt:  jwt.NewNumericDate(now),
    }

    // If you want email as a custom claim, create a custom struct or embed it in map claims.
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signed, err := token.SignedString([]byte(t.secretKey))
    if err != nil {
        return "", err
    }
    return signed, nil
}

func (t *tokenService) ValidateAccessToken(tokenStr string) (*jwt.RegisteredClaims, error) {
    parsedToken, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(t.secretKey), nil
    })
    if err != nil || !parsedToken.Valid {
        return nil, errors.New("invalid or expired access token")
    }

    claims, ok := parsedToken.Claims.(*jwt.RegisteredClaims)
    if !ok {
        return nil, errors.New("invalid token claims")
    }
    return claims, nil
}


package services

import (
    "context"
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "os"
    "strconv"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/google/uuid"
)

// RefreshTokenService manages creation, storage, validation, and invalidation
// of refresh tokens in Redis.
type RefreshTokenService interface {
    GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error)
    ValidateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error
    InvalidateRefreshToken(ctx context.Context, userID uuid.UUID) error
}

type refreshTokenService struct {
    rdb              *redis.Client
    refreshTokenExp  time.Duration
}

// NewRefreshTokenService configures an expiration from env REFRESH_TOKEN_EXP_HOURS,
// defaults to 24 hours if missing.
func NewRefreshTokenService(rdb *redis.Client) (RefreshTokenService, error) {
    expStr := os.Getenv("REFRESH_TOKEN_EXP_HOURS")
    if expStr == "" {
        expStr = "24" // default
    }
    hours, err := strconv.Atoi(expStr)
    if err != nil {
        hours = 24
    }

    return &refreshTokenService{
        rdb:             rdb,
        refreshTokenExp: time.Duration(hours) * time.Hour,
    }, nil
}

func (s *refreshTokenService) GenerateRefreshToken(ctx context.Context, userID uuid.UUID) (string, error) {
    // 32 random bytes -> base64 token
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        return "", fmt.Errorf("failed to generate random bytes for refresh token: %w", err)
    }
    refreshToken := base64.URLEncoding.EncodeToString(b)

    key := fmt.Sprintf("refreshToken:%s", userID.String())
    err := s.rdb.Set(ctx, key, refreshToken, s.refreshTokenExp).Err()
    if err != nil {
        return "", fmt.Errorf("failed to store refresh token in redis: %w", err)
    }
    return refreshToken, nil
}

func (s *refreshTokenService) ValidateRefreshToken(ctx context.Context, userID uuid.UUID, refreshToken string) error {
    key := fmt.Sprintf("refreshToken:%s", userID.String())
    stored, err := s.rdb.Get(ctx, key).Result()
    if err != nil {
        return fmt.Errorf("no stored refresh token or expired: %w", err)
    }
    if stored != refreshToken {
        return fmt.Errorf("refresh token mismatch or invalid")
    }
    return nil
}

func (s *refreshTokenService) InvalidateRefreshToken(ctx context.Context, userID uuid.UUID) error {
    key := fmt.Sprintf("refreshToken:%s", userID.String())
    return s.rdb.Del(ctx, key).Err()
}


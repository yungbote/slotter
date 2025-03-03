package avatar

import (
  "context"
  "fmt"
  "io"
  "net/http"
  "os"
  "github.com/google/uuid"
  s3service "github.com/yungbote/slotter/backend/services/database/internal/services/s3"
)

type AvatarService interface {
  GenerateAvatar(ctx context.Context, sprite, seed string) ([]byte, error)
  GenerateAndUploadAvatar(ctx context.Context, sprite, seed string) (string, error)
}

type avatarService struct {
  s3              s3service.S3Service
  dicebearURL     string
}

func NewAvatarService(s3 s3service.S3Service) AvatarService {
  dicebearURL := os.Getenv("DICEBEAR_URL")
  if dicebearURL == "" {
    dicebearURL = "http://localhost:3000"
  }
  return &avatarService{
    s3:             s3,
    dicebearURL:    dicebearURL,
  }
}

func (a *avatarService) GenerateAvatar(ctx context.Context, sprite, seed string) ([]byte, error) {
  url := fmt.Sprintf("%s/png/%s?seed=%s", a.dicebearURL, sprite, seed)
  req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
  if err != nil {
    return nil, fmt.Errorf("failed to create request to dicebear: %w", err)
  }
  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    return nil, fmt.Errorf("failed to call dicebear endpoint: %w", err)
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("dicebear returned non-200 status: %d", resp.StatusCode)
  }
  data, err := io.ReadAll(resp.Body)
  if err != nil {
    return nil, fmt.Errorf("failed to read dicebear response: %w", err)
  }
  return data, nil
}

func (a *avatarService) GenerateAndUploadAvatar(ctx context.Context, sprite, seed string) (string, error) {
  avatarBytes, err := a.GenerateAvatar(ctx, sprite, seed)
  if err != nil {
    return "", err
  }
  fileName := fmt.Sprintf("dicebear-%s.png", uuid.NewString())
  contentType := "image/png"
  url, err := a.s3.UploadImage(ctx, fileName, avatarBytes, contentType)
  if err != nil {
    return "", fmt.Errorf("failed to upload avatar to s3: %w", err)
  }
  return url, nil
}



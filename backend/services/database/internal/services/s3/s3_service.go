package s3

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

// S3Service is an interface for all S3 file operations you need.
type S3Service interface {
	UploadImage(ctx context.Context, fileName string, fileData []byte, contentType string) (string, error)
	UploadFile(ctx context.Context, filePath string, data []byte) (string, error)
	UpdateImage(ctx context.Context, existingURL string, newFileData []byte, newContentType string) (string, error)
	UpdateFile(ctx context.Context, existingURL string, newFilePath string) (string, error)
	RetrieveFile(ctx context.Context, fileKey string) ([]byte, error)
}

// s3Service implements S3Service.
type s3Service struct {
	client     *awsS3.Client
	bucketName string
	region     string
}

// NewS3Service loads environment variables and returns an S3Service.
func NewS3Service() (S3Service, error) {
	bucketName := os.Getenv("S3_BUCKET")
	if bucketName == "" {
		return nil, fmt.Errorf("S3_BUCKET environment variable not set")
	}
	region := os.Getenv("AWS_REGION")
	if region == "" {
		region = "us-east-1"
	}

	// Load AWS config (credentials, region, etc.)
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	svc := awsS3.NewFromConfig(cfg)

	return &s3Service{
		client:     svc,
		bucketName: bucketName,
		region:     region,
	}, nil
}

// UploadImage uploads raw image data (in memory) and returns the publicly accessible URL.
func (s *s3Service) UploadImage(ctx context.Context, fileName string, fileData []byte, contentType string) (string, error) {
	uniqueID := uuid.New().String()

	ext := filepath.Ext(fileName)
	if ext == "" {
		// If we can guess from mime type
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			ext = exts[0]
		} else {
			ext = ".bin"
		}
	}

	key := fmt.Sprintf("uploads/%s%s", uniqueID, ext)

	_, err := s.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        bytes.NewReader(fileData),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}

	// Construct the public URL. (For a public-read object, you can often just link directly)
	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)
	return fileURL, nil
}

// UploadFile uploads a file from disk (given filePath). 
// You might actually want to pass the raw data directly from your code. 
func (s *s3Service) UploadFile(ctx context.Context, filePath string, data []byte) (string, error) {
	// If you want to open the file from disk, do so. We also accept 'data' here, so either approach works.
	fileName := filepath.Base(filePath)
	contentType := detectContentType(fileName)

	uniqueID := uuid.New().String()
	ext := filepath.Ext(fileName)
	if ext == "" {
		if exts, _ := mime.ExtensionsByType(contentType); len(exts) > 0 {
			ext = exts[0]
		} else {
			ext = ".bin"
		}
	}

	key := fmt.Sprintf("uploads/%s%s", uniqueID, ext)

	_, err := s.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file to s3: %w", err)
	}

	fileURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.bucketName, s.region, key)
	return fileURL, nil
}

// UpdateImage overwrites an existing object in S3 with newFileData, returning the same or new URL.
func (s *s3Service) UpdateImage(ctx context.Context, existingURL string, newFileData []byte, newContentType string) (string, error) {
	key, err := parseS3KeyFromURL(existingURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse existing S3 URL: %w", err)
	}

	_, err = s.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        bytes.NewReader(newFileData),
		ContentType: &newContentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update file in s3: %w", err)
	}

	// Return the same S3 URL, since the object key did not change:
	return existingURL, nil
}

// UpdateFile overwrites an existing S3 object with a new local file's contents.
func (s *s3Service) UpdateFile(ctx context.Context, existingURL, newFilePath string) (string, error) {
	key, err := parseS3KeyFromURL(existingURL)
	if err != nil {
		return "", err
	}

	data, err := os.ReadFile(newFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read local file: %w", err)
	}
	contentType := detectContentType(newFilePath)

	_, err = s.client.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:      &s.bucketName,
		Key:         &key,
		Body:        bytes.NewReader(data),
		ContentType: &contentType,
		ACL:         types.ObjectCannedACLPublicRead,
	})
	if err != nil {
		return "", fmt.Errorf("failed to update file in s3: %w", err)
	}

	return existingURL, nil
}

// RetrieveFile fetches the raw bytes of an S3 object by key. 
// If you only have a URL, parse out the key with parseS3KeyFromURL.
func (s *s3Service) RetrieveFile(ctx context.Context, fileKey string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &awsS3.GetObjectInput{
		Bucket: &s.bucketName,
		Key:    &fileKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object data: %w", err)
	}
	return data, nil
}

// parseS3KeyFromURL extracts the object key from an S3 URL like "https://bucket.s3.us-east-1.amazonaws.com/uploads/abc123.png"
func parseS3KeyFromURL(fileURL string) (string, error) {
	u, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}
	// The path typically starts with "/". We remove the leading slash.
	key := strings.TrimPrefix(u.Path, "/")
	return key, nil
}

// detectContentType is a naive approach. Adjust as you wish, or read a small portion to guess the real MIME type.
func detectContentType(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	switch ext {
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	default:
		return "application/octet-stream"
	}
}


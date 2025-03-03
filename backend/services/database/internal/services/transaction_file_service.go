package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)


type TFSvc interface {
  //GENERAL CRUD
  CreateTransactionFile(file models.TransactionFile) (*models.TransactionFile, error)
  UpdateTransactionFileName(fileID uuid.UUID, newName string) error
  UpdateTransactionFileExtension(fileID uuid.UUID, newExt string) error
  UpdateTransactionFilePathURL(fileID uuid.UUID, newPathURL string) error
  GetTransactionFileByID(fileID uuid.UUID) (*models.TransactionFile, error)
  DeleteTransactionFile(fileID uuid.UUID) error

  LinkToLocation(fileID, locationID uuid.UUID) error
  UnlinkFromLocation(fileID, locationID uuid.UUID) error
  LinkToItem(fileID, itemID uuid.UUID) error
  UnlinkFromItem(fileID, itemID uuid.UUID) error

  ListTransactionFiles(f repos.TransactionFileFilter) ([]*models.TransactionFile, error)
}

type tfSvc struct {
  repo            repos.TFRepo
}

func NewTFSvc(repo repos.TFRepo) TFSvc {
  return &tfSvc{repo: repo}
}

func (s *tfSvc) CreateTransactionFile(file models.TransactionFile) (*models.TransactionFile, error) {
  if file.FileName == "" {
    return nil, fmt.Errorf("transaction file name is required")
  }
  created, err := s.repo.Create(file)
  if err != nil {
    return nil, fmt.Errorf("failed to create transaction file: %w", err)
  }
  return created, nil
}

func (s *tfSvc) UpdateTransactionFileName(fileID uuid.UUID, newName string) error {
  if fileID == uuid.Nil || newName == "" {
    return fmt.Errorf("invalid input to update file name")
  }
  if err := s.repo.UpdateName(fileID, newName); err != nil {
    return fmt.Errorf("update name repo error: %w", err)
  }
  return nil
}

func (s *tfSvc) UpdateTransactionFileExtension(fileID uuid.UUID, newExt string) error {
  if fileID == uuid.Nil || newExt == "" {
    return fmt.Errorf("invalid input to update file extension")
  }
  if err := s.repo.UpdateFileExtension(fileID, newExt); err != nil {
    return fmt.Errorf("update extension repo error: %w", err)
  }
  return nil
}

func (s *tfSvc) UpdateTransactionFilePathURL(fileID uuid.UUID, newPathURL string) error {
  if fileID == uuid.Nil || newPathURL == "" {
    return fmt.Errorf("invalid input to update file path url")
  }
  if err := s.repo.UpdateFilePathURL(fileID, newPathURL); err != nil {
    return fmt.Errorf("update file path url repo error: %w", err)
  }
  return nil
}

func (s *tfSvc) GetTransactionFileByID(fileID uuid.UUID) (*models.TransactionFile, error) {
  if fileID == uuid.Nil {
    return nil, fmt.Errorf("Invalid FileID")
  }
  file, err := s.repo.GetByID(fileID)
  if err != nil {
    return nil, fmt.Errorf("repo getByID error: %w", err)
  }
  return file, nil
}

func (s *tfSvc) DeleteTransactionFile(fileID uuid.UUID) error {
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  if err := s.repo.Delete(fileID); err != nil {
    return fmt.Errorf("failed to delete transaction file: %w", err)
  }
  return nil
}

func (s *tfSvc) LinkToLocation(fileID, locationID uuid.UUID) error {
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  return s.repo.LinkToLocation(fileID, locationID)
}

func (s *tfSvc) UnlinkFromLocation(fileID, locationID uuid.UUID) error {
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  if locationID == uuid.Nil {
    return fmt.Errorf("Invalid LocationID")
  }
  return s.repo.UnlinkFromLocation(fileID, locationID)
}

func (s *tfSvc) LinkToItem(fileID, itemID uuid.UUID) error {
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  return s.repo.LinkToItem(fileID, itemID)
}

func (s *tfSvc) UnlinkFromItem(fileID, itemID uuid.UUID) error {
  if fileID == uuid.Nil {
    return fmt.Errorf("Invalid FileID")
  }
  if itemID == uuid.Nil {
    return fmt.Errorf("Invalid ItemID")
  }
  return s.repo.UnlinkFromItem(fileID, itemID)
}

func (s *tfSvc) ListTransactionFiles(f repos.TransactionFileFilter) ([]*models.TransactionFile, error) {
  files, err := s.repo.ListTransactionFiles(f)
  if err != nil {
    return nil, fmt.Errorf("list transaction files files repo error: %w", err)
  }
  return files, nil
}

                      

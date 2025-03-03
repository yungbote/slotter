package services

import (
  "fmt"
  "github.com/google/uuid"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repos"
)

type CSvc interface {
  CreateCompany(company models.Company) (*models.Company, error)
  UpdateCompanyName(companyID uuid.UUID, newName string) error
  UpdateCompanyAvatarURL(companyID uuid.UUID, newAvatarURL string) error
  GetCompanyByID(companyID uuid.UUID) (*models.Company, error)
  DeleteCompany(companyID uuid.UUID) error
}

type cSvc struct {
  repo          repos.CRepo
}

func NewCSvc(repo repos.CompanyRepo) CSvc {
  return &cSvc{repo: repo}
}

func (s *cSvc) CreateCompany(company models.Company) (*models.Company, error) {
  if company.Name == "" {
    return nil, fmt.Errorf("Company name cannot be empty")
  }
  if company.AvatarURL == "" {
    return nil, fmt.Errorf("Company AvatarURL cannot be empty")
  }
  created, err := s.repo.Create(company)
  if err != nil {
    return nil, fmt.Errorf("Failed to create company: %w", err)
  }
  return created, nil
}

func (s *cSvc) UpdateCompanyName(companyID uuid.UUID, newName string) error {
  if companyID == uuid.Nil {
    return fmt.Errorf("invalid company ID")
  }
  if newName == "" {
    return fmt.Errorf("newName cannot be empty")
  }
  if err := s.repo.UpdateName(companyID, newName); err != nil {
    return fmt.Errorf("Failed to update company name: %w", err)
  }
  return nil
}

func (s *cSvc) UpdateCompanyAvatarURL(companyID uuid.UUID, newAvatarURL string) error {
  if companyID == uuid.Nil {
    return fmt.Errorf("invalid company ID")
  }
  if newAvatarURL == "" {
    return fmt.Errorf("newAvatarURL cannot be empty")
  }
  if err := s.repo.UpdateAvatarURL(companyID, newAvatarURL); err != nil {
    return fmt.Errorf("Failed to update avatar url: %w", err)
  }
  return nil
}

func (s *cSvc) GetCompanyByID(companyID uuid.UUID) (*models.Company, error) {
  if companyID == uuid.Nil {
    return nil, fmt.Errorf("Invalid company ID")
  }
  c, err := s.repo.GetByID(companyID)
  if err != nil {
    return nil, fmt.Errorf("could not fetch company: %w", err)
  }
  return c, nil
}

func (s *cSvc) DeleteCompany(companyID uuid.UUID) error {
  if companyID == uuid.Nil {
    return fmt.Errorf("Invalid company ID")
  }
  if err := s.repo.Delete(companyID); err != nil {
    return fmt.Errorf("Failed to delete company: %w", err)
  }
  return nil
}







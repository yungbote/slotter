package services

import (
  "fmt"
  "github.com/yungbote/slotter/backend/services/database/internal/models"
  "github.com/yungbote/slotter/backend/services/database/internal/repositories"
)

type CompanyService interface {
  CreateCompany(input *models.Company) (*models.Company, error)
  GetCompanyID(id uint) (*models.Company, error)
  UpdateCompany(existing *models.Company, updates *models.Company) (*models.Company, error)
  DeleteCompany(existing *models.Company) error
}

type companyService struct {
  repo repositories.CompanyRepository
}

func NewCompanyService(repo repositories.CompanyRepository) CompanyService {
  return &companyService{repo: repo}
}

func (s *companyService) CreateCompany(input *models.Company) (*models.Company, error) {
  if input.Name == "" {
    return nil, fmt.Errorf("company name cannot be empty")
  }
  err := s.repo.Create(input)
  if err != nil {
    return nil, err
  }
  return input, nil
}

func (s *companyService) GetCompanyByID(id uint) (*models.Company, error) {
  return s.repo.GetByID(id)
}

func (s *companyService) UpdateCompany(existing *models.Company, updates *models.Company) (*models.Company, error) {
  existing.Name = updates.Name
  if err := s.repo.Update(existing); err != nil {
    return nil, err
  }
  return existing, nil
}

func (s *companyService) DeleteCompany(existing *models.Company) error {
  return s.repo.Delete(existing)
}

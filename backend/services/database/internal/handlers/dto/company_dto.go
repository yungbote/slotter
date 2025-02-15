package dto

import (
    "github.com/yungbote/slotter/backend/services/database/internal/models"
)

//Input DTO's
type CreateCompanyRequest struct {
    Name        string          `json:"name" binding:"required"`
}

type UpdateCompanyRequest struct {
    Name        *string         `json:"name,omitempty"`
}

type GetCompanyByIDRequest struct {
    ID          uint            `json:"id" binding:"required"`
}

type DeleteCompanyRequest struct {
    ID          uint            `json:"id" binding:"required"`
}

//Output DTO's
type CompanyResponse struct {
    ID      uint    `json:"id"`
    Name    string  `json:"name"`
}

//Mapping Functions
func (req *CreateCompanyRequest) ToModel() *models.Company {
    return &models.Company{
        Name: req.Name,
    }
}

func (req *UpdateCompanyRequest) ApplyToModel(c *models.Company) {
    if req.Name != nil {
        c.Name = *req.Name
    }
}

func ToCompanyResponse(c *models.Company) *CompanyResponse {
    if c == nil {
        return nil
    }
    return &CompanyResponse{
        ID:     c.ID,
        Name:   c.Name,
    }
}

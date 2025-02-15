package dto

import (
    "github.com/yungbote/slotter/backend/services/database/internal/models"
)

//Input DTO's
type CreateUserRequest struct {
    Email       string          `json:"email"       binding:"required,email"`
    Password    string          `json:"password"    binding:"required"`
    FirstName   string          `json:"first_name"  binding:"required"`
    LastName    string          `json:"last_name"   binding:"required"`
    CompanyID   uint            `json:"company_id"  binding:"required"`
}

type UpdateUserRequest struct {
    Email       *string         `json:"email,omitempty"     binding:"omitempty,email"`
    Password    *string         `json:"password,omitempty"`
    FirstName   *string         `json:"first_name,omitempty"`
    LastName    *string         `json:"last_name,omitempty"`
    CompanyID   *uint           `json:"company_id,omitempty"`
}

type GetUserByIDRequest struct {
    ID          uint          `json:"id"      binding:"required"`
}

type GetUserByEmailRequest struct {
    Email       string          `json:"email"   binding:"required,email"`
}

type DeleteUserRequest struct {
    ID          uint            `json:"id" binding:"required"`

}

type ListUsersByCompanyIDRequest struct {
    CompanyID   uint    `json:"company_id" binding:"required"`
}

type CountUsersByCompanyIDRequest struct {
    CompanyID   uint    `json:"company_id" binding:"required"`
}

//Output DTO's
type UserResponse struct {
    ID          uint        `json:"id"`
    Email       string      `json:"email"`
    FirstName   string      `json:"first_name"`
    LastName    string      `json:"last_name"`
    CompanyID   uint        `json:"company_id"`
}

//Mapping Functions
func (req *CreateUserRequest) ToModel() *models.User {
    return &models.User{
        Email:              req.Email,
        PasswordHash:       req.Password,
        FirstName:          req.FirstName,
        LastName:           req.LastName,
        CompanyID:          req.CompanyID,
    }
}

func (req *UpdateUserRequest) ApplyToModel(u *models.User) {
    if req.Email != nil {
        u.Email = *req.Email
    }
    if req.Password != nil {
        u.PasswordHash = *req.Password
    }
    if req.FirstName != nil {
        u.FirstName = *req.FirstName
    }
    if req.LastName != nil {
        u.LastName = *req.LastName
    }
    if req.CompanyID != nil {
        u.CompanyID = *req.CompanyID
    }
}

func ToUserResponse(u *models.User) *UserResponse {
    if u == nil {
        return nil
    }
    return &UserResponse{
        ID:         u.ID,
        Email:      u.Email,
        FirstName:  u.FirstName,
        LastName:   u.LastName,
        CompanyID:  u.CompanyID,
    }
}

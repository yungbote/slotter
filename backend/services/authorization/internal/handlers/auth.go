package handlers

import (
  "bytes"
  "encoding/json"
  "errors"
  "fmt"
  "log"
  "net/http"
  "os"
  "strings"
  "github.com/gin-gonic/gin"
  "golang.org/x/crypto/bcrypt"
  "github.com/yungbote/slotter/backend/services/authorization/internal/tokens"
  "time"
)

var dbServiceURL = getEnv("DATABASE_SERVICE_URL", "http://slotter-database:8080")

type SignUpRequest struct {
  Email       string          `json:"email" binding:"required"`
  Password    string          `json:"password" binding:"required"`
  FullName    string          `json:"fullName" binding:"required"`
  CompanyID   *uint           `json:"companyId,omitempty"`
  CompanyName string          `json:"companyName,omitempty"`
}

func SignUp(c *gin.Context) {
  var req SignUpRequest
  if err := c.ShouldBindJSON(&req); err != nil {
    log.Println("ERROR: SignUp payload:", err)
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid sign-up data"})
    return
  }
  hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
  if err != nil {
    log.Println("ERROR: bcrypt:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not hash password"})
    return
  }
  var companyID uint
  if req.CompanyID != nil {
    exists, err := checkCompanyExists(*req.CompanyID)
    if err != nil {
      log.Println("ERROR: checkCompanyExists:", err)
      c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error checking company exists"})
      return
    }
    if !exists {
      c.JSON(http.StatusBadRequest, gin.H{"error": "No company found with that id"})
      return
    }
    companyID = *req.CompanyID
  } else if req.CompanyName != "" {
    newCompany, err := createCompany(req.CompanyName)
    if err != nil {
      log.Println("ERROR: createCompany:", err)
      c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create company"})
      return
    }
    companyID = newCompany.ID
  }
  createdUser, err := createUser(req.Email, string(hashed), req.FullName, companyID)
  if err != nil {
    log.Println("ERROR: createUser:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create user"})
    return
  }
  c.JSON(http.StatusCreated, createdUser)
}

type User struct {
  ID            uint        `json:"id"`
  Email         string      `json:"email"`
  PasswordHash  string      `json:"passwordHash"`
  FullName      string      `json:"fullName"`
  Status        string      `json:"status"`
  CompanyID     uint        `json:"companyId"`
}

func Login(c *gin.Context) {
  var req struct {
    Email     string    `json:"email" binding:"required"`
    Password  string    `json:"password" binding:"required"`
  }
  if err := c.ShouldBindJSON(&req); err != nil {
    c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login data"})
    return
  }
  user, err := getUserByEmail(req.Email)
  if err != nil {
    log.Println("ERROR: getUserByEmail:", err)
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
    return
  }
  if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
    c.JSON(http.StatusUnauthorized, gin.H{"error": "Incorrect password"})
    return
  }
  token, err := tokens.GenerateJWT(user.ID, time.Hour)
  if err != nil {
    log.Println("ERROR: GenerateJWT:", err)
    c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
    return
  }
  c.JSON(http.StatusOK, gin.H{
    "access_token": token,
    "user":         user,
  })
}

func checkCompanyExists(id uint) (bool, error) {
  url := fmt.Sprintf("%s/v1/companies/%d", dbServiceURL, id)
  resp, err := http.Get(url)
  if err != nil {
    return false, err
  }
  defer resp.Body.Close()
  if resp.StatusCode == http.StatusNotFound {
    return false, nil
  }
  if resp.StatusCode != http.StatusOK {
    return false, fmt.Errorf("ERROR: Unexpected status: %d", resp.StatusCode)
  }
  return true, nil
}

type Company struct {
  ID    uint    `json:"id"`
  Name  string  `json:"name"`
}

func createCompany(name string) (*Company, error) {
  url := fmt.Sprintf("%s/v1/companies", dbServiceURL)
  payload := map[string]string{
    "name": name,
  }
  body, _ := json.Marshal(payload)
  resp, err := http.Post(url, "application/json", bytes.NewReader(body))
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
    return nil, fmt.Errorf("ERROR: createCompany got status %d", resp.StatusCode)
  }
  var comp Company
  if err := json.NewDecoder(resp.Body).Decode(&comp); err != nil {
    return nil, err
  }
  return &comp, nil
}

func createUser(email, passHash, fullName string, companyID uint) (*User, error) {
  url := fmt.Sprintf("%s/v1/users", dbServiceURL)
  payload := map[string]interface{}{
    "email":        email
    "passwordHash": passHash,
    "fullName":     fullName,
    "companyID":    companyID,
  }
  body, _ := json.Marshal(payload)
}

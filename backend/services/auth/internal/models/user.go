package models

type User struct {
  ID              uint      `json:"id"`
  Email           string    `json:"email"`
  PasswordHash    string    `json:"passwordHash"`
  FullName        string    `json:"fullName"`
}

package repositories

import (
  "errors"
  "fmt"
  "github.com/jackc/pgconn"
  "gorm.io/gorm"
)

type DomainErrorCode string

const (
  ErrCodeNotFound         DomainErrorCode = "NOT_FOUND"
  ErrCodeDuplicate        DomainErrorCode = "DUPLICATE"
  ErrCodeForeignKey       DomainErrorCode = "FOREIGN_KEY_VIOLATION"
  ErrCodeCheckViolation   DomainErrorCode = "CHECK_VIOLATION"
  ErrCodeValidation       DomainErrorCode = "VALIDATION"
  ErrCodeUnknown          DomainErrorCode = "UNKNOWN"
)

type DomainError struct {
  Code          DomainErrorCode
  Message       string
  Op            string
  Underlying    error
}

func (e *DomainError) Error() string {
  return fmt.Sprintf("[%s] %s: %s", e.Op, e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
  return e.Underlying
}

func NewDomainError(op string, code DomainErrorCode, msg string, underlying error) error {
  return &DomainError{
    Code:         code,
    Message:      msg,
    Op:           op,
    Underlying:   underlying,
  }
}

func ParseDBError(op string, err error) error {
  if err == nil {
    return nil
  }
  if errors.Is(err, gorm.ErrRecordNotFound) {
    return NewDomainError(op, ErrCodeNotFound, "record not found", err)
  }
  var pgErr *pgconn.PgError
  if errors.As(err, &pgErr) {
    return mapPostgresError(op, pgErr)
  }
  return NewDomainError(op, ErrCodeUnknown, err.Error(), err)
}

func mapPostgresError(op string, pgErr *pgconn.PgError) error {
  switch pgErr.Code {
  case "23503":
    return NewDomainError(op, ErrCodeForeignKey, "foreign key constraint violation", pgErr)
  case "23505":
    return NewDomainError(op, ErrCodeDuplicate, "duplicate key", pgErr)
  case "23514":
    return NewDomainError(op, ErrCodeCheckViolation, "check constraint violation", pgErr)
  case "23502":
    return NewDomainError(op, ErrCodeValidation, "column cannot be null", pgErr)
  default:
    return NewDomainError(op, ErrCodeUnknown, pgErr.Message, pgErr)
  }
}

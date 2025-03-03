package errors

import (
  "fmt"
)

type CodeError struct {
  Code    string
  Msg     string
  Err     error
  Entity  string
}

func (ce *CodeError) Error() string {
  if ce.Err != nil {
    return fmt.Sprintf("[%s] %s: %v", ce.Code, ce.Msg, ce.Err)
  }
  return fmt.Sprintf("[%s] %s", ce.Code, ce.Msg)
}

func (cd *CodeError) Unwrap() error {
  return ce.Error
}

func NewCodeError(code, msg string, err error) error {
  return &CodeError{Code: code, Msg: msg, Err: err, Entity: entity}
}

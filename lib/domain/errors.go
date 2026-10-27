// Package domain contains values and error categories shared across LGC.
package domain

import (
	"errors"
	"fmt"
)

// Stable error categories allow callers to classify failures without parsing
// error messages.
var (
	ErrUsage              = errors.New("usage")
	ErrUnauthenticated    = errors.New("unauthenticated")
	ErrInvalidLanguage    = errors.New("invalid language")
	ErrO2Unsupported      = errors.New("O2 unsupported")
	ErrSubmissionRejected = errors.New("submission rejected")
	ErrCaptchaRequired    = errors.New("CAPTCHA required")
	ErrMalformedUpstream  = errors.New("malformed upstream response")
	ErrJudgeTimeout       = errors.New("judge timeout")
	ErrPlugin             = errors.New("plugin failure")
	ErrNotImplemented     = errors.New("not implemented")
)

// Error wraps a stable category with safe operation context.
type Error struct {
	Category  error
	Operation string
	Err       error
}

// Error implements error.
func (e *Error) Error() string {
	switch {
	case e.Operation != "" && e.Err != nil:
		return fmt.Sprintf("%s: %v", e.Operation, e.Err)
	case e.Operation != "" && e.Category != nil:
		return fmt.Sprintf("%s: %v", e.Operation, e.Category)
	case e.Err != nil:
		return e.Err.Error()
	case e.Category != nil:
		return e.Category.Error()
	default:
		return "lgc error"
	}
}

// Unwrap exposes both the stable category and the underlying cause.
func (e *Error) Unwrap() []error {
	if e.Category == nil {
		return []error{e.Err}
	}
	if e.Err == nil {
		return []error{e.Category}
	}
	return []error{e.Category, e.Err}
}

// Wrap associates a stable category with safe operation context.
func Wrap(category error, operation string, cause error) error {
	return &Error{
		Category:  category,
		Operation: operation,
		Err:       cause,
	}
}

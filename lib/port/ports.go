// Package port declares the stable capability interfaces implemented by
// plugins and consumed by LGC core workflows.
package port

import (
	"context"
	"net/http"

	"github.com/haomingbai/lgc/lib/domain"
)

// AuthorizedRequest is a same-origin request carrying session authorization.
type AuthorizedRequest struct {
	Method string
	Path   string
	Header http.Header
	Body   []byte
}

// AuthorizedResponse is a bounded response returned to a selected provider.
type AuthorizedResponse struct {
	StatusCode int
	Header     http.Header
	Body       []byte
}

// AuthorizedRequester is the narrow session capability exposed to providers.
type AuthorizedRequester interface {
	Do(ctx context.Context, request AuthorizedRequest) (AuthorizedResponse, error)
}

// CaptchaSolver converts one fetched challenge into one answer.
type CaptchaSolver interface {
	Solve(ctx context.Context, challenge domain.CaptchaChallenge) (string, error)
}

// SessionStore persists an opaque, versioned session envelope.
type SessionStore interface {
	Load(ctx context.Context, profile string) ([]byte, bool, error)
	Save(ctx context.Context, profile string, envelope []byte) error
	Delete(ctx context.Context, profile string) error
}

// RecordUpdateSource supplies advisory judge-record updates.
type RecordUpdateSource interface {
	Follow(ctx context.Context, rid string) (<-chan domain.RecordUpdate, error)
	Close() error
}

// ProblemLoader loads optional problem content.
type ProblemLoader interface {
	Load(ctx context.Context, problemID string) (domain.Problem, error)
}

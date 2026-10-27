// Package multimodal provides the multimodal-model CAPTCHA solver skeleton.
package multimodal

import (
	"context"
	"time"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

const (
	// ProviderID is the stable configuration identifier for this provider.
	ProviderID = "lgc.captcha.multimodal"
	// APIVersion is the extension API version implemented by this provider.
	APIVersion = 1
)

// Config contains the multimodal CAPTCHA provider settings.
type Config struct {
	Endpoint            string
	Model               string
	APIKeyRef           string
	Timeout             time.Duration
	MaxResponseBytes    int64
	MaxAnswerBytes      int
	MaxOutputTokens     int
	MaxTransportRetries int
}

// Provider registers the multimodal CAPTCHA solver.
type Provider struct{}

// NewProvider creates the bundled provider entry point.
func NewProvider() *Provider {
	return &Provider{}
}

// Descriptor returns the provider compatibility identity.
func (*Provider) Descriptor() port.Descriptor {
	return port.Descriptor{
		ID:         ProviderID,
		Version:    "0.0.0",
		APIVersion: APIVersion,
		Kind:       port.ProviderCaptcha,
	}
}

// Register registers the provider's typed factory.
func (*Provider) Register(ctx context.Context, registrar port.Registrar) error {
	return domain.ErrNotImplemented
}

// NewSolver opens one configured multimodal solver.
func NewSolver(
	ctx context.Context,
	providerConfig port.ProviderConfig,
) (port.CaptchaSolver, error) {
	return nil, domain.ErrNotImplemented
}

// Solver implements the multimodal CAPTCHA capability.
type Solver struct {
	config Config
}

// Solve sends one fetched challenge to the configured model endpoint.
func (*Solver) Solve(
	ctx context.Context,
	challenge domain.CaptchaChallenge,
) (string, error) {
	return "", domain.ErrNotImplemented
}

var (
	_ port.Provider      = (*Provider)(nil)
	_ port.CaptchaSolver = (*Solver)(nil)
)

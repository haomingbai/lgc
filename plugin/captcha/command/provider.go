// Package command provides the external-command CAPTCHA solver skeleton.
package command

import (
	"context"
	"time"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

const (
	// ProviderID is the stable configuration identifier for this provider.
	ProviderID = "lgc.captcha.command"
	// APIVersion is the extension API version implemented by this provider.
	APIVersion = 1
)

// Config contains the external-command solver settings.
type Config struct {
	Executable     string
	Args           []string
	Timeout        time.Duration
	WorkingDir     string
	EnvAllowlist   []string
	MaxStdoutBytes int64
	MaxStderrBytes int64
}

// Provider registers the external-command CAPTCHA solver.
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

// NewSolver opens one configured external-command solver.
func NewSolver(
	ctx context.Context,
	providerConfig port.ProviderConfig,
) (port.CaptchaSolver, error) {
	return nil, domain.ErrNotImplemented
}

// Solver implements the command-based CAPTCHA capability.
type Solver struct {
	config Config
}

// Solve sends one fetched challenge to the configured trusted program.
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

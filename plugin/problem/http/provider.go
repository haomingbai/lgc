// Package problemhttp provides the optional HTTP problem-loader skeleton.
package problemhttp

import (
	"context"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

const (
	// ProviderID is the stable configuration identifier for this provider.
	ProviderID = "lgc.problem.http"
	// APIVersion is the extension API version implemented by this provider.
	APIVersion = 1
)

// Config contains the optional problem-loader settings.
type Config struct {
	MaxResponseBytes int64
}

// Provider registers the optional HTTP problem loader.
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
		Kind:       port.ProviderProblemLoader,
	}
}

// Register registers the provider's typed factory.
func (*Provider) Register(ctx context.Context, registrar port.Registrar) error {
	return domain.ErrNotImplemented
}

// NewLoader opens one configured problem loader.
func NewLoader(
	ctx context.Context,
	requester port.AuthorizedRequester,
	providerConfig port.ProviderConfig,
) (port.ProblemLoader, error) {
	return nil, domain.ErrNotImplemented
}

// Loader implements optional problem-content loading.
type Loader struct {
	config Config
}

// Load fetches and parses one problem statement.
func (*Loader) Load(
	ctx context.Context,
	problemID string,
) (domain.Problem, error) {
	return domain.Problem{}, domain.ErrNotImplemented
}

var (
	_ port.Provider      = (*Provider)(nil)
	_ port.ProblemLoader = (*Loader)(nil)
)

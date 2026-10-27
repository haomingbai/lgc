// Package websocket provides the advisory judge-update source skeleton.
package websocket

import (
	"context"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

const (
	// ProviderID is the stable configuration identifier for this provider.
	ProviderID = "lgc.recordupdate.websocket"
	// APIVersion is the extension API version implemented by this provider.
	APIVersion = 1
)

// Config contains the WebSocket update-source settings.
type Config struct {
	Endpoint string
}

// Provider registers the WebSocket update source.
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
		Kind:       port.ProviderRecordUpdate,
	}
}

// Register registers the provider's typed factory.
func (*Provider) Register(ctx context.Context, registrar port.Registrar) error {
	return domain.ErrNotImplemented
}

// NewSource opens one configured update source.
func NewSource(
	ctx context.Context,
	requester port.AuthorizedRequester,
	providerConfig port.ProviderConfig,
) (port.RecordUpdateSource, error) {
	return nil, domain.ErrNotImplemented
}

// Source implements an optional pushed-update subscription.
type Source struct {
	config Config
}

// Follow starts one advisory update subscription.
func (*Source) Follow(
	ctx context.Context,
	rid string,
) (<-chan domain.RecordUpdate, error) {
	return nil, domain.ErrNotImplemented
}

// Close closes the source and its owned connection.
func (*Source) Close() error {
	return domain.ErrNotImplemented
}

var (
	_ port.Provider           = (*Provider)(nil)
	_ port.RecordUpdateSource = (*Source)(nil)
)

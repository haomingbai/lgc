// Package file provides the local session-store plugin skeleton.
package file

import (
	"context"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

const (
	// ProviderID is the stable configuration identifier for this provider.
	ProviderID = "lgc.sessionstore.file"
	// APIVersion is the extension API version implemented by this provider.
	APIVersion = 1
)

// Config contains the local session-store settings.
type Config struct {
	Root string
}

// Provider registers the local file session store.
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
		Kind:       port.ProviderSessionStore,
	}
}

// Register registers the provider's typed factory.
func (*Provider) Register(ctx context.Context, registrar port.Registrar) error {
	return domain.ErrNotImplemented
}

// NewStore opens one configured session store.
func NewStore(
	ctx context.Context,
	providerConfig port.ProviderConfig,
) (port.SessionStore, error) {
	return nil, domain.ErrNotImplemented
}

// Store implements atomic, permission-restricted local storage.
type Store struct {
	config Config
}

// Load returns an opaque envelope and distinguishes absence from read failure.
func (*Store) Load(
	ctx context.Context,
	profile string,
) ([]byte, bool, error) {
	return nil, false, domain.ErrNotImplemented
}

// Save atomically replaces one profile envelope.
func (*Store) Save(ctx context.Context, profile string, envelope []byte) error {
	return domain.ErrNotImplemented
}

// Delete removes one profile envelope idempotently.
func (*Store) Delete(ctx context.Context, profile string) error {
	return domain.ErrNotImplemented
}

var (
	_ port.Provider     = (*Provider)(nil)
	_ port.SessionStore = (*Store)(nil)
)

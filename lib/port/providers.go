package port

import (
	"context"
	"encoding/json"
)

// ProviderKind identifies one typed provider capability.
type ProviderKind string

const (
	ProviderCaptcha       ProviderKind = "captcha"
	ProviderSessionStore  ProviderKind = "session_store"
	ProviderRecordUpdate  ProviderKind = "record_update"
	ProviderProblemLoader ProviderKind = "problem_loader"
)

// Descriptor identifies a provider and its extension API compatibility.
type Descriptor struct {
	ID         string
	Version    string
	APIVersion int
	Kind       ProviderKind
}

// ProviderConfig is the provider-scoped, non-secret configuration document.
type ProviderConfig map[string]json.RawMessage

// Factories open selected capabilities after their dependencies are ready.
type (
	CaptchaFactory func(
		ctx context.Context,
		config ProviderConfig,
	) (CaptchaSolver, error)

	SessionStoreFactory func(
		ctx context.Context,
		config ProviderConfig,
	) (SessionStore, error)

	RecordUpdateSourceFactory func(
		ctx context.Context,
		requester AuthorizedRequester,
		config ProviderConfig,
	) (RecordUpdateSource, error)

	ProblemLoaderFactory func(
		ctx context.Context,
		requester AuthorizedRequester,
		config ProviderConfig,
	) (ProblemLoader, error)
)

// Registrar is implemented by the command-side typed provider registry.
type Registrar interface {
	RegisterCaptcha(Descriptor, CaptchaFactory) error
	RegisterSessionStore(Descriptor, SessionStoreFactory) error
	RegisterRecordUpdateSource(Descriptor, RecordUpdateSourceFactory) error
	RegisterProblemLoader(Descriptor, ProblemLoaderFactory) error
}

// Provider is a bundled plugin entry point that registers typed factories.
type Provider interface {
	Descriptor() Descriptor
	Register(ctx context.Context, registrar Registrar) error
}

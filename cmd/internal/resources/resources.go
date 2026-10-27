// Package resources generates live dependencies for exactly one command run.
package resources

import (
	"context"
	"sync"

	"github.com/haomingbai/lgc/cmd/internal/arguments"
	"github.com/haomingbai/lgc/cmd/internal/configuration"
	"github.com/haomingbai/lgc/lib/config"
	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/judge"
	"github.com/haomingbai/lgc/lib/port"
	"github.com/haomingbai/lgc/lib/session"
	"github.com/haomingbai/lgc/lib/submission"
)

type registeredCaptcha struct {
	descriptor port.Descriptor
	factory    port.CaptchaFactory
}

type registeredSessionStore struct {
	descriptor port.Descriptor
	factory    port.SessionStoreFactory
}

type registeredRecordUpdate struct {
	descriptor port.Descriptor
	factory    port.RecordUpdateSourceFactory
}

type registeredProblemLoader struct {
	descriptor port.Descriptor
	factory    port.ProblemLoaderFactory
}

// Registry is the immutable-after-startup, typed provider registry.
type Registry struct {
	mu             sync.RWMutex
	captcha        map[string]registeredCaptcha
	sessionStores  map[string]registeredSessionStore
	recordUpdates  map[string]registeredRecordUpdate
	problemLoaders map[string]registeredProblemLoader
}

// NewRegistry creates an empty typed provider registry.
func NewRegistry() *Registry {
	return &Registry{
		captcha:        make(map[string]registeredCaptcha),
		sessionStores:  make(map[string]registeredSessionStore),
		recordUpdates:  make(map[string]registeredRecordUpdate),
		problemLoaders: make(map[string]registeredProblemLoader),
	}
}

// RegisterCaptcha registers one CAPTCHA factory.
func (*Registry) RegisterCaptcha(
	descriptor port.Descriptor,
	factory port.CaptchaFactory,
) error {
	return domain.ErrNotImplemented
}

// RegisterSessionStore registers one session-store factory.
func (*Registry) RegisterSessionStore(
	descriptor port.Descriptor,
	factory port.SessionStoreFactory,
) error {
	return domain.ErrNotImplemented
}

// RegisterRecordUpdateSource registers one update-source factory.
func (*Registry) RegisterRecordUpdateSource(
	descriptor port.Descriptor,
	factory port.RecordUpdateSourceFactory,
) error {
	return domain.ErrNotImplemented
}

// RegisterProblemLoader registers one problem-loader factory.
func (*Registry) RegisterProblemLoader(
	descriptor port.Descriptor,
	factory port.ProblemLoaderFactory,
) error {
	return domain.ErrNotImplemented
}

// ExecutionResources contains all dependencies owned by one command run.
type ExecutionResources struct {
	Invocation         arguments.Invocation
	Configuration      configuration.Configuration
	LuoguSession       *session.LuoguSession
	Catalog            *config.Catalog
	Captcha            port.CaptchaSolver
	SessionStore       port.SessionStore
	RecordUpdateSource port.RecordUpdateSource
	SubmissionBuilder  *submission.SubmissionBuilder
	SubmissionExecutor *submission.SubmissionExecutor
	JudgeObserver      *judge.JudgeObserver

	mu     sync.Mutex
	closed bool
}

// Close releases all successfully acquired resources.
func (r *ExecutionResources) Close() error {
	return domain.ErrNotImplemented
}

// Generator composes provider, session, catalog, and command input resources.
type Generator struct {
	registry *Registry
}

// NewGenerator creates a resource generator with one typed registry.
func NewGenerator(registry *Registry) *Generator {
	return &Generator{registry: registry}
}

// Generate creates a complete execution bundle or closes every partial
// resource before returning an error.
func (g *Generator) Generate(
	ctx context.Context,
	invocation arguments.Invocation,
	configuration configuration.Configuration,
) (*ExecutionResources, error) {
	return nil, domain.ErrNotImplemented
}

var _ port.Registrar = (*Registry)(nil)

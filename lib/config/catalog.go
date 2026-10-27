// Package config owns runtime configuration snapshots and lookups.
package config

import (
	"context"
	"sync"
	"time"

	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/session"
)

// SnapshotSource identifies where a runtime snapshot came from.
type SnapshotSource string

const (
	SnapshotEmbedded SnapshotSource = "embedded"
	SnapshotRuntime  SnapshotSource = "runtime"
)

// RuntimeSnapshot is one coherent runtime route, language, and status view.
type RuntimeSnapshot struct {
	ObservedAt time.Time
	Source     SnapshotSource
	Routes     map[string]string
	Languages  map[domain.LanguageID]domain.Language
	Statuses   map[int]domain.RecordStatus
}

// Loader refreshes a validated runtime snapshot through an authenticated
// session.
type Loader interface {
	Refresh(
		ctx context.Context,
		session *session.LuoguSession,
	) (RuntimeSnapshot, error)
}

// Catalog publishes immutable runtime snapshots.
type Catalog struct {
	mu      sync.RWMutex
	current RuntimeSnapshot
}

// NewCatalog validates and publishes an initial snapshot.
func NewCatalog(initial RuntimeSnapshot) (*Catalog, error) {
	return nil, domain.ErrNotImplemented
}

// Current returns the current immutable snapshot.
func (c *Catalog) Current() RuntimeSnapshot {
	return RuntimeSnapshot{}
}

// Replace validates and atomically publishes a complete candidate snapshot.
func (c *Catalog) Replace(snapshot RuntimeSnapshot) error {
	return domain.ErrNotImplemented
}

// LookupLanguage performs one exact ID-or-name lookup.
func (c *Catalog) LookupLanguage(value string) (domain.Language, bool) {
	return domain.Language{}, false
}

// LanguageByID looks up one language by its runtime numeric ID.
func (c *Catalog) LanguageByID(id domain.LanguageID) (domain.Language, bool) {
	return domain.Language{}, false
}

// RecordStatus looks up one status by its runtime numeric ID.
func (c *Catalog) RecordStatus(id int) (domain.RecordStatus, bool) {
	return domain.RecordStatus{}, false
}

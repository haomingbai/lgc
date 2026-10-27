// Package configuration owns non-secret configuration discovery and
// validation for the LGC command.
package configuration

import (
	"context"
	"time"

	"github.com/haomingbai/lgc/cmd/internal/arguments"
	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
)

// PluginSelection contains provider IDs selected by configuration.
type PluginSelection struct {
	Captcha       string
	SessionStore  string
	RecordUpdate  string
	ProblemLoader string
}

// Configuration is the validated, non-secret command configuration.
type Configuration struct {
	Profile        string
	BaseURL        string
	RequestTimeout time.Duration
	JudgeTimeout   time.Duration
	PollInterval   time.Duration
	Plugins        PluginSelection
	PluginConfig   map[string]port.ProviderConfig
}

// Reader selects, loads, and validates one configuration file.
type Reader struct{}

// NewReader creates a configuration reader.
func NewReader() *Reader {
	return &Reader{}
}

// Read consumes the argument stage's invocation and returns validated config.
func (*Reader) Read(
	ctx context.Context,
	invocation arguments.Invocation,
) (Configuration, error) {
	return Configuration{}, domain.ErrNotImplemented
}

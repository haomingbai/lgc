// Package arguments owns process-argument interpretation for the LGC command.
package arguments

import "github.com/haomingbai/lgc/lib/domain"

// Command is one supported command-line operation.
type Command string

const (
	CommandLogin   Command = "login"
	CommandSubmit  Command = "submit"
	CommandHelp    Command = "help"
	CommandVersion Command = "version"
)

// OutputMode selects the structured or human-readable command presentation.
type OutputMode string

const (
	OutputText OutputMode = "text"
	OutputJSON OutputMode = "json"
)

// Invocation is the immutable result of argument interpretation.
type Invocation struct {
	Command       Command
	ConfigPath    string
	ConfigPathSet bool
	Profile       string
	Account       string
	ProblemID     string
	Language      string
	SourcePath    string
	SourceStdin   bool
	RequestO2     bool
	OutputMode    OutputMode
}

// Reader validates command grammar without performing I/O.
type Reader struct{}

// NewReader creates an argument reader.
func NewReader() *Reader {
	return &Reader{}
}

// Read converts raw process arguments into one invocation.
func (*Reader) Read(args []string) (Invocation, error) {
	return Invocation{}, domain.ErrNotImplemented
}

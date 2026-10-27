// Package executor runs exactly one prepared LGC command.
package executor

import (
	"context"
	"errors"
	"io"

	"github.com/haomingbai/lgc/cmd/internal/arguments"
	"github.com/haomingbai/lgc/cmd/internal/resources"
	"github.com/haomingbai/lgc/lib/domain"
)

// ExitCode is the stable process exit category.
type ExitCode int

const (
	ExitSuccess            ExitCode = 0
	ExitInternal           ExitCode = 1
	ExitUsage              ExitCode = 2
	ExitAuthentication     ExitCode = 3
	ExitSubmission         ExitCode = 4
	ExitObservationTimeout ExitCode = 5
)

// ErrorInfo is the safe error portion of the structured result.
type ErrorInfo struct {
	Category string `json:"category"`
	Message  string `json:"message"`
}

// Result is the versioned machine-readable command result.
type Result struct {
	SchemaVersion int               `json:"schema_version"`
	Command       arguments.Command `json:"command"`
	OK            bool              `json:"ok"`
	Data          any               `json:"data,omitempty"`
	Error         *ErrorInfo        `json:"error,omitempty"`
}

// Executor owns command workflow execution, output, and exit classification.
type Executor struct {
	output io.Writer
}

// NewExecutor creates one executor with a single structured-output writer.
func NewExecutor(output io.Writer) *Executor {
	return &Executor{output: output}
}

// Execute runs the selected workflow against an already-generated bundle.
func (e *Executor) Execute(
	ctx context.Context,
	resources *resources.ExecutionResources,
) (Result, ExitCode) {
	return Result{
			SchemaVersion: 1,
			Command:       resources.Invocation.Command,
			Error: &ErrorInfo{
				Category: "not_implemented",
				Message:  domain.ErrNotImplemented.Error(),
			},
		},
		ExitInternal
}

// ExitCodeForError maps stable error categories to process exit status.
func ExitCodeForError(err error) ExitCode {
	switch {
	case errors.Is(err, domain.ErrUsage):
		return ExitUsage
	case errors.Is(err, domain.ErrUnauthenticated),
		errors.Is(err, domain.ErrPlugin):
		return ExitAuthentication
	case errors.Is(err, domain.ErrSubmissionRejected),
		errors.Is(err, domain.ErrMalformedUpstream),
		errors.Is(err, domain.ErrInvalidLanguage),
		errors.Is(err, domain.ErrO2Unsupported):
		return ExitSubmission
	case errors.Is(err, domain.ErrJudgeTimeout):
		return ExitObservationTimeout
	default:
		return ExitInternal
	}
}

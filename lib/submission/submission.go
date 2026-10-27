// Package submission owns construction policy, request serialization, and
// submission execution.
package submission

import (
	"context"

	"github.com/haomingbai/lgc/lib/config"
	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
	"github.com/haomingbai/lgc/lib/problem"
	"github.com/haomingbai/lgc/lib/session"
)

// SubmissionInput contains parsed caller values.
type SubmissionInput struct {
	Problem           problem.Problem
	RequestedLanguage string
	Code              string
	RequestO2         bool
}

// Submission contains final values ready for encoding.
type Submission struct {
	Problem    problem.Problem
	LanguageID domain.LanguageID
	Code       string
	EnableO2   int
}

// SubmissionReceipt identifies an accepted submission.
type SubmissionReceipt struct {
	RID string
}

// SubmissionEncoder serializes one final submission.
type SubmissionEncoder interface {
	Encode(submission Submission, captcha *string) ([]byte, error)
}

// SubmissionRoutes contains one coherent submission and challenge route pair.
type SubmissionRoutes struct {
	Submit            string
	SubmissionCaptcha string
}

// SubmissionBuilder applies language, O2, and source construction policy.
type SubmissionBuilder struct {
	catalog *config.Catalog
}

// NewSubmissionBuilder creates a builder over one runtime catalog.
func NewSubmissionBuilder(catalog *config.Catalog) *SubmissionBuilder {
	return &SubmissionBuilder{catalog: catalog}
}

// Build applies all construction policy exactly once.
func (b *SubmissionBuilder) Build(
	input SubmissionInput,
) (Submission, error) {
	return Submission{}, domain.ErrNotImplemented
}

// JSONSubmissionEncoder is the internal seam for the verified JSON request
// shape.
type JSONSubmissionEncoder struct{}

// NewSubmissionEncoder creates the verified submission encoder.
func NewSubmissionEncoder() SubmissionEncoder {
	return &JSONSubmissionEncoder{}
}

// Encode serializes final submission values.
func (*JSONSubmissionEncoder) Encode(
	submission Submission,
	captcha *string,
) ([]byte, error) {
	return nil, domain.ErrNotImplemented
}

// SubmissionExecutor performs the initial POST and verified CAPTCHA retry.
type SubmissionExecutor struct {
	encoder SubmissionEncoder
	captcha port.CaptchaSolver
	routes  SubmissionRoutes
}

// NewSubmissionExecutor creates an executor from already-selected resources.
func NewSubmissionExecutor(
	encoder SubmissionEncoder,
	captcha port.CaptchaSolver,
	routes SubmissionRoutes,
) *SubmissionExecutor {
	return &SubmissionExecutor{
		encoder: encoder,
		captcha: captcha,
		routes:  routes,
	}
}

// Execute submits one already-built submission and returns its RID.
func (e *SubmissionExecutor) Execute(
	ctx context.Context,
	session *session.LuoguSession,
	submission Submission,
) (SubmissionReceipt, error) {
	return SubmissionReceipt{}, domain.ErrNotImplemented
}

var _ SubmissionEncoder = (*JSONSubmissionEncoder)(nil)

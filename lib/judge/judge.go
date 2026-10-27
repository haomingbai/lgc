// Package judge owns judge-record acquisition, reduction, and observation.
package judge

import (
	"context"
	"time"

	"github.com/haomingbai/lgc/lib/config"
	"github.com/haomingbai/lgc/lib/domain"
	"github.com/haomingbai/lgc/lib/port"
	"github.com/haomingbai/lgc/lib/session"
)

// RecordUpdate is an advisory update supplied by an optional provider.
type RecordUpdate = domain.RecordUpdate

// CompileResult is the safe compile portion of a judge record.
type CompileResult struct {
	Success bool
	Message string
	O2      bool
}

// JudgeResult is the safe case-progress portion of a judge record.
type JudgeResult struct {
	FinishedCaseCount int
	Status            int
	Time              int
	Memory            int
	Score             int
	Subtasks          []any
}

// JudgeRecord is one coherent judge-record view.
type JudgeRecord struct {
	RID        string
	Status     domain.RecordStatus
	LanguageID domain.LanguageID
	EnableO2   bool
	Compile    CompileResult
	Judge      JudgeResult
	ReceivedAt time.Time
}

// ObservationPolicy bounds polling and total observation time.
type ObservationPolicy struct {
	InitialPollDelay time.Duration
	MaxPollDelay     time.Duration
	PollTimeout      time.Duration
}

// JudgeObserver acquires records, merges updates, and polls to terminality.
type JudgeObserver struct {
	session      *session.LuoguSession
	catalog      *config.Catalog
	recordRoute  string
	updateSource port.RecordUpdateSource
	policy       ObservationPolicy
}

// NewJudgeObserver creates an observer from one runtime route snapshot.
func NewJudgeObserver(
	session *session.LuoguSession,
	catalog *config.Catalog,
	recordRoute string,
	updateSource port.RecordUpdateSource,
	policy ObservationPolicy,
) *JudgeObserver {
	return &JudgeObserver{
		session:      session,
		catalog:      catalog,
		recordRoute:  recordRoute,
		updateSource: updateSource,
		policy:       policy,
	}
}

// Observe returns a terminal record or the latest record with a timeout error.
func (o *JudgeObserver) Observe(
	ctx context.Context,
	rid string,
) (JudgeRecord, error) {
	return JudgeRecord{}, domain.ErrNotImplemented
}

// RecordReducer merges canonical snapshots and advisory updates.
type RecordReducer struct{}

// NewRecordReducer creates an empty reducer.
func NewRecordReducer() *RecordReducer {
	return &RecordReducer{}
}

// ApplySnapshot applies one canonical record snapshot.
func (r *RecordReducer) ApplySnapshot(record JudgeRecord) error {
	return domain.ErrNotImplemented
}

// ApplyUpdate applies one advisory record update.
func (r *RecordReducer) ApplyUpdate(update RecordUpdate) error {
	return domain.ErrNotImplemented
}

// Snapshot returns one coherent reduced record.
func (r *RecordReducer) Snapshot() JudgeRecord {
	return JudgeRecord{}
}

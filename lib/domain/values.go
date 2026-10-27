package domain

import (
	"encoding/json"
	"time"
)

// LanguageID is the runtime identifier assigned to a Luogu language.
type LanguageID int

// Language describes one runtime language entry.
type Language struct {
	ID       LanguageID
	Name     string
	Family   string
	Disabled bool
	CanO2    bool
	FileExts []string
}

// RecordStatus describes one runtime judge status.
type RecordStatus struct {
	ID        int
	Type      string
	Short     string
	PollAgain bool
}

// Problem is a problem identity with optional, already-loaded content.
type Problem struct {
	ID      string
	Content *string
}

// CaptchaKind identifies why a CAPTCHA challenge was requested.
type CaptchaKind string

const (
	// CaptchaLogin is a challenge issued during authentication.
	CaptchaLogin CaptchaKind = "login"
	// CaptchaSubmission is a challenge issued during source submission.
	CaptchaSubmission CaptchaKind = "submission"
)

// CaptchaChallenge is the bounded input passed to a CAPTCHA solver.
type CaptchaChallenge struct {
	Kind      CaptchaKind
	MediaType string
	Image     []byte
}

// RecordUpdate is an advisory update from an optional update source.
type RecordUpdate struct {
	RID        string
	Fields     map[string]json.RawMessage
	ReceivedAt time.Time
}

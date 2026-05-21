// Package locale provides internationalization support for gformiac.
// It auto-detects the display language from environment variables and allows
// runtime switching between supported languages via Set.
package locale

import (
	"os"
	"strings"
	"sync"
)

// Lang represents a supported language.
type Lang string

// Supported language codes.
const (
	// EN selects English messages.
	EN Lang = "en"
	// JA selects Japanese messages.
	JA Lang = "ja"
)

var (
	mu      sync.RWMutex
	current = EN
)

// M is the active message catalog. Use this to access localized strings.
var M = &En

func init() {
	if detect() == JA {
		Set(JA)
	}
}

// Set switches the active language. It is safe to call concurrently.
func Set(l Lang) {
	mu.Lock()
	defer mu.Unlock()
	current = l
	switch l {
	case JA:
		M = &Ja
	default:
		M = &En
	}
}

// Get returns the current language. It is safe to call concurrently.
func Get() Lang {
	mu.RLock()
	defer mu.RUnlock()
	return current
}

func detect() Lang {
	if v := os.Getenv("GFORMIAC_LANG"); v != "" {
		if strings.HasPrefix(v, "ja") {
			return JA
		}
		return EN
	}
	for _, env := range []string{"LC_ALL", "LC_MESSAGES", "LANG"} {
		if v := os.Getenv(env); strings.HasPrefix(v, "ja") {
			return JA
		}
	}
	return EN
}

// Messages holds all user-facing strings for a given language.
// Each field corresponds to a distinct UI message; format strings follow
// the fmt package conventions (%s, %d, %q, %w).
type Messages struct {
	// CLI descriptions used in cobra command help output.
	CmdLong     string
	PlanShort   string
	ApplyShort  string
	ImportShort string

	// Flag descriptions shown in --help output.
	FlagFile        string
	FlagCredentials string
	FlagToken       string
	FlagState       string
	FlagAutoApprove string
	FlagOutput      string
	FlagLang        string

	// Plan / Diff output messages.
	NoChanges     string
	NoChangesLong string
	FormInfo      string
	DiffSummary   string // "+%d ~%d -%d"

	// New form creation summary messages.
	NewFormHeader string
	CreateSummary string // "%d items to create"

	// Apply command progress and result messages.
	Creating     string
	Applied      string
	Cancelled    string
	ConfirmApply string
	FormIDLabel  string
	URLLabel     string
	StateLabel   string

	// Import command result messages.
	Imported      string
	SpecFileLabel string
	ItemCount     string

	// OAuth2 browser authentication messages.
	AuthBrowser     string
	AuthSuccessHTML string
	AuthFailedHTML  string
	AuthCodeMissing string
	AuthTimeout     string
	AuthOpenURL     string

	// Validation error messages (format strings with item prefix %s).
	ValErrors      string
	ValTitle       string
	ValItems       string
	ValItemTitle   string
	ValItemType    string
	ValTypeUnknown string
	ValChoiceReq   string
	ValChoiceType  string
	ValChoiceOpts  string
	ValChoiceEmpty string
	ValScaleReq    string
	ValScaleLow    string
	ValScaleHigh   string
	ValScaleRange  string

	// Error message fragments wrapped via fmt.Errorf.
	ErrReadFile    string
	ErrParseYAML   string
	ErrTitleReq    string
	ErrCredentials string
	ErrCredsParse  string
	ErrAuth        string
	ErrTokenSave   string
	ErrFormsInit   string
	ErrFormCreate  string
	ErrFormSetup   string
	ErrFormGet     string
	ErrFormUpdate  string
	ErrLocalServer string
	ErrStateSave   string
	ErrStateRead   string
	ErrSpecSave    string
	ErrCredsHint   string
}

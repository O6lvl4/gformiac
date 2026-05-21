package locale

import (
	"os"
	"strings"
	"sync"
)

// Lang represents a supported language.
type Lang string

const (
	EN Lang = "en"
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

// Set switches the active language.
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

// Get returns the current language.
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
type Messages struct {
	// CLI
	CmdLong     string
	PlanShort   string
	ApplyShort  string
	ImportShort string

	// Flags
	FlagFile        string
	FlagCredentials string
	FlagToken       string
	FlagState       string
	FlagAutoApprove string
	FlagOutput      string
	FlagLang        string

	// Plan / Diff
	NoChanges     string
	NoChangesLong string
	FormInfo      string
	DiffSummary   string // "+%d ~%d -%d"

	// New form
	NewFormHeader string
	CreateSummary string // "%d items to create"

	// Apply
	Creating     string
	Applied      string
	Cancelled    string
	ConfirmApply string
	FormIDLabel  string
	URLLabel     string
	StateLabel   string

	// Import
	Imported      string
	SpecFileLabel string
	ItemCount     string

	// Auth
	AuthBrowser     string
	AuthSuccessHTML string
	AuthFailedHTML  string
	AuthCodeMissing string
	AuthTimeout     string
	AuthOpenURL     string

	// Validation (format strings with prefix %s)
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

	// Errors (format strings for fmt.Errorf)
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

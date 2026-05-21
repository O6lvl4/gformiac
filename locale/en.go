package locale

// En is the English message catalog.
var En = Messages{
	// CLI
	CmdLong:     "Declaratively manage Google Forms from YAML definitions",
	PlanShort:   "Preview changes (dry-run)",
	ApplyShort:  "Apply form definition",
	ImportShort: "Import an existing Google Form to YAML",

	// Flags
	FlagFile:        "Form definition file",
	FlagCredentials: "OAuth2 credentials file",
	FlagToken:       "OAuth token file",
	FlagState:       "State file",
	FlagAutoApprove: "Skip confirmation prompt",
	FlagOutput:      "Output file (defaults to --file value)",
	FlagLang:        "Language (en, ja)",

	// Plan / Diff
	NoChanges:     "No changes",
	NoChangesLong: "No changes — form is up to date",
	FormInfo:      "Form info:",
	DiffSummary:   "Total: +%d ~%d -%d",

	// New form
	NewFormHeader: "Creating new form:",
	CreateSummary: "Total: %d items to create",

	// Apply
	Creating:     "Creating form...",
	Applied:      "Applied!",
	Cancelled:    "Cancelled",
	ConfirmApply: "Apply?",
	FormIDLabel:  "  Form ID:    %s",
	URLLabel:     "  URL:        %s",
	StateLabel:   "  State file: %s",

	// Import
	Imported:      "Import completed!",
	SpecFileLabel: "  Spec file:  %s",
	ItemCount:     "  Items:      %d",

	// Auth
	AuthBrowser:     "Authenticating via browser...",
	AuthSuccessHTML: "<h2>Authentication successful!</h2><p>You can close this tab and return to the terminal.</p>",
	AuthFailedHTML:  "<h2>Authentication failed</h2><p>Auth code was not received.</p>",
	AuthCodeMissing: "auth code not found in callback",
	AuthTimeout:     "authentication timeout (2 min)",
	AuthOpenURL:     "Open this URL in your browser:\n%s",

	// Validation
	ValErrors:      "Validation errors (%d):\n  %s",
	ValTitle:       "title is required",
	ValItems:       "items must have at least one entry",
	ValItemTitle:   "%s: title is required",
	ValItemType:    "%s: type is required",
	ValTypeUnknown: "%s: unknown type %q (valid: %s)",
	ValChoiceReq:   "%s: choice field is required for type=choice",
	ValChoiceType:  "%s: unknown choice.type %q (valid: %s)",
	ValChoiceOpts:  "%s: choice.options requires at least one entry",
	ValChoiceEmpty: "%s: choice.options[%d] must not be empty",
	ValScaleReq:    "%s: scale field is required for type=scale",
	ValScaleLow:    "%s: scale.low must be 0 or 1 (got %d)",
	ValScaleHigh:   "%s: scale.high must be 2-10 (got %d)",
	ValScaleRange:  "%s: scale.low (%d) must be less than scale.high (%d)",

	// Errors (descriptive text only, no %w — wrapped via errWrap)
	ErrReadFile:    "reading file",
	ErrParseYAML:   "parsing YAML",
	ErrTitleReq:    "title is required",
	ErrCredentials: "reading credentials",
	ErrCredsParse:  "parsing credentials",
	ErrAuth:        "authentication failed",
	ErrTokenSave:   "saving token",
	ErrFormsInit:   "initializing Forms API",
	ErrFormCreate:  "creating form",
	ErrFormSetup:   "configuring form",
	ErrFormGet:     "fetching form",
	ErrFormUpdate:  "updating form",
	ErrLocalServer: "starting local server",
	ErrStateSave:   "saving state",
	ErrStateRead:   "reading state",
	ErrSpecSave:    "saving spec",
	ErrCredsHint:   "\n\nHint: run `gcloud auth application-default login --scopes=%s,https://www.googleapis.com/auth/cloud-platform` or place an OAuth2 credentials.json file",
}

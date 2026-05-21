package locale

// Ja is the Japanese message catalog.
var Ja = Messages{
	// CLI
	CmdLong:     "YAML定義からGoogle Formsを宣言的に管理するIaCツール",
	PlanShort:   "変更のプレビュー（dry-run）",
	ApplyShort:  "フォーム定義を適用",
	ImportShort: "既存のGoogle FormをYAML定義にインポート",

	// Flags
	FlagFile:        "フォーム定義ファイル",
	FlagCredentials: "OAuth2認証情報ファイル",
	FlagToken:       "OAuthトークンファイル",
	FlagState:       "状態ファイル",
	FlagAutoApprove: "確認をスキップ",
	FlagOutput:      "出力ファイル（未指定時は --file の値）",
	FlagLang:        "言語 (en, ja)",

	// Plan / Diff
	NoChanges:     "変更なし",
	NoChangesLong: "変更なし — フォームは最新です",
	FormInfo:      "フォーム情報:",
	DiffSummary:   "合計: +%d ~%d -%d",

	// New form
	NewFormHeader: "新規フォーム作成:",
	CreateSummary: "合計: %d項目を作成",

	// Apply
	Creating:     "フォーム作成中...",
	Applied:      "適用完了!",
	Cancelled:    "キャンセルしました",
	ConfirmApply: "適用しますか？",
	FormIDLabel:  "  フォームID:  %s",
	URLLabel:     "  回答URL:     %s",
	StateLabel:   "  状態ファイル: %s",

	// Import
	Imported:      "インポート完了!",
	SpecFileLabel: "  定義ファイル: %s",
	ItemCount:     "  項目数: %d",

	// Auth
	AuthBrowser:     "ブラウザで認証してください...",
	AuthSuccessHTML: "<h2>認証成功!</h2><p>このタブを閉じてターミナルに戻ってください。</p>",
	AuthFailedHTML:  "<h2>認証失敗</h2><p>認証コードが取得できませんでした。</p>",
	AuthCodeMissing: "認証コードが見つかりません",
	AuthTimeout:     "認証タイムアウト（2分）",
	AuthOpenURL:     "以下のURLをブラウザで開いてください:\n%s",

	// Validation
	ValErrors:      "バリデーションエラー (%d件):\n  %s",
	ValTitle:       "title は必須です",
	ValItems:       "items は1つ以上必要です",
	ValItemTitle:   "%s: title は必須です",
	ValItemType:    "%s: type は必須です",
	ValTypeUnknown: "%s: 不明な type %q (有効値: %s)",
	ValChoiceReq:   "%s: type=choice には choice フィールドが必須です",
	ValChoiceType:  "%s: 不明な choice.type %q (有効値: %s)",
	ValChoiceOpts:  "%s: choice.options は1つ以上必要です",
	ValChoiceEmpty: "%s: choice.options[%d] は空にできません",
	ValScaleReq:    "%s: type=scale には scale フィールドが必須です",
	ValScaleLow:    "%s: scale.low は 0 または 1 のみ (got %d)",
	ValScaleHigh:   "%s: scale.high は 2〜10 の範囲 (got %d)",
	ValScaleRange:  "%s: scale.low (%d) は scale.high (%d) より小さくなければなりません",

	// Errors
	ErrReadFile:    "ファイル読み込み失敗",
	ErrParseYAML:   "YAML解析失敗",
	ErrTitleReq:    "titleは必須です",
	ErrCredentials: "credentials読み込み失敗",
	ErrCredsParse:  "credentials解析失敗",
	ErrAuth:        "認証失敗",
	ErrTokenSave:   "トークン保存失敗",
	ErrFormsInit:   "Forms API初期化失敗",
	ErrFormCreate:  "フォーム作成失敗",
	ErrFormSetup:   "フォーム設定失敗",
	ErrFormGet:     "フォーム取得失敗",
	ErrFormUpdate:  "フォーム更新失敗",
	ErrLocalServer: "ローカルサーバー起動失敗",
	ErrStateSave:   "状態保存失敗",
	ErrStateRead:   "状態ファイル読み込み失敗",
	ErrSpecSave:    "定義ファイル保存失敗",
	ErrCredsHint:   "\n\nヒント: gcloud auth application-default login --scopes=%s,https://www.googleapis.com/auth/cloud-platform を実行するか、OAuth2 credentials.json を配置してください",
}

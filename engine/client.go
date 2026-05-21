package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	forms "google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

type Client struct {
	svc *forms.Service
}

func NewClient(ctx context.Context, credentialsFile, tokenFile string) (*Client, error) {
	// Service account via GOOGLE_APPLICATION_CREDENTIALS
	if os.Getenv("GOOGLE_APPLICATION_CREDENTIALS") != "" {
		svc, err := forms.NewService(ctx,
			option.WithScopes("https://www.googleapis.com/auth/forms.body"),
		)
		if err != nil {
			return nil, fmt.Errorf("サービスアカウント認証失敗: %w", err)
		}
		return &Client{svc: svc}, nil
	}

	// OAuth2 flow
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		return nil, fmt.Errorf("credentials読み込み失敗 (%s): %w", credentialsFile, err)
	}

	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/forms.body")
	if err != nil {
		return nil, fmt.Errorf("credentials解析失敗: %w", err)
	}

	tok, err := tokenFromFile(tokenFile)
	if err != nil {
		tok, err = tokenFromWeb(config)
		if err != nil {
			return nil, fmt.Errorf("認証失敗: %w", err)
		}
		saveToken(tokenFile, tok)
	}

	svc, err := forms.NewService(ctx, option.WithHTTPClient(config.Client(ctx, tok)))
	if err != nil {
		return nil, fmt.Errorf("Forms API初期化失敗: %w", err)
	}

	return &Client{svc: svc}, nil
}

// Plan fetches the remote form and computes the diff against the local spec.
func (c *Client) Plan(ctx context.Context, formID string, spec *FormSpec, state *State) (*DiffResult, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("フォーム取得失敗: %w", err)
	}
	remote := formToSpec(form)
	return Diff(spec, remote, state), nil
}

// CreateForm creates a new Google Form from the spec.
func (c *Client) CreateForm(ctx context.Context, spec *FormSpec) (*State, error) {
	form, err := c.svc.Forms.Create(&forms.Form{
		Info: &forms.Info{Title: spec.Title},
	}).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("フォーム作成失敗: %w", err)
	}
	fmt.Printf("  フォーム作成完了: %s\n", form.FormId)

	requests := specToCreateRequests(spec)
	if len(requests) > 0 {
		_, err = c.svc.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{
			Requests: requests,
		}).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("フォーム設定失敗: %w", err)
		}
	}

	form, err = c.svc.Forms.Get(form.FormId).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("フォーム取得失敗: %w", err)
	}

	return buildState(form), nil
}

// UpdateForm reconciles the remote form to match the local spec.
// Strategy: delete all existing items, then recreate from spec.
func (c *Client) UpdateForm(ctx context.Context, formID string, spec *FormSpec) (*State, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("フォーム取得失敗: %w", err)
	}

	requests := specToUpdateRequests(spec, form)
	if len(requests) > 0 {
		_, err = c.svc.Forms.BatchUpdate(formID, &forms.BatchUpdateFormRequest{
			Requests: requests,
		}).Context(ctx).Do()
		if err != nil {
			return nil, fmt.Errorf("フォーム更新失敗: %w", err)
		}
	}

	form, err = c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("フォーム取得失敗: %w", err)
	}

	return buildState(form), nil
}

// ImportForm fetches a remote form and returns it as a FormSpec.
func (c *Client) ImportForm(ctx context.Context, formID string) (*FormSpec, *State, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("フォーム取得失敗: %w", err)
	}
	return formToSpec(form), buildState(form), nil
}

func tokenFromWeb(config *oauth2.Config) (*oauth2.Token, error) {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("ブラウザで以下のURLを開いて認証してください:\n%s\n\n認証コードを入力: ", authURL)
	var code string
	if _, err := fmt.Scan(&code); err != nil {
		return nil, err
	}
	return config.Exchange(context.Background(), code)
}

func tokenFromFile(path string) (*oauth2.Token, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var tok oauth2.Token
	return &tok, json.NewDecoder(f).Decode(&tok)
}

func saveToken(path string, tok *oauth2.Token) {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(tok)
}

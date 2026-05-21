package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/O6lvl4/gformiac/locale"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	forms "google.golang.org/api/forms/v1"
	"google.golang.org/api/option"
)

const formsScope = "https://www.googleapis.com/auth/forms.body"

// Client wraps the Google Forms API service.
type Client struct {
	svc *forms.Service
}

// NewClient creates an authenticated Forms API client.
func NewClient(ctx context.Context, credentialsFile, tokenFile string) (*Client, error) {
	if client, err := tryDefaultCredentials(ctx); err == nil {
		return client, nil
	}
	return newOAuth2Client(ctx, credentialsFile, tokenFile)
}

func tryDefaultCredentials(ctx context.Context) (*Client, error) {
	ts, err := google.DefaultTokenSource(ctx, formsScope)
	if err != nil {
		return nil, err
	}
	svc, err := forms.NewService(ctx, option.WithTokenSource(ts))
	if err != nil {
		return nil, err
	}
	return &Client{svc: svc}, nil
}

func newOAuth2Client(ctx context.Context, credentialsFile, tokenFile string) (*Client, error) {
	b, err := os.ReadFile(credentialsFile)
	if err != nil {
		hint := fmt.Sprintf(locale.M.ErrCredsHint, formsScope)
		return nil, fmt.Errorf("%s (%s): %w%s", locale.M.ErrCredentials, credentialsFile, err, hint)
	}

	config, err := google.ConfigFromJSON(b, formsScope)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrCredsParse, err)
	}

	tok, err := resolveToken(ctx, config, tokenFile)
	if err != nil {
		return nil, err
	}

	svc, err := forms.NewService(ctx, option.WithHTTPClient(config.Client(ctx, tok)))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrFormsInit, err)
	}
	return &Client{svc: svc}, nil
}

func resolveToken(ctx context.Context, config *oauth2.Config, tokenFile string) (*oauth2.Token, error) {
	if tok, err := tokenFromFile(tokenFile); err == nil {
		return tok, nil
	}
	tok, err := tokenFromWeb(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrAuth, err)
	}
	if err := saveToken(tokenFile, tok); err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrTokenSave, err)
	}
	return tok, nil
}

// Plan fetches the remote form and computes the diff against the local spec.
func (c *Client) Plan(ctx context.Context, formID string, spec *FormSpec, state *State) (*DiffResult, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrFormGet, err)
	}
	return Diff(spec, formToSpec(form), state), nil
}

// CreateForm creates a new Google Form from the spec and returns the resulting state.
func (c *Client) CreateForm(ctx context.Context, spec *FormSpec) (*State, error) {
	form, err := c.svc.Forms.Create(&forms.Form{
		Info: &forms.Info{Title: spec.Title},
	}).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrFormCreate, err)
	}

	if requests := specToCreateRequests(spec); len(requests) > 0 {
		if _, err := c.svc.Forms.BatchUpdate(form.FormId, &forms.BatchUpdateFormRequest{
			Requests: requests,
		}).Context(ctx).Do(); err != nil {
			return nil, fmt.Errorf("%s: %w", locale.M.ErrFormSetup, err)
		}
	}

	return c.fetchState(ctx, form.FormId)
}

// UpdateForm reconciles the remote form to match the local spec.
func (c *Client) UpdateForm(ctx context.Context, formID string, spec *FormSpec) (*State, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrFormGet, err)
	}

	if requests := specToUpdateRequests(spec, form); len(requests) > 0 {
		if _, err := c.svc.Forms.BatchUpdate(formID, &forms.BatchUpdateFormRequest{
			Requests: requests,
		}).Context(ctx).Do(); err != nil {
			return nil, fmt.Errorf("%s: %w", locale.M.ErrFormUpdate, err)
		}
	}

	return c.fetchState(ctx, formID)
}

// ImportForm fetches a remote form and converts it to a FormSpec and State.
func (c *Client) ImportForm(ctx context.Context, formID string) (*FormSpec, *State, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, nil, fmt.Errorf("%s: %w", locale.M.ErrFormGet, err)
	}
	return formToSpec(form), buildState(form), nil
}

func (c *Client) fetchState(ctx context.Context, formID string) (*State, error) {
	form, err := c.svc.Forms.Get(formID).Context(ctx).Do()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrFormGet, err)
	}
	return buildState(form), nil
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) (*oauth2.Token, error) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", locale.M.ErrLocalServer, err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	config.RedirectURL = fmt.Sprintf("http://localhost:%d", port)

	codeCh := make(chan string, 1)
	errCh := make(chan error, 1)
	srv := &http.Server{Handler: callbackHandler(codeCh, errCh)}
	go srv.Serve(listener)
	defer srv.Close()

	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Println(locale.M.AuthBrowser)
	openBrowser(authURL)

	return waitForToken(ctx, config, codeCh, errCh)
}

func callbackHandler(codeCh chan<- string, errCh chan<- error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if code == "" {
			fmt.Fprintln(w, locale.M.AuthFailedHTML)
			errCh <- fmt.Errorf("%s", locale.M.AuthCodeMissing)
			return
		}
		fmt.Fprintln(w, locale.M.AuthSuccessHTML)
		codeCh <- code
	})
}

func waitForToken(ctx context.Context, config *oauth2.Config, codeCh <-chan string, errCh <-chan error) (*oauth2.Token, error) {
	select {
	case code := <-codeCh:
		return config.Exchange(ctx, code)
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(2 * time.Minute):
		return nil, fmt.Errorf("%s", locale.M.AuthTimeout)
	}
}

func openBrowser(url string) {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "darwin":
		cmd = "open"
	case "linux":
		cmd = "xdg-open"
	case "windows":
		cmd, args = "rundll32", []string{"url.dll,FileProtocolHandler"}
	default:
		fmt.Printf(locale.M.AuthOpenURL+"\n", url)
		return
	}
	exec.Command(cmd, append(args, url)...).Start()
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

func saveToken(path string, tok *oauth2.Token) error {
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return json.NewEncoder(f).Encode(tok)
}

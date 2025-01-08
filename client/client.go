package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eldius/onedrive-client/internal/configs"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

const (
	graphApiEndpoint = "https://graph.microsoft.com/v1.0/me/"
)

type Client interface {
	Authenticate() (*TokenData, error)
	AuthenticatedUser() (*CurrentUser, error)
	GetAppDriveInfo() error
}

type client struct {
	c     *http.Client
	creds struct {
		id          string
		secret      string
		token       *TokenData
		scopes      []string
		redirectURL string
	}
}

func (c *client) getScopes() []string {
	if c.creds.scopes == nil || len(c.creds.scopes) == 0 {
		return configs.DefaultAuthScopes
	}
	return c.creds.scopes
}

func (c *client) getRedirectURL() string {
	if c.creds.redirectURL == "" {
		return configs.DefaultRedirectURL
	}
	return c.creds.redirectURL
}

func (c *client) getAuthPath() (string, error) {
	if c.creds.redirectURL == "" {
		return "GET /authentication", nil
	}
	u, err := url.Parse(c.creds.redirectURL)
	if err != nil {
		return "", err
	}
	return u.Path, nil
}

// Option is used to configure
// client properties
type Option func(*client)

// WithSecretID sets up the client secret id
func WithSecretID(id string) Option {
	return func(c *client) {
		if id == "" {
			return
		}
		c.creds.id = id
	}
}

// WithSecretKey sets up the client secret key
func WithSecretKey(secret string) Option {
	return func(c *client) {
		if secret == "" {
			return
		}
		c.creds.secret = secret
	}
}

// WithAuthenticationTokenData sets up the authentication token
// (when you already authenticated to OneDrive API)
func WithAuthenticationTokenData(token *TokenData) Option {
	return func(c *client) {
		if token == nil {
			return
		}
		c.creds.token = token
	}
}

// WithHttpClient is to define a custom http.Client
func WithHttpClient(hc *http.Client) Option {
	return func(c *client) {
		if hc != nil {
			c.c = hc
		}
	}
}

// WithScopes defines scopes to be used
// for authentication
func WithScopes(scopes ...string) Option {
	return func(c *client) {
		if scopes == nil || len(scopes) == 0 {
			return
		}
		c.creds.scopes = scopes
	}
}

func WithRedirectURL(redirectURL string) Option {
	return func(c *client) {
		if redirectURL == "" {
			return
		}
		c.creds.redirectURL = redirectURL
	}
}

func New(opts ...Option) Client {
	c := &client{}
	for _, opt := range opts {
		opt(c)
	}

	if c.c == nil {
		c.c = &http.Client{}
	}

	return c
}

func (c *client) Authenticate() (*TokenData, error) {
	td, err := newAuthenticator(c).Authenticate()
	if err != nil {
		return nil, err
	}
	c.creds.token = td.TokenData
	return td.TokenData, err
}

func (c *client) AuthenticatedUser() (*CurrentUser, error) {
	req, err := http.NewRequest(http.MethodGet, graphApiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	var resp CurrentUser
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return &resp, nil
}

func (c *client) GetAppDriveInfo() error {
	req, err := http.NewRequest(http.MethodGet, graphApiEndpoint+"/drive/special/approot", nil)
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	var res AppFolderInfo
	if err := c.do(req, &res); err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	return nil
}

func (c *client) addAuthHeaders(r *http.Request) error {
	if c.creds.token == nil {
		return errors.New("no token")
	}
	r.Header.Set("Authorization", "Bearer "+c.creds.token.AccessToken)
	return nil
}

func (c *client) do(req *http.Request, resp APIResponse) error {
	if err := c.addAuthHeaders(req); err != nil {
		return fmt.Errorf("add auth headers: %w", err)
	}
	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	l := slog.With(
		slog.String("url", req.URL.String()),
		slog.String("method", req.Method),
	)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	l.With(slog.String("body", string(b))).Debug("response body")

	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&resp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	resp.SetRawBody(string(b))
	resp.SetStatusCode(res.StatusCode)

	return nil
}

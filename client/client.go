package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eldius/onedrive-client/client/types"
	"github.com/eldius/onedrive-client/internal/configs"
	"io"
	"net/http"
	"net/url"
)

const (
	graphApiEndpoint = "https://graph.microsoft.com/v1.0/me/"
)

type Client interface {
	Authenticate(ctx context.Context) (*types.TokenData, error)
	AuthenticatedUser(ctx context.Context) (*types.CurrentUser, error)
	GetAppDriveInfo(ctx context.Context) (*types.AppFolderInfo, error)
	CreateFolder(ctx context.Context, dirName, parentID, driveID string) (*types.CreateFile, error)
}

type client struct {
	c     *http.Client
	creds struct {
		id          string
		secret      string
		token       *types.TokenData
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
func WithAuthenticationTokenData(token *types.TokenData) Option {
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

func (c *client) Authenticate(ctx context.Context) (*types.TokenData, error) {
	td, err := newAuthenticator(c).Authenticate(ctx)
	if err != nil {
		return nil, err
	}
	c.creds.token = td.TokenData
	return td.TokenData, err
}

func (c *client) AuthenticatedUser(ctx context.Context) (*types.CurrentUser, error) {
	req, err := http.NewRequest(http.MethodGet, graphApiEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)

	var resp types.CurrentUser
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}

	return &resp, nil
}

func (c *client) GetAppDriveInfo(ctx context.Context) (*types.AppFolderInfo, error) {
	req, err := http.NewRequest(http.MethodGet, graphApiEndpoint+"/drive/special/approot", nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	var res types.AppFolderInfo
	if err := c.do(req, &res); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &res, nil
}

func (c *client) CreateFolder(ctx context.Context, dirName, parentID, driveID string) (*types.CreateFile, error) {
	b, err := json.Marshal(map[string]interface{}{
		"name":                              dirName,
		"file":                              make(map[string]interface{}),
		"@microsoft.graph.conflictBehavior": "rename",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, graphApiEndpoint+fmt.Sprintf("/drives/%s/items/%s/children", driveID, parentID), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	var res types.CreateFile
	if err := c.do(req, &res); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &res, nil
}

func (c *client) addAuthHeaders(r *http.Request) error {
	if c.creds.token == nil {
		return errors.New("no token")
	}
	r.Header.Set("Authorization", "Bearer "+c.creds.token.AccessToken)
	return nil
}

func (c *client) do(req *http.Request, resp types.APIResponse) error {
	if err := c.addAuthHeaders(req); err != nil {
		return fmt.Errorf("add auth headers: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	debugResponse(res)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&resp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	resp.SetRawBody(string(b))
	resp.SetStatusCode(res.StatusCode)

	return nil
}

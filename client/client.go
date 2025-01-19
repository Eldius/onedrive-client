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
	"strings"
)

const (
	graphApiEndpoint = "https://graph.microsoft.com/v1.0/me/"
)

type Client interface {
	Authenticate(ctx context.Context) (*types.TokenData, error)
	AuthenticatedUser(ctx context.Context) (*types.CurrentUser, error)
	GetAppDriveInfo(ctx context.Context) (*types.AppFolderInfo, error)
	ListFiles(ctx context.Context, driveID, itemID string) (*types.ListFiles, error)

	CreateFolder(
		ctx context.Context,
		dirName,
		parentID,
		driveID string,
	) (*types.CreateFile, error)

	UploadFile(
		ctx context.Context,
		dirName,
		parentID,
		driveID string,
	) (*types.CreateFile, error)

	CreateUploadSession(
		ctx context.Context,
		driveID,
		itemID,
		inputFile,
		outputFile string,
	) (*types.ListFiles, error)
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
	if len(c.creds.scopes) == 0 {
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
		if len(scopes) == 0 {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var resp types.CurrentUser
	if err := c.doWithRefreshTokenIfUnauthorized(ctx, req, &resp, true, true); err != nil {
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
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var res types.AppFolderInfo
	if err := c.doWithRefreshTokenIfUnauthorized(ctx, req, &res, true, true); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &res, nil
}

func (c *client) CreateFolder(ctx context.Context, dirName, parentID, driveID string) (*types.CreateFile, error) {
	b, err := json.Marshal(folderPayload{
		Name:              dirName,
		Folder:            folder{},
		ConflictBehaviour: "rename",
	})
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, graphApiEndpoint+fmt.Sprintf("/drives/%s/items/%s/children", driveID, parentID), bytes.NewBuffer(b))
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var res types.CreateFile
	if err := c.doWithRefreshTokenIfUnauthorized(ctx, req, &res, true, true); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &res, nil
}

func (c *client) ListFiles(ctx context.Context, driveID, itemID string) (*types.ListFiles, error) {
	req, err := http.NewRequest(http.MethodGet, graphApiEndpoint+fmt.Sprintf("/drives/%s/items/%s/children", driveID, itemID), nil)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var resp types.ListFiles
	if err := c.doWithRefreshTokenIfUnauthorized(ctx, req, &resp, true, true); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &resp, nil
}

func (c *client) CreateUploadSession(ctx context.Context, driveID, itemID, inputFile, outputFile string) (*types.ListFiles, error) {
	b, err := json.Marshal(createUploadSession{
		MicrosoftGraphConflictBehavior: "rename",
		Description:                    "",
		FileSystemInfo:                 createUploadSessionFileSystemInfo{},
		Name:                           outputFile,
	})
	if err != nil {
		return nil, fmt.Errorf("marshal json: %w", err)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		graphApiEndpoint+fmt.Sprintf("/drives/%s/items/%s/createUploadSession", driveID, itemID),
		bytes.NewBuffer(b),
	)
	if err != nil {
		return nil, fmt.Errorf("new request: %w", err)
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	var resp types.ListFiles
	if err := c.doWithRefreshTokenIfUnauthorized(ctx, req, &resp, true, true); err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	return &resp, nil
}

func (c *client) addAuthHeaders(r *http.Request) error {
	if c.creds.token == nil {
		return errors.New("no token")
	}
	r.Header.Set("Authorization", "Bearer "+c.creds.token.AccessToken)
	return nil
}

func (c *client) do(ctx context.Context, req *http.Request, resp types.APIResponse, authenticated bool) error {
	return c.doWithRefreshTokenIfUnauthorized(ctx, req, resp, authenticated, false)
}

func (c *client) doWithRefreshTokenIfUnauthorized(ctx context.Context, req *http.Request, resp types.APIResponse, authenticated, refresh bool) error {
	reqB := make([]byte, 0)
	if req.Body != nil {
		var err error
		reqB, err = io.ReadAll(req.Body)
		if err != nil {
			return fmt.Errorf("read body: %w", err)
		}
	}

	req.Body = io.NopCloser(bytes.NewBuffer(reqB))

	if authenticated {
		if err := c.addAuthHeaders(req); err != nil {
			return fmt.Errorf("add auth headers: %w", err)
		}
	}

	res, err := c.c.Do(req)
	if err != nil {
		return fmt.Errorf("do request: %w", err)
	}
	defer func() {
		_ = res.Body.Close()
	}()

	debugResponse(ctx, res, reqB)

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if err := json.NewDecoder(bytes.NewReader(b)).Decode(&resp); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}

	resp.SetRawBody(string(b))
	resp.SetStatusCode(res.StatusCode)

	if res.StatusCode == http.StatusUnauthorized && refresh {
		if err := c.refreshToken(ctx); err != nil {
			return fmt.Errorf("refresh token: %w", err)
		}
		return c.doWithRefreshTokenIfUnauthorized(ctx, req, resp, authenticated, false)
	}
	if res.StatusCode/100 != 2 {
		return fmt.Errorf("unexpected status code: %d (%d)", res.StatusCode, res.StatusCode/100)
	}

	return nil
}

func (c *client) refreshToken(ctx context.Context) error {
	v := url.Values{}
	v.Set("client_id", configs.GetSecretID())
	v.Set("scope", strings.Join(c.getScopes(), " "))
	v.Set("refresh_token", c.creds.token.RefreshToken)
	v.Set("redirect_uri", c.getRedirectURL())
	v.Set("grant_type", "refresh_token")
	req, err := http.NewRequest(http.MethodPost, tokenEndpoint, strings.NewReader(v.Encode()))
	if err != nil {
		return fmt.Errorf("new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	if err := c.do(ctx, req, c.creds.token, false); err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	return nil
}

func (c *client) UploadFile(ctx context.Context, dirName, parentID, driveID string) (*types.CreateFile, error) {
	return nil, nil
}

type folderPayload struct {
	Name              string `json:"name"`
	Folder            folder `json:"folder"`
	ConflictBehaviour string `json:"@microsoft.graph.conflictBehavior"`
}

type folder struct{}

type createUploadSession struct {
	MicrosoftGraphConflictBehavior string                            `json:"@microsoft.graph.conflictBehavior"`
	Description                    string                            `json:"description"`
	FileSystemInfo                 createUploadSessionFileSystemInfo `json:"fileSystemInfo"`
	Name                           string                            `json:"name"`
}

type createUploadSessionFileSystemInfo struct {
	OdataType string `json:"@odata.type"`
}

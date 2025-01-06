package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/eldius/onedrive-client/internal/configs"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client interface {
	Authenticate() (*TokenData, error)
	AuthenticatedUser() error
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

// WithAuthenticationToken sets up the authentication token
// (when you already authenticated to OneDrive API)
func WithAuthenticationToken(token *TokenData) Option {
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

func (c *client) AuthenticatedUser() error {
	if c.creds.token == nil {
		return errors.New("no token")
	}
	return nil
}

func (c *client) authListener() error {
	mux := http.NewServeMux()
	path, err := c.getAuthPath()
	if err != nil {
		return err
	}
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("path", r.URL.Path),
			"query", r.URL.Query(),
		).Debug("auth handler")
	})

	return http.ListenAndServe(":9999", mux)
}

type authenticator struct {
	c     *client
	state int
}

func newAuthenticator(c *client) *authenticator {
	return &authenticator{
		c:     c,
		state: rand.Int(),
	}
}

func (a *authenticator) authListener() (*authData, error) {
	var response *authData
	var s *http.Server
	mux := http.NewServeMux()
	path, err := a.c.getAuthPath()
	if err != nil {
		return nil, err
	}
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		slog.With(
			slog.String("path", r.URL.Path),
			"query", r.URL.Query(),
		).Debug("auth handler")
		q := r.URL.Query()
		response = &authData{
			Code:             q.Get("code"),
			SentState:        strconv.Itoa(a.state),
			ReceivedState:    q.Get("state"),
			Error:            q.Get("error"),
			ErrorDescription: q.Get("error_description"),
		}
		if response.ReceivedState != strconv.Itoa(a.state) {
			w.WriteHeader(http.StatusUnauthorized)
			renderAuthPage(w, *response)
			return
		}
		response, err = a.generateToken(response)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		renderAuthPage(w, *response)
		go func() {
			time.Sleep(5 * time.Second)
			_ = s.Shutdown(context.Background())
		}()
	})

	//mux.HandleFunc("/auth/shutdown", func(w http.ResponseWriter, r *http.Request) {
	//	slog.With(
	//		slog.String("path", r.URL.Path),
	//		"query", r.URL.Query(),
	//	).Debug("shutting down")
	//	if s != nil {
	//		go func() {
	//			if s != nil {
	//				_ = s.Shutdown(context.Background())
	//			}
	//		}()
	//	}
	//})

	s = &http.Server{
		Addr:    ":9999",
		Handler: mux,
	}

	s.RegisterOnShutdown(func() {
		slog.Info("shutting down...")
		time.Sleep(5 * time.Second)
	})

	_ = s.ListenAndServe()
	return response, nil
}

func (a *authenticator) Authenticate() (*authData, error) {
	u, err := url.Parse("https://login.microsoftonline.com/common/oauth2/v2.0/authorize")
	if err != nil {
		return nil, fmt.Errorf("authenticate: parse login.microsoftonline.com url: %w", err)
	}

	slog.With(
		slog.String("url", u.String()),
		slog.String("id", a.c.creds.id),
		slog.String("secret", a.c.creds.secret),
	).Debug("authenticate")

	q := u.Query()
	q.Add("client_id", a.c.creds.id)
	q.Add("client_secret", a.c.creds.secret)
	q.Add("response_type", "code")
	q.Add("redirect_uri", a.c.getRedirectURL())
	q.Add("response_mode", "query")
	q.Add("scope", strings.Join(a.c.getScopes(), " "))
	q.Add("state", strconv.Itoa(a.state))

	u.RawQuery = q.Encode()

	fmt.Printf("Please, authenticate here: %s\n\n\n\n", u.String())

	return a.authListener()
}

func (a *authenticator) generateToken(d *authData) (*authData, error) {
	v := url.Values{}
	v.Set("client_id", a.c.creds.id)
	v.Set("scope", strings.Join(a.c.getScopes(), " "))
	v.Set("code", d.Code)
	v.Set("redirect_uri", a.c.getRedirectURL())
	v.Set("grant_type", "authorization_code")
	res, err := http.PostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", v)
	if err != nil {
		return d, fmt.Errorf("generateToken: create request: %w", err)
	}
	defer func(r io.ReadCloser) {
		_ = r.Close()
	}(res.Body)

	//debugResponse(res)

	var t TokenData
	if err := json.NewDecoder(res.Body).Decode(&t); err != nil {
		slog.With("error", err).Error("generateToken")
		return d, fmt.Errorf("generateToken: decode response body: %w", err)
	}

	d.TokenData = &t

	slog.With("token", d.TokenData, "status_code", res.StatusCode).Debug("generateToken")

	return d, nil
}

func debugResponse(res *http.Response) {
	b, _ := io.ReadAll(res.Body)
	slog.With("status_code", res.StatusCode, "response_body", string(b)).Debug("generateToken")
	res.Body = io.NopCloser(bytes.NewReader(b))
}

type authData struct {
	Code             string
	SentState        string
	ReceivedState    string
	Error            string
	ErrorDescription string
	Raw              string
	TokenData        *TokenData
}

func (a authData) TokenAsString() string {
	if a.TokenData == nil {
		return ""
	}

	b, _ := json.Marshal(a.TokenData)
	return string(b)
}

type TokenData struct {
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	ExtExpiresIn int    `json:"ext_expires_in"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	IDToken      string `json:"id_token"`
}

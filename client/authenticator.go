package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/eldius/onedrive-client/client/types"
	"github.com/eldius/onedrive-client/internal/static"
	"html/template"
	"io"
	"log/slog"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

var (
	tmpl *template.Template
)

func init() {
	var err error
	tmpl, err = template.ParseFS(static.HandlerTemplates, "templates/**")
	if err != nil {
		panic(fmt.Errorf("error parsing templates: %w", err))
	}
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

func (a *authenticator) authListener(ctx context.Context) (*authData, error) {
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

func (a *authenticator) Authenticate(ctx context.Context) (*authData, error) {
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

	return a.authListener(ctx)
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

	var t types.TokenData
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
	TokenData        *types.TokenData
}

func (a authData) TokenAsString() string {
	if a.TokenData == nil {
		return ""
	}

	b, _ := json.Marshal(a.TokenData)
	return string(b)
}

func renderAuthPage(w http.ResponseWriter, d authData) {
	if err := tmpl.ExecuteTemplate(w, "authentication_response.html", d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
	}
}

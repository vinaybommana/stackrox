package authproviders

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	"github.com/stackrox/rox/pkg/auth/tokens"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/grpc/requestinfo"
	"github.com/stackrox/rox/pkg/httputil"
)

const (
	providersPath    = "providers"
	loginPath        = "login"
	tokenRefreshPath = "tokenrefresh"
)

func (r *registryImpl) URLPathPrefix() string {
	return r.urlPathPrefix
}

func (r *registryImpl) errorURL(err error, typ string, clientState string) *url.URL {
	return &url.URL{
		Path: r.redirectURL,
		Fragment: url.Values{
			"error": {err.Error()},
			"type":  {typ},
			"state": {clientState},
		}.Encode(),
	}
}

func (r *registryImpl) tokenURL(rawToken string, typ string, clientState string) *url.URL {
	return &url.URL{
		Path: r.redirectURL,
		Fragment: url.Values{
			"token": {rawToken},
			"type":  {typ},
			"state": {clientState},
		}.Encode(),
	}
}

func (r *registryImpl) providersURLPrefix() string {
	return path.Join(r.urlPathPrefix, providersPath) + "/"
}

func (r *registryImpl) loginURLPrefix() string {
	return path.Join(r.urlPathPrefix, loginPath) + "/"
}

func (r *registryImpl) tokenRefreshPath() string {
	return path.Join(r.urlPathPrefix, tokenRefreshPath)
}

func (r *registryImpl) initHTTPMux() {
	r.HandleFunc(r.providersURLPrefix(), r.providersHTTPHandler)
	r.HandleFunc(r.loginURLPrefix(), r.loginHTTPHandler)
	if features.RefreshTokens.Enabled() {
		r.HandleFunc(r.tokenRefreshPath(), httputil.RESTHandler(r.tokenRefreshEndpoint))
	}
}

func (r *registryImpl) loginHTTPHandler(w http.ResponseWriter, req *http.Request) {
	prefix := r.loginURLPrefix()
	if !strings.HasPrefix(req.URL.Path, prefix) {
		log.Errorf("UNEXPECTED: received HTTP request for invalid URL %v", req.URL)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	providerID := req.URL.Path[len(prefix):]
	clientState := req.URL.Query().Get("clientState")

	provider := r.getAuthProvider(providerID)
	if provider == nil {
		http.Error(w, fmt.Sprintf("Unknown auth provider ID %q", providerID), http.StatusNotFound)
		return
	}

	ri := requestinfo.FromContext(req.Context())
	loginURL := provider.Backend().LoginURL(clientState, &ri)
	if loginURL == "" {
		http.Error(w, "could not get login URL", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Location", loginURL)
	w.WriteHeader(http.StatusSeeOther)
}

type tokenRefreshResponse struct {
	Token  string    `json:"token,omitempty"`
	Expiry time.Time `json:"expiry,omitempty"`
}

func (r *registryImpl) tokenRefreshEndpoint(req *http.Request) (interface{}, error) {
	refreshTokenCookie, err := req.Cookie(refreshTokenCookieName)
	if err != nil {
		return nil, httputil.Errorf(http.StatusBadRequest, "could not obtain refresh token cookie: %v", err)
	}

	var cookieData refreshTokenCookieData
	if err := cookieData.Decode(refreshTokenCookie.Value); err != nil {
		return nil, httputil.Errorf(http.StatusBadRequest, "unparseable data in refresh token cookie: %v", err)
	}

	provider := r.getAuthProvider(cookieData.ProviderID)
	if provider == nil {
		return nil, httputil.Errorf(http.StatusBadRequest, "refresh token cookie references invalid auth provider %q", cookieData.ProviderID)
	}

	if provider.Type() != cookieData.ProviderType {
		return nil, httputil.Errorf(http.StatusBadRequest, "refresh token cookie references auth provider %q of wrong type %q (expected: %q)", cookieData.ProviderID, provider.Type(), cookieData.ProviderType)
	}

	providerBackend := provider.Backend()
	if !provider.Enabled() || providerBackend == nil {
		return nil, httputil.Errorf(http.StatusInternalServerError, "auth provider %q is not currently active", provider.ID())
	}

	refreshTokenEnabledBackend, _ := providerBackend.(RefreshTokenEnabledBackend)
	if refreshTokenEnabledBackend == nil {
		return nil, httputil.Errorf(http.StatusBadRequest, "auth provider backend of type %q does not support refresh tokens", provider.Type())
	}

	authResp, err := refreshTokenEnabledBackend.RefreshAccessToken(req.Context(), cookieData.RefreshToken)
	if err != nil {
		return nil, httputil.Errorf(http.StatusInternalServerError, "failed to obtain new access token for refresh token: %v", err)
	}

	token, err := issueTokenForResponse(req.Context(), provider, authResp)
	if err != nil {
		return nil, httputil.Errorf(http.StatusInternalServerError, "failed to issue Rox token: %v", err)
	}

	return &tokenRefreshResponse{
		Token:  token.Token,
		Expiry: token.Expiry(),
	}, nil
}

func (r *registryImpl) loginURL(providerID string) string {
	return path.Join(r.loginURLPrefix(), providerID)
}

func (r *registryImpl) providersHTTPHandler(w http.ResponseWriter, req *http.Request) {
	prefix := r.providersURLPrefix()
	if !strings.HasPrefix(req.URL.Path, prefix) {
		log.Errorf("UNEXPECTED: received HTTP request for invalid URL %v", req.URL)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	relativePath := req.URL.Path[len(prefix):]
	parts := strings.SplitN(relativePath, "/", 2)
	if len(parts) == 0 {
		log.Debugf("Could not split URL path %q", req.URL.Path[len(prefix):])
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	typ := parts[0]

	factory := r.getFactory(typ)
	if factory == nil {
		log.Debugf("Factory with type %q not found", typ)
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	providerID, err := factory.ProcessHTTPRequest(w, req)
	var provider Provider
	if err == nil {
		provider = r.getAuthProvider(providerID)
		if provider == nil {
			err = fmt.Errorf("invalid auth provider ID %q", providerID)
		} else if provider.Type() != parts[0] {
			err = fmt.Errorf("auth provider %s is of invalid type %s", provider.Name(), provider.Type())
		}
	}
	if err != nil {
		if httpErr, ok := err.(httputil.HTTPError); ok {
			http.Error(w, httpErr.Error(), httpErr.HTTPStatusCode())
			return
		}
		w.Header().Set("Location", r.errorURL(err, typ, "").String())
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	authResp, clientState, err := provider.Backend().ProcessHTTPRequest(w, req)
	var tokenInfo *tokens.TokenInfo
	var refreshToken string

	if err == nil && authResp != nil {
		tokenInfo, err = issueTokenForResponse(req.Context(), provider, authResp)
		refreshToken = authResp.RefreshToken
	}

	if err != nil {
		if httpErr, ok := err.(httputil.HTTPError); ok {
			http.Error(w, httpErr.Error(), httpErr.HTTPStatusCode())
			return
		}
		w.Header().Set("Location", r.errorURL(err, typ, clientState).String())
		w.WriteHeader(http.StatusSeeOther)
		return
	}

	if tokenInfo == nil {
		// Assume the ProcessHTTPRequest already took care of writing a response.
		return
	}

	w.Header().Set("Location", r.tokenURL(tokenInfo.Token, typ, clientState).String())
	if refreshToken != "" && features.RefreshTokens.Enabled() {
		cookieData := refreshTokenCookieData{
			ProviderType: typ,
			ProviderID:   providerID,
			RefreshToken: refreshToken,
		}
		if encodedData, err := cookieData.Encode(); err != nil {
			log.Errorf("failed to encode refresh token cookie data: %v", err)
		} else {
			refreshTokenCookie := &http.Cookie{
				Name:     refreshTokenCookieName,
				Value:    encodedData,
				Path:     r.tokenRefreshPath(),
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			}
			http.SetCookie(w, refreshTokenCookie)
		}
	}

	w.WriteHeader(http.StatusSeeOther)
}

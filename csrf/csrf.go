package csrf

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

type ctxKeyType string

const ctxKey ctxKeyType = "CSRF Token Context key"

// FormsTokenKey is the key to use in forms for CSRF tokens.
const FormsTokenKey = "anti_csrf_token"

// Protect protects state-changing requests on the given handler from CSRF.
func Protect(h http.Handler) http.Handler {
	const cookieName = "anti-CSRF"

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(cookieName)
		var token string
		if err != nil || cookie.Value == "" {
			token = genToken()
			c := http.Cookie{
				Name:  cookieName,
				Value: token,
			}
			http.SetCookie(w, &c)
		} else {
			token = cookie.Value
		}
		r = r.WithContext(context.WithValue(r.Context(), ctxKey, token))
		if isStatePreserving(r) {
			// This is a safe request, allow
			h.ServeHTTP(w, r)
			return
		}
		// This is a state-changing request, check token.
		if v := r.FormValue(FormsTokenKey); v != token {
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			return
		}
		h.ServeHTTP(w, r)
	})
}

// GetToken retrieves the anti-CSRF token for the current request.
func GetToken(r *http.Request) string {
	return r.Context().Value(ctxKey).(string)
}

// AddData adds the CSRF token to m keyed on FormsTokenKey.
func AddData(r *http.Request, m map[string]interface{}) map[string]interface{} {
	m[FormsTokenKey] = GetToken(r)
	return m
}

func genToken() string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate CSRF token: %v", err))
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

func isStatePreserving(r *http.Request) bool {
	switch r.Method {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return true
	default:
		return false
	}
}

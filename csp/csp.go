package csp

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
)

type ctxKeyType string

const ctxKey ctxKeyType = "CSP nonce context key"

// Protect protects an handler with CSP
func Protect(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		nonce := genNonce()
		w.Header().Set("Content-Security-Policy", "object-src 'none';"+
			"script-src 'nonce-"+nonce+"' 'unsafe-inline' 'unsafe-eval' 'strict-dynamic' https: http:;"+
			"base-uri 'none';"+
			"report-uri https://your-report-collector.example.com/")
		r = r.WithContext(context.WithValue(r.Context(), ctxKey, nonce))
		h.ServeHTTP(w, r)
	}
}

func genNonce() string {
	b := make([]byte, 20)
	_, err := rand.Read(b)
	if err != nil {
		panic(fmt.Errorf("failed to generate entropy using crypto/rand/RandReader: %v", err))
	}
	return base64.RawStdEncoding.EncodeToString(b)
}

// GetNonce retrieves the nonce from the current context.
func GetNonce(r *http.Request) string {
	return r.Context().Value(ctxKey).(string)
}
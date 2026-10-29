// Package middleware is the authentication / authorization layer of
// lms-catalog-api, reduced to what the W13 security story needs.
package middleware

import (
	"context"
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

// Scopes of the two kinds of token that reach lms-catalog-api.
const (
	// ScopeCopiesWrite is carried by the internal HS256 token that
	// lms-circulation-api mints for loan-copy / return-copy (HU-06, HU-07).
	ScopeCopiesWrite = "catalog:copies:write"
	// ScopeBooksWrite is carried by the Administrator RS256 session token
	// issued by lms-access-api (HU-04, HU-09).
	ScopeBooksWrite = "catalog:books:write"
)

type scopesKey struct{}

// RequireAuth is the team's existing middleware (RS256 for the Administrator,
// HS256 with INTERNAL_JWT_SECRET for service calls, algorithm pinned per key),
// plus one change: it keeps the token's scopes in the request context.
func RequireAuth(publicKey *rsa.PublicKey, internalSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			parts := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication token required")
				return
			}

			claims := jwt.MapClaims{}
			token, err := jwt.ParseWithClaims(parts[1], claims, func(t *jwt.Token) (interface{}, error) {
				switch t.Method.(type) {
				case *jwt.SigningMethodRSA:
					return publicKey, nil
				case *jwt.SigningMethodHMAC:
					return []byte(internalSecret), nil
				default:
					return nil, jwt.ErrTokenSignatureInvalid
				}
			}, jwt.WithValidMethods([]string{jwt.SigningMethodRS256.Alg(), jwt.SigningMethodHS256.Alg()}))
			if err != nil || !token.Valid {
				writeError(w, http.StatusUnauthorized, "UNAUTHORIZED", "Authentication token required")
				return
			}

			scope, _ := claims["scope"].(string)
			ctx := context.WithValue(r.Context(), scopesKey{}, strings.Fields(scope))
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequireScope answers 403 when an authenticated token was not issued for
// this route — the missing control of STRIDE row E (PR #36).
func RequireScope(scope string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			scopes, _ := r.Context().Value(scopesKey{}).([]string)
			for _, s := range scopes {
				if s == scope {
					next.ServeHTTP(w, r)
					return
				}
			}
			writeError(w, http.StatusForbidden, "FORBIDDEN", "Token not allowed on this route")
		})
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(`{"error":"` + code + `","message":"` + message + `"}`))
}

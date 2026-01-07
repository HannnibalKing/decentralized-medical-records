package middleware

import (
	"net/http"
)

// DPoP is a placeholder middleware for DPoP proof verification.
// TODO: implement DPoP proof validation and nonce replay protection.
func DPoP(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("DPoP")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"error":"dpop_required"}`))
			return
		}
		// TODO: validate JWK thumbprint binding and nonce replay protection.
		next.ServeHTTP(w, r)
	})
}

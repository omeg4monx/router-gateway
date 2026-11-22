package auth

import (
	"net/http"
	"strings"
	"time"
)

// SessionCookieName is the HTTP cookie name carrying the signed admin session token.
const SessionCookieName = "om4_session"

func isSecureRequest(r *http.Request) bool {
	if r == nil {
		return false
	}
	if r.TLS != nil {
		return true
	}
	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		first, _, _ := strings.Cut(proto, ",")
		if strings.EqualFold(strings.TrimSpace(first), "https") {
			return true
		}
	}
	return false
}

func setSessionCookie(w http.ResponseWriter, r *http.Request, token string, ttl time.Duration, expiresAt time.Time) {
	maxAge := int(ttl.Seconds())
	if maxAge < 0 {
		maxAge = 0
	}
	//nolint:gosec // G124: Secure is set dynamically for TLS/X-Forwarded-Proto=https so localhost HTTP works.
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expiresAt.UTC(),
		MaxAge:   maxAge,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func clearSessionCookie(w http.ResponseWriter, r *http.Request) {
	//nolint:gosec // G124: Secure is set dynamically for TLS/X-Forwarded-Proto=https so localhost HTTP works.
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    "",
		Path:     "/",
		Expires:  time.Unix(0, 0).UTC(),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isSecureRequest(r),
		SameSite: http.SameSiteLaxMode,
	})
}

func readSessionCookie(r *http.Request) (string, error) {
	if r == nil {
		return "", ErrUnauthenticated
	}
	cookie, err := r.Cookie(SessionCookieName)
	if err != nil || strings.TrimSpace(cookie.Value) == "" {
		return "", ErrUnauthenticated
	}
	return cookie.Value, nil
}

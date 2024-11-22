package plauth

import (
	"errors"
	"net/http"
	"time"
)

func SetSessionTokenCookie(w http.ResponseWriter, token string, expiresAt time.Time, isProd bool) {
	cookie := http.Cookie{
		Name:     "plauth.session-token",
		Value:    token,
		Path:     "/",
		Expires:  expiresAt,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
	}

	http.SetCookie(w, &cookie)
}

func GetSessionTokenFromCookie(r *http.Request) (string, error) {
	cookie, err := r.Cookie("plauth.session-token")
	if err != nil || cookie.Value == "" {
		return "", errors.New("invalid cookie")
	}
	return cookie.Value, nil
}

func DeleteSessionTokenCookie(w http.ResponseWriter, isProd bool) {
	cookie := http.Cookie{
		Name:     "plauth.session-token",
		Path:     "/",
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   0,
	}

	http.SetCookie(w, &cookie)
}

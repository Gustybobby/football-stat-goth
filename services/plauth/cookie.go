package plauth

import (
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
	if err != nil {
		return "", err
	}
	return cookie.Value, err
}

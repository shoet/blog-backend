package handlers

import "net/http"

type Cookie struct {
	Key   string
	Value string
}

func SetCookie(w http.ResponseWriter, key string, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   1000 * 60 * 60 * 24 * 14,
		Path:     "/",
	})
}

func ClearCookie(w http.ResponseWriter, key string) {
	http.SetCookie(w, &http.Cookie{
		Name:   key,
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	})
}

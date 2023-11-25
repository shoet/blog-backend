package handlers

import (
	"fmt"
	"net/http"

	"github.com/shoet/blog/config"
)

func SetCookie(cfg *config.Config, w http.ResponseWriter, key string, value string) error {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   1000 * 60 * 60 * 24 * 7,
		Path:     "/",
	}
	if cfg.Env == "prod" {
		if cfg.SiteDomain == "" {
			fmt.Println("site domain is not set")
		}
		cookie.Secure = true
		cookie.Domain = cfg.SiteDomain
	}
	http.SetCookie(w, cookie)
	return nil
}

func ClearCookie(cfg *config.Config, w http.ResponseWriter, key string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
	}
	if cfg.Env == "prod" {
		if cfg.SiteDomain == "" {
			fmt.Println("site domain is not set")
		}
		cookie.Secure = true
		cookie.Domain = cfg.SiteDomain
	}
	http.SetCookie(w, cookie)
}

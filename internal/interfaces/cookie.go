package interfaces

import (
	"fmt"
	"net/http"
)

type CookieManager struct {
	Env        string
	SiteDomain string
}

func NewCookieManager(env string, siteDomain string) *CookieManager {
	return &CookieManager{
		Env:        env,
		SiteDomain: siteDomain,
	}
}

func (c *CookieManager) SetCookie(w http.ResponseWriter, key string, value string) error {
	cookie := &http.Cookie{
		Name:     key,
		Value:    value,
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   1000 * 60 * 60 * 24 * 7,
		Path:     "/",
	}
	if c.Env == "prod" {
		if c.SiteDomain == "" {
			fmt.Println("site domain is not set")
		}
		cookie.Secure = true
		cookie.Domain = c.SiteDomain
	}
	http.SetCookie(w, cookie)
	return nil
}

func (c *CookieManager) ClearCookie(w http.ResponseWriter, key string) {
	cookie := &http.Cookie{
		Name:     key,
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
		Path:     "/",
	}
	if c.Env == "prod" {
		if c.SiteDomain == "" {
			fmt.Println("site domain is not set")
		}
		cookie.Secure = true
		cookie.Domain = c.SiteDomain
	}
	http.SetCookie(w, cookie)
}

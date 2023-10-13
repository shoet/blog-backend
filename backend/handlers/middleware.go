package handlers

import "net/http"

var whiteList = []string{
	"http://localhost:5173",
	"http://localhost:3000",
	"http://localhost:6006",
}

func CORSMiddleWare(next http.Handler) http.Handler {
	// TODO: ブラッシュアップ
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if originAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET,PUT,POST,DELETE,UPDATE,OPTIONS")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func originAllowed(origin string) bool {
	for _, allowedOrigin := range whiteList {
		if allowedOrigin == origin {
			return true
		}
	}
	return false
}

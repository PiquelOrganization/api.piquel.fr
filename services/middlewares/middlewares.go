package middlewares

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/PiquelChips/piquel.fr/services/auth"
	"github.com/PiquelChips/piquel.fr/services/config"
	"github.com/gorilla/mux"
)

func SetupMiddlewares(router *mux.Router) {
	//router.Use(authMiddleware)
    router.Use(mux.CORSMethodMiddleware(router))
    router.Use(cORSMiddleware)
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.HasPrefix(r.URL.Path, "/auth/") {
            next.ServeHTTP(w, r)
            return
        }

        _, err := auth.GetSessionUser(r)
        if err != nil {
            http.Error(w, "You are not authenticated", http.StatusForbidden)
            return
        }

        next.ServeHTTP(w, r)
	})
}

func cORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
        log.Printf("[CORS] %s requesting...", origin)

		if origin == "" {
			if r.Host != r.Header.Get("Host") {
				http.Error(w, "Missing Origin Header", http.StatusForbidden)
				return
			}

			// Same origin
			next.ServeHTTP(w, r)
			return
		}

        isValidOrigin, allowCredentials := validateOrigin(origin, config.Config.CORS.AllowedOrigins)

		if !isValidOrigin {
			http.Error(w, "Origin not allowed", http.StatusUnauthorized)
            log.Printf("[CORS] Just rejected %s! This origin is unauthorized!", origin)
            return
		}

        w.Header().Set("Access-Control-Allow-Origin", origin)
        w.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.Config.CORS.MaxAge))
        if allowCredentials {
            w.Header().Set("Access-Control-Allow-Credentials", "true")
        }

		if r.Method == http.MethodOptions {
			// Return immediately for OPTIONS requests
			w.WriteHeader(http.StatusOK)
			return
		}

        log.Printf("[CORS] Just allowed request from %s!", origin)

		next.ServeHTTP(w, r)
	})
}

func validateOrigin(origin string, allowedOrigins map[string]bool) (bool, bool) {
	if origin == "" {
		return false, false
	}

	for allowed, credentials := range allowedOrigins {
		if allowed == origin {
			return true, credentials
		}

		// For example *.piquel.fr
		if strings.Contains(allowed, "*.") {
			// Would then be .piquel.fr
			domain := strings.Split(allowed, "*.")[1]
			if strings.HasSuffix(origin, domain) {
				return true, credentials
			}
		}
	}

	return false, false
}

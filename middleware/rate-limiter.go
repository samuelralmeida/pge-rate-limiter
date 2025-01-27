package middleware

import (
	"context"
	"log"
	"net/http"
	"strings"
)

type rateLimiter interface {
	IsAllow(ctx context.Context, ip, token string) (bool, error)
}

func RateLimit(ctx context.Context, limiter rateLimiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := strings.Split(r.RemoteAddr, ":")[0]
			token := r.Header.Get("API_KEY")

			allowed, err := limiter.IsAllow(ctx, ip, token)

			if err != nil {
				log.Println(err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("limiter internal error"))
				return
			}

			if !allowed {
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

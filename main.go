package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
	"github.com/samuelralmeida/pge-rate-limiter/limiter"
	mw "github.com/samuelralmeida/pge-rate-limiter/middleware"
	"github.com/samuelralmeida/pge-rate-limiter/storage/redis"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("no .env file loaded")
	}
}

func main() {
	ctx := context.Background()

	redisStorage := redis.NewRedisStorage()
	rateLimit := limiter.NewLimiter(redisStorage)

	rateLimiterMiddleware := mw.RateLimit(ctx, rateLimit)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(rateLimiterMiddleware)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	log.Println("listening on port", os.Getenv("APP_PORT"))
	http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("APP_PORT")), r)
}

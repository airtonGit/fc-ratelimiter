package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/airtongit/fc-ratelimiter/internal/ratelimiter"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {

	if _, err := os.Stat(".env"); err == nil {
		err := godotenv.Load()
		if err != nil {
			fmt.Println("Error loading .env file", err)
			return
		}
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		fmt.Println("Error loading redis DB config", err)
		return
	}

	fmt.Println(fmt.Sprintf("Loaded config redis_host:%s, redis_db:%d", redisHost, redisDB))

	redisCache := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       redisDB,
	})

	if _, err := redisCache.Ping(context.Background()).Result(); err != nil {
		fmt.Println("Error connecting to redis", err)
		return
	}

	redisRepository := database.NewRedisRepository(redisCache)

	IPLimitSec, err := strconv.Atoi(os.Getenv("IP_LIMIT_SEC"))
	if err != nil {
		panic("Error loading IP Limit Sec config")
	}

	csvlimitTokensList := os.Getenv("TOKENS_LIMIT_SEC")
	// convert csv list to []string
	csvReader := csv.NewReader(strings.NewReader(csvlimitTokensList))
	csvReader.Comma = ','
	// parse csv (csvlimitTokensList)
	tokensList, err := csvReader.Read()
	if err != nil {
		panic("Error loading tokens list")
	}

	tokenRateLimitMapc := make(map[string]int)
	for _, tokenLimitPair := range tokensList {
		values := strings.Split(tokenLimitPair, "=")
		tokenItem := values[0]
		tokenLimit, err := strconv.Atoi(values[1])
		if err != nil {
			panic("Error loading token limit")
		}
		tokenRateLimitMapc[tokenItem] = tokenLimit
	}

	rateLimiterUsecase := ratelimiter.NewRateLimiterUsecase(redisRepository)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				fmt.Println("Error parsing remote address:", err)
				return
			}

			usecaseInputDTO := ratelimiter.AllowRateLimitInputDTO{
				IPRequestBySecondLimit:     IPLimitSec,
				IP:                         host,
				Token:                      r.Header.Get("API_KEY"),
				TokenRequestsBySecondLimit: tokenRateLimitMapc,
				TokenDuration:              time.Second,
				IpDuration:                 time.Second,
			}

			rateLimitOutput := rateLimiterUsecase.Execute(r.Context(), usecaseInputDTO)
			if !rateLimitOutput.Allow {
				//fmt.Println("Rate limit exceeded 429")
				w.WriteHeader(http.StatusTooManyRequests)
				w.Write([]byte("you have reached the maximum number of requests or actions allowed within a certain time frame"))
				return
			}

			next.ServeHTTP(w, r)
		})
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})
	fmt.Println("Listening on :8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		fmt.Println(err)
		return
	}
}

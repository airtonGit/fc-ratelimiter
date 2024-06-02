package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/airtongit/fc-ratelimiter/internal/infrastructure/database"
	"github.com/airtongit/fc-ratelimiter/internal/ratelimiter"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	redisHost := os.Getenv("REDIS_HOST")
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		panic("Error loading redis DB config")
	}

	redisCache := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "",
		DB:       redisDB,
	})

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
	type TokenRateLimit struct {
		Token string
		Limit int
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
	//r.Use(middleware.Logger)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			usecaseInputDTO := ratelimiter.AllowRateLimitInputDTO{
				IpLimit:       IPLimitSec,
				IP:            strings.Split(r.RemoteAddr, ":")[0],
				Token:         r.Header.Get("API_KEY"),
				TokenLimit:    tokenRateLimitMapc,
				TokenDuration: time.Second,
				IpDuration:    time.Second,
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

package middleware

import (
	"net/http"
	"time"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/zyncc/ecommerce-microservice/services/api-gateway/pkg/utils"
)

func RateLimiter(
	rdb *redis.Client,
	capacity int,
	refillRate float64, // tokens per second
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			ip := chimiddleware.GetClientIP(ctx)

			allowed, err := tokenBucketScript.Run(
				ctx,
				rdb,
				[]string{"rate_limit:" + ip},
				capacity,
				refillRate,
				float64(time.Now().Unix()),
			).Bool()
			if err != nil {
				utils.ErrorResponse(w, http.StatusInternalServerError, "something went wrong")
				return
			}

			if !allowed {
				utils.ErrorResponse(w, http.StatusTooManyRequests, "too many requests")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

var tokenBucketScript = redis.NewScript(`
local key = KEYS[1]

local capacity = tonumber(ARGV[1])
local refillRate = tonumber(ARGV[2]) -- tokens per second
local now = tonumber(ARGV[3])

local data = redis.call("HMGET", key, "tokens", "last")

local tokens = tonumber(data[1])
local last = tonumber(data[2])

if tokens == nil then
    tokens = capacity
    last = now
end

local elapsed = math.max(0, now - last)
tokens = math.min(capacity, tokens + elapsed * refillRate)

if tokens < 1 then
    redis.call("HMSET", key,
        "tokens", tokens,
        "last", now
    )
    redis.call("EXPIRE", key, math.ceil(capacity / refillRate))
    return 0
end

tokens = tokens - 1

redis.call("HMSET", key,
    "tokens", tokens,
    "last", now
)

redis.call("EXPIRE", key, math.ceil(capacity / refillRate))

return 1
`)

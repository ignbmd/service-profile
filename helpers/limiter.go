package helpers

import (
	"context"

	"github.com/redis/go-redis/v9"
	"smartbtw.com/services/profile/db"
)

type RateLimiter struct {
	client *redis.ClusterClient
	limit  int64
	window int64
}

// Convert second to millisecond
func InSecond(second int64) int64 {
	return second * 1000
}

// Convert minute to millisecond
func InMinute(minute int64) int64 {
	return InSecond(minute * 60)
}

// Convert hour to millisecond
func InHour(hour int64) int64 {
	return InMinute(hour * 60)
}

// Create new rate limiter instance
func NewLimiter(limit int64, window int64) *RateLimiter {
	return &RateLimiter{
		client: db.NewRedisCluster(),
		limit:  limit,
		window: window,
	}
}

// Check if key is allowed to access
func (r *RateLimiter) IsAllowed(key string) (bool, error) {
	script := `
		local current
		current = redis.call("incr",KEYS[1])
		if tonumber(current) == 1 then
			redis.call("pexpire",KEYS[1], ARGV[1])
		end
		return current
	`

	res, error := r.client.Eval(context.Background(), script, []string{key}, r.window).Result()

	if error != nil {
		return false, error
	}

	if res.(int64) > r.limit {
		return false, nil
	}

	return true, nil
}

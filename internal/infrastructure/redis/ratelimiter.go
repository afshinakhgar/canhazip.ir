package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RateLimiter implements a sliding window rate limiter backed by Redis sorted sets.
type RateLimiter struct {
	rdb *redis.Client
}

// NewRateLimiter creates a RateLimiter using the provided Redis client.
func NewRateLimiter(rdb *redis.Client) *RateLimiter {
	return &RateLimiter{rdb: rdb}
}

// Allow returns true if the request should be allowed under the given limit
// within the specified window. The key should uniquely identify the caller
// (e.g. IP address + route prefix).
func (r *RateLimiter) Allow(ctx context.Context, key string, limit int, window time.Duration) (bool, error) {
	now := time.Now().UnixMilli()
	windowStart := now - window.Milliseconds()
	redisKey := fmt.Sprintf("ratelimit:%s", key)

	pipe := r.rdb.Pipeline()

	// Remove entries outside the window.
	pipe.ZRemRangeByScore(ctx, redisKey, "0", fmt.Sprintf("%d", windowStart))

	// Count remaining entries in the current window.
	countCmd := pipe.ZCount(ctx, redisKey, fmt.Sprintf("%d", windowStart), "+inf")

	// Add current request as a member (member = timestamp in ms as string).
	member := fmt.Sprintf("%d", now)
	pipe.ZAdd(ctx, redisKey, redis.Z{Score: float64(now), Member: member})

	// Set key expiry to window duration so Redis cleans up idle keys.
	pipe.Expire(ctx, redisKey, window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, fmt.Errorf("ratelimiter: pipeline exec: %w", err)
	}

	count, err := countCmd.Result()
	if err != nil {
		return false, fmt.Errorf("ratelimiter: count result: %w", err)
	}

	return count < int64(limit), nil
}

package requestlog

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

const key = "canhazip:requests"
const maxEntries = 500

// Entry represents a single captured HTTP request.
type Entry struct {
	Timestamp string `json:"ts"`
	IP        string `json:"ip"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Status    int    `json:"status"`
	LatencyMs int64  `json:"latency_ms"`
	UserAgent string `json:"ua"`
}

// Logger is a Redis-backed ring buffer of recent requests.
type Logger struct {
	rdb *redis.Client
}

// New creates a new Logger backed by the given Redis client.
func New(rdb *redis.Client) *Logger {
	return &Logger{rdb: rdb}
}

// Push stores one entry (LPUSH + LTRIM to maxEntries). Best-effort, ignore errors.
func (l *Logger) Push(entry Entry) {
	data, err := json.Marshal(entry)
	if err != nil {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	pipe := l.rdb.Pipeline()
	pipe.LPush(ctx, key, data)
	pipe.LTrim(ctx, key, 0, maxEntries-1)
	_, _ = pipe.Exec(ctx)
}

// Recent returns up to n latest entries (LRANGE 0 n-1), newest first.
func (l *Logger) Recent(n int64) []Entry {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	vals, err := l.rdb.LRange(ctx, key, 0, n-1).Result()
	if err != nil {
		return nil
	}
	entries := make([]Entry, 0, len(vals))
	for _, v := range vals {
		var e Entry
		if err := json.Unmarshal([]byte(v), &e); err == nil {
			entries = append(entries, e)
		}
	}
	return entries
}

// TotalCount returns LLEN of the list.
func (l *Logger) TotalCount() int64 {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	n, _ := l.rdb.LLen(ctx, key).Result()
	return n
}

package banlist

import (
	"context"
	"sort"
	"time"

	"github.com/redis/go-redis/v9"
)

const key = "canhazip:banned"

// BanList is a Redis-backed set of banned IPs.
type BanList struct {
	rdb *redis.Client
}

// New creates a new BanList backed by the given Redis client.
func New(rdb *redis.Client) *BanList {
	return &BanList{rdb: rdb}
}

func ctx() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 2*time.Second)
}

// Ban adds ip to the ban set. Returns error only on Redis failure.
func (b *BanList) Ban(ip string) error {
	c, cancel := ctx()
	defer cancel()
	return b.rdb.SAdd(c, key, ip).Err()
}

// Unban removes ip from the ban set.
func (b *BanList) Unban(ip string) error {
	c, cancel := ctx()
	defer cancel()
	return b.rdb.SRem(c, key, ip).Err()
}

// IsBanned returns true if ip is in the ban set. Fail-open (returns false on error).
func (b *BanList) IsBanned(ip string) bool {
	c, cancel := ctx()
	defer cancel()
	ok, err := b.rdb.SIsMember(c, key, ip).Result()
	if err != nil {
		return false
	}
	return ok
}

// List returns all banned IPs sorted alphabetically.
func (b *BanList) List() ([]string, error) {
	c, cancel := ctx()
	defer cancel()
	members, err := b.rdb.SMembers(c, key).Result()
	if err != nil {
		return nil, err
	}
	sort.Strings(members)
	return members, nil
}

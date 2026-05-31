package blocklist

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	ipdomain "iplocation.sabaai.ir/internal/domain/ip"
)

// Checker checks IP reputation against local Firehol blocklists.
type Checker struct {
	mu     sync.RWMutex
	level1 []*net.IPNet
	level2 []*net.IPNet
	l1Path string
	l2Path string
}

// NewChecker loads level1 and level2 netset files into memory.
func NewChecker(level1Path, level2Path string) (*Checker, error) {
	l1, err := loadNetset(level1Path)
	if err != nil {
		return nil, fmt.Errorf("blocklist: load level1: %w", err)
	}
	l2, err := loadNetset(level2Path)
	if err != nil {
		return nil, fmt.Errorf("blocklist: load level2: %w", err)
	}
	return &Checker{
		level1: l1,
		level2: l2,
		l1Path: level1Path,
		l2Path: level2Path,
	}, nil
}

// Reload re-reads both files and swaps the slices atomically.
func (c *Checker) Reload() error {
	l1, err := loadNetset(c.l1Path)
	if err != nil {
		return fmt.Errorf("blocklist: reload level1: %w", err)
	}
	l2, err := loadNetset(c.l2Path)
	if err != nil {
		return fmt.Errorf("blocklist: reload level2: %w", err)
	}
	c.mu.Lock()
	c.level1 = l1
	c.level2 = l2
	c.mu.Unlock()
	return nil
}

// Stats returns current CIDR counts and file paths.
func (c *Checker) Stats() (l1Count, l2Count int, l1Path, l2Path string) {
	c.mu.RLock()
	l1Count = len(c.level1)
	l2Count = len(c.level2)
	c.mu.RUnlock()
	return l1Count, l2Count, c.l1Path, c.l2Path
}

// Level1Path returns the path to the level1 netset file.
func (c *Checker) Level1Path() string { return c.l1Path }

// Level2Path returns the path to the level2 netset file.
func (c *Checker) Level2Path() string { return c.l2Path }

// Lookup satisfies ipdomain.Repository (stub — use LookupReputation).
func (c *Checker) Lookup(_ string) (*ipdomain.IPInfo, error) {
	return nil, fmt.Errorf("blocklist: Lookup not implemented")
}

// LookupASN satisfies ipdomain.Repository (stub).
func (c *Checker) LookupASN(_ string) (uint, string, error) {
	return 0, "", fmt.Errorf("blocklist: LookupASN not implemented")
}

// LookupReputation checks IP against firehol_level1 then level2.
// level1 → malicious (score 100), level2 → suspicious (score 50), else clean.
func (c *Checker) LookupReputation(ipStr string) (*ipdomain.Reputation, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return &ipdomain.Reputation{Score: 0, Label: "unknown"}, nil
	}

	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, cidr := range c.level1 {
		if cidr.Contains(ip) {
			return &ipdomain.Reputation{Score: 100, Label: "malicious"}, nil
		}
	}
	for _, cidr := range c.level2 {
		if cidr.Contains(ip) {
			return &ipdomain.Reputation{Score: 50, Label: "suspicious"}, nil
		}
	}
	return &ipdomain.Reputation{Score: 0, Label: "clean"}, nil
}

func loadNetset(path string) ([]*net.IPNet, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var nets []*net.IPNet
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if !strings.Contains(line, "/") {
			line = line + "/32"
		}
		_, ipnet, err := net.ParseCIDR(line)
		if err != nil {
			continue
		}
		nets = append(nets, ipnet)
	}
	return nets, scanner.Err()
}

package redis

import (
	"strings"
	"time"

	goredis "github.com/go-redis/redis"
)

var (
	// DefaultTTL used if config ttl is zero value to prevent non-expiring keys
	DefaultTTL = time.Duration(time.Second * 86400)

	// DefaultRetries used if zero value in config
	DefaultRetries = 10

	// Error returned by goredis if client config is incorrect (cluster vs default)
	BadClientType = "ERR This instance has cluster support disabled"
)

// NewClient creates either default or clustered client
func NewClient(c Config) Client {
	if c.TTL == 0 {
		c.TTL = DefaultTTL
	}
	if c.Retries == 0 {
		c.Retries = DefaultRetries
	}
	addr := check(c.Addr)
	if c.Clustered {
		return goredis.NewClusterClient(&ClusterOptions{
			Addrs:      addr,
			MaxRetries: c.Retries,
			Password:   "",
		})
	}
	return goredis.NewClient(&Options{
		Addr:       addr[0],
		MaxRetries: c.Retries,
		Password:   "",
	})
}

// Client interface to satisfy regular redis client and clustered client
type Client interface {
	Get(string) *StringCmd
	MGet(...string) *SliceCmd
	Scan(uint64, string, int64) *ScanCmd
	Set(string, interface{}, time.Duration) *StatusCmd
	Del(...string) *IntCmd
	Ping() *StatusCmd
	Close() error
}
type (
	Options        = goredis.Options
	ClusterOptions = goredis.ClusterOptions

	IntCmd    = goredis.IntCmd
	SliceCmd  = goredis.SliceCmd
	StatusCmd = goredis.StatusCmd
	StringCmd = goredis.StringCmd
	ScanCmd   = goredis.ScanCmd
)

func check(addr []string) []string {
	for i, v := range addr {
		if !strings.HasSuffix(v, ":6379") {
			addr[i] = v + ":6379"
		}
	}
	return addr
}

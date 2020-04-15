package redis

import (
	"errors"
	"log"
	"time"

	goredis "github.com/go-redis/redis"
)

// Open creates a new redis agent
func Open(c Config) (*Agent, error) {
	agent := &Agent{
		client: NewClient(c),
		config: c,
	}
	return agent, agent.dial()
}

// Agent controls and maintains db data source
type Agent struct {
	client Client
	config Config
}

// Config struct for agent
type Config struct {
	Addr      []string      `json:"addr"`
	Clustered bool          `json:"cluster"`
	Retries   int           `json:"retries"`
	Scope     string        `json:"scope"`
	TTL       time.Duration `json:"ttl"`
}

func (a *Agent) dial() error {
	log.Println("dba: dialing...")
	if _, err := a.client.Ping().Result(); err != nil {
		switch err.Error() {
		case BadClientType:
			log.Println("dba: dial err: clustered client failed; trying default client")
			a.config.Clustered = false
			a.client = NewClient(a.config)
			if _, err := a.client.Ping().Result(); err != nil {
				log.Printf("dba: dial err; retry with default client err: %s", err.Error())
				return err
			}
			log.Println("dba: dial ok; successful retry with default client")
		default:
			log.Printf("dba: dial err: %s", err.Error())
			return err
		}
	}
	return nil
}

var (
	ErrBadAddr   = errors.New("dba: bad addr")
	ErrBadConfig = errors.New("dba: bad config")
	ErrBadKey    = errors.New("dba: bad key")
)

// GetBytes returns byte array based on provided key
func (a *Agent) GetBytes(key string) ([]byte, error) {
	b, err := a.client.Get(key).Bytes()
	if err == nil {
		return b, nil
	}
	switch err {
	case goredis.Nil:
		return nil, ErrBadKey
	default:
		return nil, err
	}
}

// SetBytes will put the kvp into redis at the exp time in seconds
func (a *Agent) SetBytes(key string, data []byte) error {
	return a.client.Set(key, data, a.config.TTL).Err()
}

// Scan for all keys based on provided pattern
func (a *Agent) Scan(value string, count int64) ([]string, error) {
	if count == 0 {
		count = 10
	}
	var data []string
	var cursor uint64
	for {
		var err error
		data, cursor, err = a.client.Scan(cursor, value, count).Result()
		if err != nil {
			panic(err)
		}
		if cursor == 0 {
			break
		}
	}
	return data, nil
}

// MGet fetches records from db
func (a *Agent) MGet(key ...string) ([]interface{}, error) {
	return a.client.MGet(key...).Result()
}

// Del removes db record by key
func (a *Agent) Del(k ...string) error {
	return a.client.Del(k...).Err()
}

// Close func stops agent
func (a *Agent) Close() {
	if a.client != nil {
		a.client.Close()
	} else {
		log.Println("close: nil agent... closing")
	}
}

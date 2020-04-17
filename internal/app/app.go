package app

import (
	"context"
	"encoding/gob"
	"os"

	"github.com/BethesdaNet/friends-go/internal/db/redis"
	"github.com/BethesdaNet/friends-go/internal/handler"
	"github.com/BethesdaNet/friends-go/internal/metric/bio"
	"github.com/BethesdaNet/friends-go/internal/metric/relic"
	"github.com/BethesdaNet/friends-go/internal/provider"
)

const (
	// DisableSplunkLogging re-enables the human-readable service logs and disables
	// splunk output. This is useful when you are debugging the system and want to
	// know what the internals are doing.
	DisableSplunkLogging = false

	// DefaultRedisTTL is how long the record should persist in redis
	DefaultRedisTTL = 18600

	// APIVersion specifies the allowed api version on request
	APIVersion = "v1"
)

// Open creates the presence service
func Open(conf Config, dba *redis.Agent, nra *relic.Agent) (*Friends, error) {
	if DisableSplunkLogging {
		conf.SimpleLog = true
	}
	if conf.Redis.TTL == 0 {
		conf.Redis.TTL = DefaultRedisTTL
	}
	if conf.Provider == nil {
		conf.Provider = make(map[string]interface{})
	}
	f := &Friends{
		config: conf,
		done:   make(chan struct{}),
	}

	// create manager instance and set relic db agent. NOTE(xc): create "new" func
	// to handle creation of the manager if any new complexities (maps, cache) are
	// implemented in the future.
	f.manager = &Manager{
		dba: dba,
		//notes: make(chan Notification),
		done: make(chan struct{}),
	}

	// create new logger and redirect it to stderr or the wanted pipe output
	f.Log = bio.NewLogger(nil, os.Stderr)

	// if the new relic agent is not nil attach to service
	if nra != nil {
		f.nra = nra
	}

	f.handler = &handler.Handler{f}

	return f, f.init()
}

// Presence is the handler, processor, and data layer. It brokers inbound http
// requests from clients to maintain presence status.
type Friends struct {
	ctx        context.Context
	done       chan struct{}
	root       string
	file, path string
	httpok     int
	err        error

	// handler
	handler *handler.Handler

	// manager brokers communications between inbound http handlers, external data
	// providers (notifications, identity), and external data sources (redis).
	manager *Manager

	// state contains attrs set by inbound presence status http requests which can
	// change system behavior based on data/parameters.
	state State

	// config contains all required settings to run and maintain presence system
	// and components accessing external resources
	config Config

	// Log (bio.Logger) handles system logging efficently (simple, or verbose) data. Bio
	// logger was built to handle batch processing or realtime based on low/hi water.
	// Setting Config.SimpleLog / DisableSplunkLogging (for local dev) attrs to reduce
	// volume and noise.
	Log *bio.Logger

	// nra wraps newrelics agent to handle error cases where the default agent panics
	// when invalid keys are set causing fatal tasks
	nra *relic.Agent
}

// Config struct ...
type Config struct {
	Name string `json:"name"`
	Addr string `json:"addr"`
	Env  string `json:"env"`

	// Target is the registered ingress ALB target group endpoint. This attribute
	// gets a suffix appended based on private bool (*-private)
	Target string `json:"target"`

	// Private changes target to contain the "-private" prefix for the task and
	// modifies the behavior of the presence service.
	Private bool `json:"private"`

	// Redis contains settings for redis agent
	Redis redis.Config

	// Relic contains settings for newrelic agent
	Relic relic.Config

	// Provider map holds onto provider configurations by name
	Provider map[string]interface{} `json:"provider"`

	// Meta contains attributes specific for aws ecr / docker which are then added
	// to http request logging. These fields are important in debugging services
	// and their specific problematic containers.
	//Meta *aws.Meta `json:"-"`

	// SimpleLog controls the log output. If true the output is simplified to help
	// reduce log volume during development or dev testing. If the bool flag is
	// false (default) the log output will be the standardized bnet splunk format
	SimpleLog bool `json:"-"`
}

// Close the friends service
func (f *Friends) Close() {
	close(f.manager.done)
	close(f.done)
}

// EnableDaemon enables background running routine
const EnableDaemon = true

// init will setup and configure providers, register gob enc/dec on interface
// types and any other required steps on init
func (f *Friends) init() error {
	for name, kind := range map[string]interface{}{} {
		gob.RegisterName(name, kind)
	}
	for _, pc := range f.config.Provider {
		np, err := provider.Open(pc)
		if err != nil {
			return err
		}
		if err := f.manager.addProvider(np); err != nil {
			return err
		}
	}
	if EnableDaemon {
		go f.manager.run()
	}
	return nil
}

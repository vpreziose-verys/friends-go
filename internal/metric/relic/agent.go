package relic

import (
	"log"

	newrelic "github.com/newrelic/go-agent"
)

var instance *Agent

// Open wraps Agent struct around the newrelic application
func Open(c Config) (*Agent, error) {
	instance = &Agent{Config: c}
	if !c.Enabled {
		return instance, nil // do not throw error on disabled relic
	}
	if c.Name == "" || c.Key == "" {
		instance.Config.Enabled = false
		return instance, ErrBadCreds
	}
	app, err := newrelic.NewApplication(newrelic.NewConfig(c.Name, c.Key))
	if err != nil {
		instance.Config.Enabled = false
	} else {
		instance.app = app
	}
	log.Println("relic: ready")
	return instance, err
}

// Agent is a wrapper for the newrelic agent/application
type Agent struct {
	Config

	app newrelic.Application
}

// Config used to control relic Agent
type Config struct {
	Name    string `json:"addr"`
	Key     string `json:"key"`
	Enabled bool   `json:"enabled"`
}

// Close handles Agent teardown and is deferred after main init
func (a *Agent) Close() {
	log.Println("relic: closing")
}

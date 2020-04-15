package client

import (
	"context"
	"encoding/json"
	"net/http"
)

// New creates a new bnet client used typically by providers
func New(addr, key, env string, c Cache) *Client {
	if c == nil {
		c = NoCache{}
	}
	return &Client{
		Addr:  addr,
		Cache: c,
		key:   key,
	}
}

// Client brokers http requests
type Client struct {
	Name  string
	Addr  string
	Cache Cache

	key string
}

// Do sends the http request
func (s *Client) Do(r *http.Request, data interface{}) error {
	// defer newrelic.StartSegment(tx, fmt.Sprintf("provider_%s_client_%s", s.Name, "request")).End()
	return s.fetch(r.Context(), r, &data)
}

func (s *Client) fetch(ctx context.Context, cr *http.Request, data interface{}) error {
	r, err := http.DefaultClient.Do(cr)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	rep := Reply{}
	rep.Platform.Message = data
	return json.NewDecoder(r.Body).Decode(&rep)
}

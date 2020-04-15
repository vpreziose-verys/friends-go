package provider

import "time"

const (
	// DefaultContentType used as default response content type
	DefaultContentType = "application/json"

	// DefaultTimeout used as default timeout on provider request
	DefaultTimeout = time.Duration(time.Second * 3)
)

// Config is an embedded struct in provider config structs to share generalized
// fields or functionality. Add common/generic fields here in future.
type Config struct {

	// ID is the providers system id (SystemID: from python->go conversion)
	ID string `json:"id"`

	// Name is the providers system name (SystemName: from python->go conversion)
	Name string `json:"name"`

	// Env used to point to the correct bnet environment (local, dev, test, ...)
	Env string `json:"env"`

	// Addr is the base url of service
	Addr string `json:"addr"`

	// Enabled controls provider state
	Enabled bool `json:"enabled"`

	// EnableLog turns on the provider logger
	EnableLog bool `json:"enable_logging"`

	// Timeout duration for provider requests (http_request_timeout_sec)
	Timeout time.Duration `json:"timeout"`

	// key holds onto the providers required key
	Key string `json:"-"`
}

// SetKey sets the key used by the provider
func (c *Config) SetKey(key string) {
	c.Key = key
}

package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/BethesdaNet/friends-go/internal/platform"
)

// Key controls server key validation
type Presence struct {
	PresenceConfig
	provider
}

// KeyConfig for key provider
type PresenceConfig struct {
	Config
	PresenceURL        string `json:"presence_url"`
	PresencePrivateURL string `json:"presence_private_url"`
}

// Close method will be called during teardown
func (p *Presence) Close() {}

// Check validates the key and returns the buid of the caller
func (p *Presence) Check(key, product string, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()
	return p.client.Do(p.check(ctx, key, product), out)
}

func (p *Presence) check(ctx context.Context, key, product string) *http.Request {
	const bodyfmt = `{"system_name":%q,"system_key":%q}`
	r, _ := http.NewRequest(http.MethodPost, p.Addr+p.PresenceURL, strings.NewReader(fmt.Sprintf(bodyfmt, key, product)))
	r.Header.Set(platform.HeaderKeyServer, p.Key)
	r.Header.Set("Content-Type", DefaultContentType)
	return r.WithContext(ctx)
}

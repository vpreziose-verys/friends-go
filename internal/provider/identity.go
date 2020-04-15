package provider

import (
	"context"
	"net/http"
	"strings"

	"github.com/BethesdaNet/friends-go/internal/platform"
	"github.com/BethesdaNet/friends-go/internal/provider/identity"
)

// Account shorthand
type Account = identity.Account

// Identity ...
type Identity struct {
	IdentityConfig
	provider
}

// IdentityConfig ...
type IdentityConfig struct {
	Config
	LookupURL string `json:"lookup_url"`
}

// Close method for cleaning tearing down the identity provider.
func (p *Identity) Close() {}

// GetAccount retrieves accounts by buid
func (p *Identity) GetAccount(id string, in *Account) error {
	data := []*Account{}
	if err := p.getAccounts([]string{id}, data); err != nil {
		return err
	}
	return nil
}

// GetAccounts retrieves accounts by buid array
func (p *Identity) GetAccounts(id []string, data []*Account) error {
	if err := p.getAccounts(id, data); err != nil {
		return err
	}
	return nil
}

func (p *Identity) getAccounts(id []string, data []*Account) error {
	r, _ := http.NewRequest(http.MethodGet, p.Addr+p.LookupURL+strings.Join(id, ","), nil)
	r.Header.Set(platform.HeaderKey, p.Key)
	r.Header.Set("Content-Type", DefaultContentType)
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()
	return p.client.Do(r.WithContext(ctx), data)
}

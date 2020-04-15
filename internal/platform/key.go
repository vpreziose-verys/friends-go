package platform

import (
	"errors"
	"net/http"
)

// Key represents decoded BNET key data
type Key struct {
	KeyID     int     `json:"api_key_id"`
	Key       string  `json:"api_key"`
	Kind      string  `json:"api_key_type"`
	Secret    string  `json:"secret"`
	Principal string  `json:"principal"`
	Product   Product `json:"product"`
	Ctime     int     `json:"ctime"`
}

// Product is used in bnet key data type
type Product struct {
	ID        int         `json:"product_id"`
	Name      string      `json:"name"`
	Platform  string      `json:"platform"`
	Datascope string      `json:"data_scope"`
	Ctime     int         `json:"ctime"`
	Mtime     int         `json:"mtime"`
	FPTID     interface{} `json:"first_party_title_id"`
}

// ErrInvalidKey returned when invalid bnet key data
var ErrInvalidKey = errors.New("invalid key")

// ApplyTo func updates request headers and set product in api log
func (k *Key) ApplyTo(r *http.Request) error {
	if k == nil {
		return ErrInvalidKey
	}
	if k.Key == "" || k.Product.Datascope == "" {
		return ErrInvalidKey
	}
	r.Header.Set(HeaderScope, k.Product.Datascope)
	r.Header.Set(HeaderProduct, k.Product.Name)
	r.Header.Set(HeaderService, k.Kind)
	return nil
}

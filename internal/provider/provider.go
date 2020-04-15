package provider

import (
	"net/http"
)

// Open returns a service provider based on the provided config
func Open(i interface{}) (Service, error) {
	if i == nil {
		return nil, ErrConfigNil
	}
	switch c := i.(type) {
	case *IdentityConfig:
		return mkIdentity(c)
	//case *KeyConfig:
	//	return mkKey(c)
	//case *NoteConfig:
	//	return mkNote(c)
	//case *StorageConfig:
	//	return mkStorage(c)
	//case *LanguageConfig:
	//	return mkLanguage(c)
	default:
		return nil, ErrConfigInvalid
	}
}

// Service interface to provide control over provider struct
type Service interface {
	SetClient(Client) error
	Close()
}

// Client interface used to allow better testing purposes
type Client interface {
	Do(*http.Request, interface{}) error
}

// Provider used as base embedded struct to broker calls to bnet services
type provider struct {
	client Client
}

func (p *provider) SetClient(c Client) error {
	if c == nil {
		return ErrNilClient
	}
	p.client = c
	return nil
}

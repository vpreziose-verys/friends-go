package provider

import (
	"strings"
	"time"

	"github.com/BethesdaNet/friends-go/internal/client"
)

func mkIdentity(c *IdentityConfig) (Service, error) {
	c.Timeout = ckTimeout(c.Timeout)
	return &Identity{
		IdentityConfig: *c,
		provider: provider{
			client: client.New(ckScheme(c.Addr), c.Key, c.Env, nil),
		},
	}, nil
}

func mkPresence(c *PresenceConfig) (Service, error) {
	c.Timeout = ckTimeout(c.Timeout)
	return &Presence{
		PresenceConfig: *c,
		provider: provider{
			client: client.New(ckScheme(c.Addr), c.Key, c.Env, nil),
		},
	}, nil
}

func mkNote(c *NoteConfig) (Service, error) {
	c.Timeout = ckTimeout(c.Timeout)
	return &Note{
		NoteConfig: *c,
		provider: provider{
			client: client.New(ckScheme(c.Addr), c.Key, c.Env, nil),
		},
	}, nil
}

//func mkLanguage(c *LanguageConfig) (Service, error) {
//	c.Timeout = ckTimeout(c.Timeout)
//	return &Language{
//		LanguageConfig: *c,
//		provider:       provider{},
//	}, nil
//}

//func mkStorage(c *StorageConfig) (Service, error) {
//	c.Timeout = ckTimeout(c.Timeout)
//	return &Storage{
//		StorageConfig: *c,
//		provider:      provider{},
//	}, nil
//}

func ckScheme(addr string) string {
	if !strings.HasPrefix(addr, "https://") {
		return "https://" + addr
	}
	return addr
}

func ckTimeout(d time.Duration) time.Duration {
	if d <= 0 {
		return DefaultTimeout
	}
	return d
}

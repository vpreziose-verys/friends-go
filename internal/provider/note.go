package provider

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/BethesdaNet/friends-go/internal/platform"
)

// Note provider brokers notification service transactions and caching
type Note struct {
	NoteConfig
	provider
}

// NoteConfig holds all required information for the notification provider
type NoteConfig struct {
	Config
	AnnounceURL string `json:"announcement_url"`
	NoteURL     string `json:"notification_url"`
}

// Notification struct used to send outbound annoucements or notifications
type Notification struct {
	Title string
	BUID  string
	Data  interface{}

	Announce bool
}

// Close method will be called during teardown
func (p *Note) Close() {}

// Announcement method handles sending out announcements to notification service
func (p *Note) Announcement(mt, msg, buid string, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()
	return p.client.Do(p.announce(ctx, mt, msg, buid, 1), out)
}

// Notification method handles sending out messages to notification service
func (p *Note) Notification(mt, msg, buid string, out interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), p.Timeout)
	defer cancel()
	return p.client.Do(p.note(ctx, mt, msg, buid), out)
}

func (p *Note) announce(ctx context.Context, mt, msg, buid string, version int) *http.Request {
	const bodyfmt = `{"message_type":%q,"payload":%q,"payload_version":%d}`
	r, _ := http.NewRequest(http.MethodPost, p.Addr+p.AnnounceURL, strings.NewReader(fmt.Sprintf(bodyfmt, mt, msg, version)))
	r.Header.Set(platform.HeaderKeyServer, p.Key)
	r.Header.Set("Content-Type", DefaultContentType)
	return r.WithContext(ctx)
}

func (p *Note) note(ctx context.Context, mt, msg, buid string) *http.Request {
	const bodyfmt = `{"message_type":%q,"payload":%s}`
	r, _ := http.NewRequest(http.MethodPost, p.Addr+p.NoteURL, strings.NewReader(fmt.Sprintf(bodyfmt, mt, msg)))
	r.Header.Set(platform.HeaderKeyServer, p.Key)
	r.Header.Set("Content-Type", DefaultContentType)
	return r.WithContext(ctx)
}

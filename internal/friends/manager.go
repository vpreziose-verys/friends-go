package friends

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/BethesdaNet/friends-go/internal/db/redis"
	"github.com/BethesdaNet/friends-go/internal/friends/status"
	"github.com/BethesdaNet/friends-go/internal/provider"
	"github.com/BethesdaNet/friends-go/internal/provider/identity"
)

const (
	// EnableSetOfflineOnDecodeErr flag if true will save users presence status as
	// offline (default) on decoder error to db preventing returning eronous errors
	EnableSetOfflineOnDecodeErr = false

	// EnableLegacyIdleTimestampFlow when enabled will query redis per buid for the
	// correlated idle timestamp (global or product/platform) and after parsing set
	// the value (if valid) onto the Status.Idle attr and not update Status.Time
	EnableLegacyIdleTimestampFlow = true

	// MaxScanCount limits how many keys can be requested from redis
	MaxScanCount = 50
)

type (
	// Manager brokers presence service and providers allowing for additional feats
	// to be implemented (rate limiting, retry logic, authentication, etc) and most
	// important a synchronization point where cache can be controlled correctly in
	// one spot rather than several components all handling data caching seperately
	Manager struct {

		// redis agent for db operations / handling
		dba *redis.Agent

		// platform subsystems providing functionality for each service allowing each
		// service provider its own configuration and http client instance
		identity *provider.Identity
		//key      *provider.Key
		//note     *provider.Note
		//storage  *provider.Storage
		//lang     *provider.Language

		// account map contains accounts retrieved from identity service
		account sync.Map

		// notechan specific for provider.Note outbound msgs
		//notes chan Notification

		// done chan closes the managers daemon which controls channel operations
		done chan struct{}
	}

	// Group used by manager to query a buid for all possible statuses using scan
	// and fill to query/populate buid keys based on wildcard searches
	Group struct {
		// Key array contains returned scanned keys used during loops as maps are not
		// in order to returned redis keys (especially if we need to create a default)
		Key []string

		// Data map contains scanned keys and data set by the fill func after scan
		// using map/interface for future iterations/structs representing statuses
		Data map[string]interface{}

		// Query is the original search value of scan
		Query string
	}

	// Account shorthand for identity.Account
	Account = identity.Account

	// Notification shorthand for provider.Notification
	//Notification = provider.Notification
)

// DefaultDaemonInterval is the minimum rate at which the background daemon sweep
const DefaultDaemonInterval = time.Second * 10

// Run controls managers channel data while providing synchronization between all
// providers. Using one centralized routine allows for Rate limiting (bucket),
// msg queueing for unreachable subsystems, or the ability to block all channels
// if internal components for presence need to be reinitialized.
//func (m *Manager) run() {
//	hz := newTicker(DefaultDaemonInterval)
//	defer hz.Stop()
//	for {
//		select {
//		case <-m.done:
//			return
//		case n, ok := <-m.notes:
//			if ok {
//				m.deliver(n)
//			}
//		case <-hz.C:
//			/* no-op */
//		}
//	}
//}

// SendNotification sends note msg through manager notechan if enabled
//func (m *Manager) SendNotification(title, buid string, data interface{}, annouce bool) {
//	if !EnableDaemon {
//		return
//	}
//	m.notes <- Notification{
//		Title:   title,
//		BUID:    buid,
//		Data:    data,
//		Annouce: annouce,
//	}
//}

// deliver processes inbound notifications from notes channel
//func (m *Manager) deliver(n Notification) {
//	b, _ := json.Marshal(n.Data)
//	if n.Annouce {
//		m.note.Announcement(n.Title, string(b), n.BUID, nil)
//	} else {
//		m.note.Notification(n.Title, string(b), n.BUID, nil)
//	}
//}

/* -------------------------------------------------------------------------- */

// GetBUID retrieves status by key. If key not found set offline status
// in the DB and return to handler.
//func (m *Manager) GetBUID(buid, product, platform, language, country string) (*Group, error) {
//	group, err := m.scan(buid)
//	if err == nil {
//		if err = m.fill(group); err != nil {
//			return nil, err
//		}
//		return group, nil
//	}
//	switch err {
//	case ErrScanZero:
//		log.Printf("group: scan err: %s", err)
//	default:
//		return nil, err
//	}
//	return group, nil
//}

// ErrScanZero returned if no results returned on scan of key
var ErrScanZero = errors.New("no records found")

//func (m *Manager) scan(value string) (*Group, error) {
//	wild := fmt.Sprintf(status.PaWildBUID, value)
//	keys, err := m.dba.Scan(wild, MaxScanCount)
//	if err != nil {
//		return nil, err
//	}
//	if len(keys) == 0 {
//		return nil, ErrScanZero
//	}
//	group := &Group{
//		Key:   keys,
//		Query: wild,
//		Data:  make(map[string]interface{}, len(keys)),
//	}
//	for _, k := range keys {
//		group.Data[k] = &Status{}
//	}
//	return group, nil
//}

var (
	// ErrGroupNil returned if provided group ptr is nil
	ErrGroupNil = errors.New("nil group")

	// ErrBadEncoding returned if decoding failed
	ErrBadEncoding = errors.New("unsupported data encoding retrieved")
)

func (m *Manager) fill(group *Group) error {
	if group == nil {
		return ErrGroupNil
	}
	if len(group.Key) > 1 {
		rows, err := m.dba.MGet(group.Key...)
		if err != nil {
			return err
		}
		for i, row := range rows {
			switch v := row.(type) {
			case []byte:
				if err := gob.NewDecoder(bytes.NewBuffer(v)).Decode(group.Data[group.Key[i]]); err != nil {
					continue
				}
			default:
				return ErrBadEncoding
			}
		}
	} else if len(group.Data) == 1 {
		key := group.Key[0]
		raw, err := m.dba.GetBytes(key)
		if err != nil {
			switch err {
			case redis.ErrBadKey:
				/* no-op */
			default:
				return err
			}
		} else {
			if err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(group.Data[key]); err != nil {
				return err
			}
		}
	}
	return nil
}

/* -------------------------------------------------------------------------- */

// GetSingleStatus retrieves status by key. If key not found set offline status
// in the DB and return to handler.
//func (m *Manager) GetSingleStatus(buid, product, platform, language, country string, in *Status) error {
//	key := status.Key(buid, product, platform, false)
//	raw, err := m.dba.GetBytes(key)
//	if err != nil {
//		switch err {
//		case redis.ErrBadKey:
//			return m.SetStatus(buid, product, platform, language, country, in.Set(status.Offline))
//		default:
//			return err
//		}
//	}
//	if err := gob.NewDecoder(bytes.NewBuffer(raw)).Decode(&in); err != nil {
//		// set offline status to prevent an invalid response to the client.
//		in.Set(status.Offline)
//
//		// Enable flag to save the offline status to redis. Decoding errors could be
//		// a result of bad data or other unforseen issues.
//		if EnableSetOfflineOnDecodeErr {
//			return m.SetStatus(buid, product, platform, language, country, in)
//		}
//	}
//
//	// if EnableLegacyIdleTimestampFlow is true query redis for idle timestamp for
//	// the requested key... in the future we should think of another way to set the
//	// idle timestamp instead of an additional key in redis
//	if EnableLegacyIdleTimestampFlow && in.Enum > status.AppearOffline {
//		if ts, ok := m.CheckIdle(buid); ok {
//			in.Idle = ts
//			return m.SetStatus(buid, product, platform, language, country, in.Set(status.Idle))
//		}
//	}
//	return nil
//}

/* -------------------------------------------------------------------------- */

// GetMultiStatus retrieves presence statuses from db
//func (m *Manager) GetMultiStatus(buid, product, platform, language, country string) (map[string]*Status, error) {
//	buids := strings.Split(buid, ",")
//	keys, delta := []string{}, []string{}
//	statuses := map[string]*Status{}
//
//	for _, id := range buids {
//		keys = append(keys)
//		statuses[buid] = &Status{BUID: id}
//	}
//
//	rows, err := m.dba.MGet(keys...)
//	if err != nil {
//		return statuses, err
//	}
//
//	// loop over rows and parse each row
//	for i, row := range rows {
//		b, ok := row.([]byte)
//		if !ok {
//			delta = append(delta, buids[i])
//		}
//		if err := gob.NewDecoder(bytes.NewBuffer(b)).Decode(statuses[buid]); err != nil {
//			continue
//		}
//		if len(delta) > 0 {
//			for _, id := range delta {
//				m.SetStatus(id, product, platform, language, country, statuses[id].Set(status.Offline))
//			}
//		}
//	}
//	return statuses, nil
//}

/* -------------------------------------------------------------------------- */

// SetStatus stores presence status in redis db
//func (m *Manager) SetStatus(buid, product, platform, language, country string, in *Status) error {
//	if in.Time.IsZero() {
//		in.Time = time.Now()
//	}
//	buf := &bytes.Buffer{}
//	if err := gob.NewEncoder(buf).Encode(&in); err != nil {
//		return err
//	}
//
//	key := status.Key(buid, product, platform, false)
//	if err := m.dba.SetBytes(key, buf.Bytes()); err != nil {
//		return err
//	}
//
//	// send notification through notechan to be processed
//	m.SendNotification("Presence Status Update", buid, in, false)
//
//	return nil
//}

/* -------------------------------------------------------------------------- */

// DelStatus removes presence status from db
//func (m *Manager) DelStatus(buid, product, platform string) error {
//	key := status.Key(buid, product, platform, false)
//	if EnableLegacyIdleTimestampFlow {
//		return m.dba.Del([]string{key, key + status.SxLastActivity}...)
//	}
//
//	if err := m.dba.Del(key); err != nil {
//		return err
//	}
//
//	// send notification through notechan to be processed
//	m.SendNotification("Presence Status Update", buid, nil, false)
//
//	return nil
//}

/* -------------------------------------------------------------------------- */

// CheckIdle queries redis and checks for last activity timestamp by buid
func (m *Manager) CheckIdle(buid string) (time.Time, bool) {
	key := fmt.Sprintf(status.KyBUIDLastActivity, buid)
	if raw, err := m.dba.GetBytes(key); err == nil {
		t, _ := time.Parse(time.RFC3339, string(raw))
		return t, !t.IsZero()
	}
	return time.Time{}, false
}

/* -------------------------------------------------------------------------- */

// GetAccounts returns map of accounts from identity response
func (m *Manager) GetAccounts(ids ...string) (map[string]*Account, error) {
	data := make(map[string]*Account, len(ids))
	for _, id := range ids {
		a := &Account{ID: id}
		if err := m.GetAccount(id, a); err != nil {
			return data, err
		}
		data[id] = a
	}
	return data, nil
}

// ErrAccountNotFound err returned when account not found by identity service
var ErrAccountNotFound = errors.New("account not found")

// GetAccount func returns account by buid
func (m *Manager) GetAccount(id string, in *Account) error {
	if ok := m.load(in); ok {
		return nil
	}
	if err := m.identity.GetAccount(id, in); err != nil {
		return ErrAccountNotFound
	}
	return m.store(in)
}

func (m *Manager) load(in interface{}) bool {
	switch v := in.(type) {
	case *Account:
		a, ok := m.account.Load(v.ID)
		if !ok {
			return false
		}
		in = *(a.(*Account))
		return true
	default:
		return false
	}
}

/* -------------------------------------------------------------------------- */

// ErrBadDataType returned when storing unsupported data type
var ErrBadDataType = errors.New("bad data type")

func (m *Manager) store(in interface{}) error {
	switch v := in.(type) {
	case *Account:
		m.account.Store(v.ID, v)
		return nil
	default:
		return ErrBadDataType
	}
}

func (m *Manager) delete(in interface{}) error {
	switch v := in.(type) {
	case *Account:
		m.account.Delete(v.ID)
		return nil
	default:
		return ErrBadDataType
	}
}

/* -------------------------------------------------------------------------- */

// ErrProviderNil returned when provider service interface is nil
var ErrProviderNil = errors.New("error adding provider; nil interface")

func (m *Manager) addProvider(in provider.Service) error {
	if in == nil {
		return ErrProviderNil
	}
	switch pt := in.(type) {
	case *provider.Identity:
		m.identity = pt
		//case *provider.Key:
		//	m.key = pt
		//case *provider.Note:
		//	m.note = pt
		//case *provider.Storage:
		//	m.storage = pt
		//case *provider.Language:
		//	m.lang = pt
	}
	return nil
}

/* -------------------------------------------------------------------------- */

// removeKey returns updated slice after removing provided array index. This also
// maintains slice order
func removeKey(arr []string, i int) []string {
	return append(arr[:i], arr[i+1:]...)
}

/* -------------------------------------------------------------------------- */

// ticker is a time.Ticker that knows the time.Duration it ticks at.
type ticker struct {
	time.Duration
	*time.Ticker
}

func newTicker(dur time.Duration) *ticker {
	return &ticker{
		Duration: dur,
		Ticker:   time.NewTicker(dur),
	}
}

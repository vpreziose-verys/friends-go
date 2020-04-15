package status

import (
	"encoding/json"
	"fmt"
	"math"
	"time"
)

// Status struct is the primary dataset
type Status struct {
	BUID     string `json:"buid"`
	Enum     Kind   `json:"-"`
	Product  string `json:"-"`
	Platform string `json:"-"`

	// status specific data attributes
	Global     *string `json:"status,omitempty"`
	Game       *string `json:"main_game_status,omitempty"`
	Player     *string `json:"player_status,omitempty"`
	Extended   *Ext    `json:"extended_game_status,omitempty"`
	Custom     *string `json:"custom_data,omitempty"`
	Connection *string `json:"connection_data,omitempty"`
	Joinable   *bool   `json:"joinable,omitempty"`

	Time   time.Time     `json:"-"`
	Idle   time.Time     `json:"-"`
	Expire time.Duration `json:"-"`
}

// Set updates the status
func (s *Status) Set(to Kind) *Status {
	switch s.Enum = to; s.Enum {
	case Online:
		joinable := false
		if s.Connection != nil && len(*s.Connection) > 0 {
			joinable = true
		}
		if s.Joinable == nil && joinable {
			s.Joinable = &joinable
		}
		if !s.Idle.IsZero() {
			s.Idle = time.Time{}
		}
	case Idle:
		if s.Idle.IsZero() {
			s.Idle = time.Now()
		}
	case DND:
	case AppearOffline:
	case Offline:
		return &Status{
			BUID:     s.BUID,
			Product:  s.Product,
			Platform: s.Platform,
			Time:     time.Now(),
		}
	}
	if s.Enum != Idle {
		s.Time = time.Now()
	}
	s.Expire = -1
	return s
}

// Key returns the expected db key to either fetch or store
func (s *Status) Key(global, idle bool) string {
	if global || s.Product == "" && s.Platform == "" {
		return Key(s.BUID, "", "", idle)
	}
	return Key(s.BUID, s.Product, s.Platform, idle)
}

/* -------------------------------------------------------------------------- */

// Alias used when marshalling/unmarshalling to prevent recursion
type Alias Status

// UnmarshalJSON overrides default json marhsalling for better control on the output
// in one place rather than multiple structs representating the same dataset.
func (s *Status) UnmarshalJSON(data []byte) error {
	a := (Alias)(*s)
	if err := json.Unmarshal(data, &a); err != nil {
		return err
	}
	return Check(Status(a), s)
}

// MarshalJSON overrides default json marhsalling for better control on the output
// in one place rather than multiple structs representating the same dataset.
func (s *Status) MarshalJSON() ([]byte, error) {
	if s.Enum > DND {
		return nil, ErrBadEnum
	}
	s.Global = &statuses[s.Enum]
	switch s.Enum {
	case Online:
		return json.Marshal(&OnlineResponse{(*Alias)(s)})
	case Offline, AppearOffline:
		return json.Marshal(&OfflineResponse{
			BUID:   s.BUID,
			Global: statuses[Offline],
			Info:   NewInfo(Offline, s.Time),
		})
	case Idle:
		return json.Marshal(&IdleResponse{(*Alias)(s), NewInfo(Idle, s.Time)})
	case DND:
		return json.Marshal(&DNDResponse{(*Alias)(s), NewInfo(DND, s.Time)})
	}
	return nil, ErrEncodingJSON
}

/* -------------------------------------------------------------------------- */

// Ext type required for omitempty
type Ext []ExtData

// ExtData struct used for extended game statuses item
type ExtData struct {
	ID  string        `json:"id"`
	Arg []interface{} `json:"arguments"`
}

/* -------------------------------------------------------------------------- */

// Action struct used to update a status with more constraints
type Action struct {
	To     Kind  `json:"-"`
	Enable *bool `json:"enable,omitempty"`
}

/* -------------------------------------------------------------------------- */

// NewInfo returns human readable attributes based on status kind and timestamp
func NewInfo(kind Kind, timestamp time.Time) Info {
	var elapsed time.Duration
	if timestamp.IsZero() {
		timestamp = time.Now()
	} else {
		elapsed = time.Now().Sub(timestamp)
	}
	value, unit := formatUnit(elapsed.Seconds())
	return Info{
		Seconds:   elapsed.Seconds(),
		Timestamp: timestamp.Format(time.RFC3339),
		Text:      fmt.Sprintf("%s for %d %s(s)", statuses[kind], int64(value), unit),
	}
}

// Info struct used when converting timestamp to human readable output
type Info struct {
	Seconds   float64 `json:"seconds"`
	Timestamp string  `json:"timestamp"`
	Text      string  `json:"text"`
}

// Unit consts defined below for estimated seconds used to format the Info struct
// text output. NOTE(xc): instead of including 3rd-party library to handle these
// values refactor / update how each block of time gets calculated
const (
	UnitMinute = 60
	UnitHour   = UnitMinute * 60
	UnitDay    = UnitHour * 24
	UnitWeek   = UnitDay * 7
	UnitMonth  = UnitWeek * 4
	UnitYear   = UnitMonth * 12
)

func formatUnit(duration float64) (float64, string) {
	switch {
	case duration <= UnitMinute:
		return duration, "second"
	case duration >= UnitMinute && duration <= UnitHour:
		return math.Round(duration / UnitMinute), "minute"
	case duration >= UnitHour && duration <= UnitDay:
		return math.Round(duration / UnitHour), "hour"
	case duration >= UnitDay && duration <= UnitWeek:
		return math.Round(duration / UnitDay), "day"
	case duration >= UnitWeek && duration <= UnitMonth:
		return math.Round(duration / UnitWeek), "week"
	case duration >= UnitMonth && duration < UnitYear:
		return math.Round(duration / UnitMonth), "month"
	case duration >= UnitYear:
		return math.Round(duration / UnitYear), "year"
	default:
		return duration, "forever"
	}
}

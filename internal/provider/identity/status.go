package identity

// Status enum type
type Status int

const (
	AppearOffline Status = iota
	Offline
	Online
	Idle
	DND
)

// statuses array uses enum int value as array index and string value
var statuses = [...]string{
	AppearOffline: "appear-offline",
	Offline:       "offline",
	Online:        "online",
	Idle:          "idle",
	DND:           "dnd",
}

// Key func satisfies the db and store to return the id
func (s Status) String() string {
	return statuses[s]
}

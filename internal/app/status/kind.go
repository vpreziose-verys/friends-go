package status

// Status Enums defaulting with Offline as zero
const (
	Offline       Kind = iota // 0
	AppearOffline             // 1
	Online                    // 2
	Idle                      // 3
	DND                       // 4
)

// Kind type modifies status behavior defined with the consts below
type Kind int

// String func returns status enum in human readable format
func (k Kind) String() string {
	return statuses[k]
}

var statuses = [...]string{
	Offline:       "offline",
	AppearOffline: "appear-offline",
	Online:        "online",
	Idle:          "idle",
	DND:           "dnd",
}

var keys = map[string]Kind{
	"offline":        Offline,
	"appear-offline": AppearOffline,
	"online":         Online,
	"idle":           Idle,
	"dnd":            DND,
}

// OfflineResponse returned if status.Enum == Offline/AppearOffline
type OfflineResponse struct {
	Global string `json:"status"`
	BUID   string `json:"buid"`
	Info   Info   `json:"offline_status"`
}

// OnlineResponse returned if status.Enum == Online
type OnlineResponse struct {
	*Alias
}

// IdleResponse returned if status.Enum == Idle
type IdleResponse struct {
	*Alias
	Info Info `json:"idle_status"`
}

// DNDResponse returned if status.Enum == DND
type DNDResponse struct {
	*Alias
	Info Info `json:"dnd_status"`
}

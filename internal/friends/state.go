package friends

// State struct contains attributes used during client requests. Middleware will
// add platform specific data populated from the sidecar
type State struct {

	// ID used for requests specifying a specific element by identifier
	ID string `json:"id"`

	// BUID used for bnet master account id
	BUID string `json:"buid"`

	// Key used for bnet key
	Key string `json:"key"`

	// Session used for bnet session data
	Session string `json:"session"`

	// Scope used for bnet session scope
	Scope string `json:"scope"`

	// Role used for bnet key service type
	Role string `json:"role"`

	// Platform used for bnet client platform
	Platform string `json:"platform"`

	// Product used for header product id
	Product string `json:"product"`

	// Language used for translation purposes
	Language string `json:"language"`

	// Country used for translation purposes
	Country string `json:"country"`

	// Finger used for client session fingerprint
	Finger string `json:"fp"`
}

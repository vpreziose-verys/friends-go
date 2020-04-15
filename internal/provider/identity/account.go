package identity

// Account represents a user account
type Account struct {
	ID   string `json:"account_id"`
	Name string `json:"username"`

	Type   string `json:"account_type"`
	State  State  `json:"account_state"`
	Status Status `json:"account_status"`

	Platform    string `json:"platform"`
	Environment string `json:"platform_environment"`
	QueryBUID   string `json:"query_buid"`
	Country     string `json:"country"`
	Language    string `json:"language"`
}

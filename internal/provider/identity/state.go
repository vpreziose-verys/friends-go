package identity

// State of account
type State string

const (
	StateMerged          State = "MERGED"
	StateAnonymous       State = "ANONYMOUS"
	StateIdentified      State = "IDENTIFIED"
	StateLinkedAnonymous State = "LINKED_ANONYMOUS"
)

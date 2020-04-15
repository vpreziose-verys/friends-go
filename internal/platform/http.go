package platform

const (
	// HeaderAgent header is for bnet agent type
	HeaderAgent = "X-BNET-Agent"

	// HeaderKey header is the client key (encoded jwk)
	HeaderKey = "X-BNET-Key"

	// HeaderKeyServer header for server key
	HeaderKeyServer = "X-Server-Key"

	// HeaderKeyMaster header is for master service key
	HeaderKeyMaster = "X-Master-Key"

	// HeaderBUID header is for specific buid
	HeaderBUID = "X-BNET-BUID"

	// HeaderAccount header is for master account id
	HeaderAccount = "X-BNET-Master-Account-ID"

	// HeaderSessionScope header for session scope
	HeaderSessionScope = "X-BNET-Session-Scope"

	// HeaderService header is the decoded key: type (client, admin, server)
	HeaderService = "X-BNET-Service-Type"

	// HeaderScope header is the decoded key: scope (foo)
	HeaderScope = "X-BNET-Data-Scope"

	// HeaderSession header for session token
	HeaderSession = "X-Session-Token"

	// HeaderPlatform header for client platform
	HeaderPlatform = "X-Platform"

	// HeaderPlatformBNET header is for bnet client platform
	HeaderPlatformBNET = "X-BNET-Platform"

	// HeaderProduct header is the decoded key: type (fallout, doom, etc)
	HeaderProduct = "X-BNET-Product"

	// HeaderFinger header for client fingerprint
	HeaderFinger = "X-Src-fp"

	// HeaderContent header is for the content type
	HeaderContent = "content-type"

	// ScopeBasic basic scope if no bnet-key (server, admin) not provided
	ScopeBasic = "basic"
)

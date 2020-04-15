package status

import "fmt"

// Key creates string key for get/set statuses or modifiers from redis
func Key(buid, product, platform string, active bool) string {
	var out string
	switch {
	case product != "" && platform != "":
		out = fmt.Sprintf(PaProductPlatformBUID, buid, product, platform)
	default:
		out = fmt.Sprintf(PaGlobal, buid)
	}
	if active {
		out += SxLastActivity
	}
	return out
}

// Ha=Hashes, Ky=Key, Pa=Pattern, Sx=Suffix, De=Delimeter
const (

	// PaWildBUID used for wildcard searches by buid
	PaWildBUID = "{%s}*"

	HaLocal                  = "default.local"
	HaLanguages              = "supported.languages"
	HaGameMain               = "game.statuses.hashes.main"
	HaGameExt                = "game.statuses.hashes.extended"
	KyIdle                   = "{%s}" + SxIdle
	KyOffline                = "{%s}" + SxOffline
	KyBUIDLastActivity       = "{%s}" + SxLastActivity
	KyProductPlatformIdle    = PaProductPlatform + SxIdle
	KyProductPlatformOffline = PaProductPlatform + SxOffline
	PaIdentity               = "{%s}.identity"
	PaKeys                   = "{%s}.keys"
	PaHashes                 = "{%s}.hashes"
	PaGameMain               = "{%s}:{%s}:m:{%s}"
	PaGameExt                = "{%s}:{%s}:e:{%s}"
	PaGlobal                 = "{%s}" + SxGlobal
	PaProductPlatform        = "%s.%s"
	PaProductPlatformBUID    = "{%s}." + PaProductPlatform
	PaProductPlatformBad     = "unavailable:{%s}:{%s}"
	SxGlobal                 = ".global_status"
	SxPlayer                 = ".player_status"
	SxGameMain               = ".main_game_status"
	SxGameExt                = ".extended_game_status"
	SxCustom                 = ".custom_data"
	SxJoinable               = ".joinable"
	SxConnection             = ".connection_data"
	SxDoNotDisturb           = ".dnd"
	SxLastActivity           = ".last_activity_timestamp"
	SxOffline                = ".offline"
	SxIdle                   = ".idle"
	SxKeys                   = ".keys"
	SxProcesses              = ".process_count"
	PxKeySpace               = "__keyspace@0__:"
	PsPublish                = "pmessage"
	PsExpired                = "expired"
	DeKey                    = "."
	DeValue                  = ","
	DeHash                   = ":"
	DeWild                   = "*"
)

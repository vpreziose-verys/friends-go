package status

const (
	// MaxCustomDataBytes for limiting custom data blob
	MaxCustomDataBytes = 2048

	// MaxConnectionDataSize limits the connection data size
	MaxConnectionDataSize = 2048

	// MaxGameStatusIDLength limits the game status id length
	MaxGameStatusIDLength = 50

	// MaxGameStatusDataLength limits the game status data length
	MaxGameStatusDataLength = 128

	// MaxExtendedGameStatuses limits how many extended game statuses can be set
	MaxExtendedGameStatuses = 2
)

// Check validates the status fields and data based on set rules from python
func Check(in Status, out *Status) error {

	// ensure either inbound status (from request) or out status that should contain
	// data attributes set by the manager during get/set requests.
	if in.BUID == "" && out.BUID == "" {
		return ErrBUID
	}

	/* ------------------------------------------------------------------------ */

	// Check Global (Enum.String()) Status and check against the valid key map
	if in.Global != nil && *in.Global != "" {
		if e, ok := keys[*in.Global]; ok {
			in.Enum = e
		} else {
			return ErrBadJSONStatus
		}
	}

	/* ------------------------------------------------------------------------ */

	if from, to := out.Enum, in.Enum; to != from {
		out.Enum = to
	}

	/* ------------------------------------------------------------------------ */

	if in.Player != nil && *in.Player != "" {
		out.Player = in.Player
	}

	/* ------------------------------------------------------------------------ */

	if in.Game != nil {
		if len(*in.Game) > MaxGameStatusDataLength {
			return ErrMaxGameStatusDataLength
		}
		out.Game = in.Game
	}

	/* ------------------------------------------------------------------------ */

	if in.Extended != nil {
		if len(*in.Extended) > MaxExtendedGameStatuses {
			return ErrMaxExtended
		}
		out.Extended = in.Extended
	}

	/* ------------------------------------------------------------------------ */

	if in.Custom != nil {
		if len(*in.Custom) > MaxCustomDataBytes {
			return ErrMaxCustomBytes
		}
		out.Custom = in.Custom
	}

	/* ------------------------------------------------------------------------ */

	if in.Connection != nil && len(*in.Connection) > 0 {
		out.Connection = in.Connection
		if in.Joinable != nil {
			if in.Enum > AppearOffline && *in.Joinable {
				out.Joinable = in.Joinable
			}
		}
	}

	/* ------------------------------------------------------------------------ */

	// finally set the global out status since we passed all checks
	out.Global = &statuses[in.Enum]

	return nil
}

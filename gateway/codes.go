package gateway

type DiscordCloseCode int64

const (
	CloseUnknownError DiscordCloseCode = iota + 4000
	CloseUnknownOpcode
	CloseDecodeError
	CloseNotAuthenticated
	CloseAuthenticationFailed
	CloseAlreadyAuthenticated
	_
	CloseInvalidSeq
	CloseRateLimited
	CloseSessionTimeout
	CloseInvalidShard
	CloseShardingRequired
	CloseInvalidAPIVersion
	CloseInvalidIntents
	CloseDisallowedIntents
)

func (d DiscordCloseCode) String() string {
	switch d {
	case CloseUnknownError:
		return "Unknown Error"
	case CloseUnknownOpcode:
		return "Unknown Opcode"
	case CloseDecodeError:
		return "Decode Error"
	case CloseNotAuthenticated:
		return "Not Authenticated"
	case CloseAuthenticationFailed:
		return "Authentication Failed"
	case CloseAlreadyAuthenticated:
		return "Already Authenticated"
	case CloseInvalidSeq:
		return "Invalid Sequence"
	case CloseRateLimited:
		return "Rate Limited"
	case CloseSessionTimeout:
		return "Session Timeout"
	case CloseInvalidShard:
		return "Invalid Shard"
	case CloseShardingRequired:
		return "Sharding Required"
	case CloseInvalidAPIVersion:
		return "Invalid API Version"
	case CloseInvalidIntents:
		return "Invalid Intents"
	case CloseDisallowedIntents:
		return "Disallowed Intents"
	default:
		return "Unknown"
	}
}

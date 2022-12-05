package evmfuncs

type Effect uint64

const (
	// EffectUnknown is the effect of a function call that does not change the state
	EffectUnknown Effect = 0

	// EffectRead is the effect of a function call that reads the state
	EffectRead Effect = 1 << iota

	// EffectWrite is the effect of a function call that changes any data on the blockchain
	EffectWrite

	// EffectTrigger is the effect of a function call that triggers an event
	EffectTrigger

	effectClassMax = 1 << iota >> 1
)

const (
	// EffectStateWrite is the effect of a function call that writes to a storage
	EffectStateWrite = EffectWrite | effectClassMax<<iota<<1

	// EffectAddressWrite is the effect of a function call that writes to an address
	EffectAddressWrite = EffectWrite | effectClassMax<<iota<<1

	// EffectTransfer is the effect of a function call that transfers funds
	EffectTransfer Effect = EffectStateWrite | EffectTrigger | effectClassMax<<iota<<1

	// EffectMint is the effect of a function call that mints tokens
	EffectMint = EffectTransfer | effectClassMax<<iota<<1

	// EffectBurn is the effect of a function call that burns tokens
	EffectBurn = EffectTransfer | effectClassMax<<iota<<1

	// EffectSelfdestruct is the effect of a function call that selfdestructs the contract
	EffectSelfdestruct = EffectAddressWrite | EffectStateWrite | EffectTrigger | effectClassMax<<iota<<1

	// EffectRBACUpdate is the effect of a function call that updates the RBAC
	EffectRBACUpdate = EffectStateWrite | effectClassMax<<iota<<1
)

func (e Effect) String() string {
	switch e {
	case EffectUnknown:
		return "unknown"
	case EffectRead:
		return "read"
	case EffectStateWrite:
		return "state write"
	case EffectAddressWrite:
		return "address write"
	case EffectTrigger:
		return "trigger"
	case EffectTransfer:
		return "transfer"
	case EffectMint:
		return "mint"
	case EffectBurn:
		return "burn"
	case EffectSelfdestruct:
		return "selfdestruct"
	case EffectRBACUpdate:
		return "rbac update"
	default:
		return "unknown"
	}
}

func (e Effect) IsUnknown() bool {
	return e == EffectUnknown
}

func (e Effect) Is(effect Effect) bool {
	return e&effect != 0
}

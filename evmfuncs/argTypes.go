package evmfuncs

type ArgKind uint8

// TODO: support fixed and ufixed

const (
	ArgKindUnknown ArgKind = iota
	ArgKindAddress
	ArgKindBool
	ArgKindBytes
	ArgKindInt
	ArgKindString
	ArgKindUint
)

func (k ArgKind) String() string {
	switch k {
	case ArgKindUnknown:
		return "unknown"
	case ArgKindAddress:
		return "address"
	case ArgKindBool:
		return "bool"
	case ArgKindBytes:
		return "bytes"
	case ArgKindInt:
		return "int"
	case ArgKindString:
		return "string"
	case ArgKindUint:
		return "uint"
	default:
		return "invalid"
	}
}

type ArgType struct {
	Kind      ArgKind
	OffsetLen int
}

func (t ArgType) String() string {
	return t.Kind.String()
}

var (
	ArgTypeUnknown = ArgType{Kind: ArgKindUnknown}
	ArgTypeAddress = ArgType{Kind: ArgKindAddress, OffsetLen: 20}
	ArgTypeBool    = ArgType{Kind: ArgKindBool, OffsetLen: 32}
	ArgTypeString  = ArgType{Kind: ArgKindString}

	ArgTypeUint = ArgType{Kind: ArgKindUint, OffsetLen: 32}
)

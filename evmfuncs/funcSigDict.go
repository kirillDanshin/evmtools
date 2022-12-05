package evmfuncs

import (
	"encoding/hex"
	"strings"

	"github.com/kirillDanshin/evmtools"
)

func init() {
	for k, v := range wellKnownRBACFuncs() {
		if _, ok := funcDescriptions[k]; ok {
			panic("dup function description " + k)
		}
		funcDescriptions[k] = v
	}

	for k, v := range funcDescriptions {
		v.knownMethodKey = k

		sig := v.name + "("
		for i, p := range v.inputs {
			sig += p.Type
			if i < len(v.inputs)-1 {
				sig += ","
			}
		}
		sig += ")"
		if k == sig { // only cache generic signatures
			methodID := evmtools.MethodID(sig)
			methodIDHex := hex.EncodeToString(methodID)
			v.methodIDHex = methodIDHex
			funcDescriptionByMethodID[methodIDHex] = v
		} else {
			if k != "transfer(address to, uint256 tokenID)" && k != "transfer(address to, uint256 value)" && k != "transfer(address to, uint256 amount)" {
				panic("signature mismatch: " + k + " != " + sig)
			}
		}
	}
}

type WellKnownFuncDesc struct {
	// Name is the name of the function
	name string

	knownMethodKey string

	// description is the description of the function
	description string

	// Inputs is the list of input parameters
	inputs []FuncParam

	// Outputs is the list of output parameters
	outputs []FuncParam

	// Effects is the list of side effects of the function call
	effects Effect

	methodIDHex string
}

func (desc *WellKnownFuncDesc) Name() string {
	return desc.name
}

func (desc *WellKnownFuncDesc) Description() string {
	return desc.description
}

func (desc *WellKnownFuncDesc) Inputs() []FuncParam {
	return desc.inputs
}

func (desc *WellKnownFuncDesc) Outputs() []FuncParam {
	return desc.outputs
}

func (desc *WellKnownFuncDesc) Effects() Effect {
	return desc.effects
}

func (desc *WellKnownFuncDesc) MethodIDHex() string {
	return desc.methodIDHex
}

func (desc *WellKnownFuncDesc) MethodID() []byte {
	mID, err := hex.DecodeString(desc.methodIDHex)
	if err != nil {
		sel := desc.name + "("
		for i, p := range desc.inputs {
			sel += p.Type
			if i < len(desc.inputs)-1 {
				sel += ","
			}
		}

		sel += ")"

		return evmtools.MethodID(sel)
	}

	return mID
}

func GetWellKnownFuncByMethodID(methodID []byte) (*WellKnownFuncDesc, bool) {
	mIDHex := hex.EncodeToString(methodID)
	desc, ok := funcDescriptionByMethodID[mIDHex]
	return desc, ok
}

func GetWellKnownFuncBySig(sig string) (*WellKnownFuncDesc, bool) {
	desc, ok := funcDescriptionByMethodID[sig]
	return desc, ok
}

func GetWellKnownFuncsByName(name string) []*WellKnownFuncDesc {
	descriptions := []*WellKnownFuncDesc{}

	for _, desc := range funcDescriptions {
		if strings.EqualFold(desc.name, name) {
			descriptions = append(descriptions, desc)
		}
	}

	return descriptions
}

var funcDescriptionByMethodID = map[string]*WellKnownFuncDesc{}

var funcDescriptions = map[string]*WellKnownFuncDesc{
	"balance(address)": {
		name:        "balance",
		description: "get the balance of the given account",
		inputs: []FuncParam{
			{
				Name: "accountAddress",
				Type: "address",
			},
		},
		outputs: []FuncParam{
			{
				Name: "balance",
				Type: "uint256",
			},
		},
		effects: EffectRead,
	},
	"mint(address,uint256)": {
		name:        "mint",
		description: "mint new tokens",
		inputs: []FuncParam{
			{
				Name: "receiverAddress",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectMint,
	},

	"burn(uint256)": {
		name:        "burn",
		description: "burn tokens",
		inputs: []FuncParam{
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectBurn,
	},

	"burn(address,uint256)": {
		name:        "burn",
		description: "burn tokens",
		inputs: []FuncParam{
			{
				Name: "accountAddress",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
	},

	"transfer(address,uint256)": {
		name:        "transfer",
		description: "transfer erc20 tokens or a specific erc721 to the given address",
		inputs: []FuncParam{
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"transfer(address to, uint256 tokenID)": {
		name:        "transfer",
		description: "transfer a specific erc721 to the given address",
		inputs: []FuncParam{
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "tokenID",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"transfer(address to, uint256 value)": {
		name:        "transfer",
		description: "transfer erc20 tokens to the given address",
		inputs: []FuncParam{
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "value",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"transfer(address to, uint256 amount)": {
		name:        "transfer",
		description: "transfer erc20 tokens to the given address",
		inputs: []FuncParam{
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"transferFrom(address,address,uint256)": {
		name:        "transferFrom",
		description: "transfer tokens from the given address to the given address",
		inputs: []FuncParam{
			{
				Name: "from",
				Type: "address",
			},
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"approve(address,uint256)": {
		name:        "approve",
		description: "approve the given address to spend the specified number of tokens on behalf of the message sender",
		inputs: []FuncParam{
			{
				Name: "spender",
				Type: "address",
			},
			{
				Name: "amount",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectStateWrite,
	},
	"name()": {
		name:        "name",
		description: "get the name of the token",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "name",
				Type: "string",
			},
		},
		effects: EffectRead,
	},
	"owner()": {
		name:        "owner",
		description: "get the owner of the contract",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "owner",
				Type: "address",
			},
		},
		effects: EffectRead,
	},
	"symbol()": {
		name:        "symbol",
		description: "get the symbol of the token",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "symbol",
				Type: "string",
			},
		},
		effects: EffectRead,
	},
	"totalSupply()": {
		name:        "totalSupply",
		description: "get the total supply of the token",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "totalSupply",
				Type: "uint256",
			},
		},
		effects: EffectRead,
	},
	"decimals()": {
		name:        "decimals",
		description: "get the number of decimals of the token",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "decimals",
				Type: "uint8",
			},
		},
		effects: EffectRead,
	},
	"balanceOf(address)": {
		name:        "balanceOf",
		description: "get the balance of the given address",
		inputs: []FuncParam{
			{
				Name: "accountAddress",
				Type: "address",
			},
		},
		outputs: []FuncParam{
			{
				Name: "balance",
				Type: "uint256",
			},
		},
		effects: EffectRead,
	},
	"allowance(address,address)": {
		name:        "allowance",
		description: "get the number of tokens that the given address is allowed to spend on behalf of the given address",
		inputs: []FuncParam{
			{
				Name: "owner",
				Type: "address",
			},
			{
				Name: "spender",
				Type: "address",
			},
		},
		outputs: []FuncParam{
			{
				Name: "allowance",
				Type: "uint256",
			},
		},
		effects: EffectRead,
	},
	"renounceOwnership()": {
		name:        "renounceOwnership",
		description: "renounce ownership of the contract",
		inputs:      []FuncParam{},
		outputs:     []FuncParam{},
		effects:     EffectStateWrite,
	},
	"transferOwnership(address)": {
		name:        "transferOwnership",
		description: "transfer ownership of the contract to the given address",
		inputs: []FuncParam{
			{
				Name: "newOwner",
				Type: "address",
			},
		},
		outputs: []FuncParam{},
		effects: EffectStateWrite,
	},
	"pause()": {
		name:        "pause",
		description: "pause the contract",
		inputs:      []FuncParam{},
		outputs:     []FuncParam{},
		effects:     EffectStateWrite,
	},
	"unpause()": {
		name:        "unpause",
		description: "unpause the contract",
		inputs:      []FuncParam{},
		outputs:     []FuncParam{},
		effects:     EffectStateWrite,
	},
	"paused()": {
		name:        "paused",
		description: "check if the contract is paused",
		inputs:      []FuncParam{},
		outputs: []FuncParam{
			{
				Name: "paused",
				Type: "bool",
			},
		},
		effects: EffectRead,
	},
	"ownerOf(uint256)": {
		name:        "ownerOf",
		description: "get the owner of the given token",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{
			{
				Name: "owner",
				Type: "address",
			},
		},
		effects: EffectRead,
	},
	"getApproved(uint256)": {
		name:        "getApproved",
		description: "get the approved address for the given token",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{
			{
				Name: "operator",
				Type: "address",
			},
		},
		effects: EffectRead,
	},
	"setApprovalForAll(address,bool)": {
		name:        "setApprovalForAll",
		description: "set the approval status for the given operator",
		inputs: []FuncParam{
			{
				Name: "operator",
				Type: "address",
			},
			{
				Name: "approved",
				Type: "bool",
			},
		},
		outputs: []FuncParam{},
		effects: EffectStateWrite,
	},
	"isApprovedForAll(address,address)": {
		name:        "isApprovedForAll",
		description: "check if the given operator is approved for the given owner",
		inputs: []FuncParam{
			{
				Name: "owner",
				Type: "address",
			},
			{
				Name: "operator",
				Type: "address",
			},
		},
		outputs: []FuncParam{
			{
				Name: "approved",
				Type: "bool",
			},
		},
		effects: EffectRead,
	},
	"safeTransferFrom(address,address,uint256)": {
		name:        "safeTransferFrom",
		description: "transfer the given token from the given address to the given address",
		inputs: []FuncParam{
			{
				Name: "from",
				Type: "address",
			},
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},

	"safeTransferFrom(address,address,uint256,bytes)": {
		name:        "safeTransferFrom",
		description: "transfer the given token from the given address to the given address",
		inputs: []FuncParam{
			{
				Name: "from",
				Type: "address",
			},
			{
				Name: "to",
				Type: "address",
			},
			{
				Name: "tokenId",
				Type: "uint256",
			},
			{
				Name: "data",
				Type: "bytes",
			},
		},
		outputs: []FuncParam{},
		effects: EffectTransfer,
	},
	"tokenURI(uint256)": {
		name:        "tokenURI",
		description: "get the URI of the given token",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{
			{
				Name: "tokenURI",
				Type: "string",
			},
		},
		effects: EffectRead,
	},
	"supportsInterface(bytes4)": {
		name:        "supportsInterface",
		description: "check if the contract implements the given interface",
		inputs: []FuncParam{
			{
				Name: "interfaceId",
				Type: "bytes4",
			},
		},
		outputs: []FuncParam{
			{
				Name: "supported",
				Type: "bool",
			},
		},
		effects: EffectRead,
	},
	"onERC721Received(address,address,uint256,bytes)": {
		name:        "onERC721Received",
		description: "handle the receipt of an NFT",
		inputs: []FuncParam{
			{
				Name: "operator",
				Type: "address",
			},
			{
				Name: "from",
				Type: "address",
			},
			{
				Name: "tokenId",
				Type: "uint256",
			},
			{
				Name: "data",
				Type: "bytes",
			},
		},
	},

	"isOnLegalHold(uint256)": {
		name:        "isOnLegalHold",
		description: "check if the given token is on legal hold",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{
			{
				Name: "onHold",
				Type: "bool",
			},
		},
		effects: EffectRead,
	},
	"putOnLegalHold(uint256)": {
		name:        "putOnLegalHold",
		description: "put the given token on legal hold",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectStateWrite,
	},
	"releaseLegalHold(uint256)": {
		name:        "releaseLegalHold",
		description: "release the given token from legal hold",
		inputs: []FuncParam{
			{
				Name: "tokenId",
				Type: "uint256",
			},
		},
		outputs: []FuncParam{},
		effects: EffectStateWrite,
	},
}

func getRBACFuncsForRole(role string) map[string]*WellKnownFuncDesc {

	firstCap := strings.ToUpper(role[:1]) + role[1:]

	return map[string]*WellKnownFuncDesc{
		"add" + firstCap + "(address)": {
			name:        "add" + firstCap,
			description: "add the given address to the " + role + " role",
			inputs: []FuncParam{
				{
					Name: "account",
					Type: "address",
				},
			},
			outputs: []FuncParam{},
			effects: EffectRBACUpdate,
		},
		"remove" + firstCap + "(address)": {
			name:        "remove" + firstCap,
			description: "remove the given address from the " + role + " role",
			inputs: []FuncParam{
				{
					Name: "account",
					Type: "address",
				},
			},
			outputs: []FuncParam{},
			effects: EffectRBACUpdate,
		},
		"renounce" + firstCap + "()": {
			name:        "renounce" + firstCap,
			description: "remove the sender from the " + role + " role",
			inputs:      []FuncParam{},
			outputs:     []FuncParam{},
			effects:     EffectRBACUpdate,
		},
		"has" + firstCap + "Role(address)": {
			name:        "has" + firstCap + "Role",
			description: "check if the given address has the " + role + " role",
			inputs: []FuncParam{
				{
					Name: "account",
					Type: "address",
				},
			},
			outputs: []FuncParam{
				{
					Name: "hasRole",
					Type: "bool",
				},
			},
			effects: EffectRead,
		},
		"is" + firstCap + "()": {
			name:        "is" + firstCap,
			description: "check if the sender has the " + role + " role",
			inputs:      []FuncParam{},
			outputs: []FuncParam{
				{
					Name: "is" + firstCap,
					Type: "bool",
				},
			},
			effects: EffectRead,
		},
		"is" + firstCap + "(address)": {
			name:        "is" + firstCap,
			description: "check if the given address has the " + role + " role",
			inputs: []FuncParam{
				{
					Name: "account",
					Type: "address",
				},
			},
			outputs: []FuncParam{
				{
					Name: "is" + firstCap,
					Type: "bool",
				},
			},
			effects: EffectRead,
		},
	}
}

var wellKnownRBACFuncs = func() map[string]*WellKnownFuncDesc {
	wellKnownRoles := []string{
		"admin",
		"minter",
		"pauser",
		"burner",
		"signer",
		"whitelisted",
		"blacklisted",
		"owner",
		"operator",
		"verifier",
		"legal",
		"legalOperator",
		"legalHolder",
		"legalHoldOperator",
	}
	funcs := map[string]*WellKnownFuncDesc{}
	for _, role := range wellKnownRoles {
		for k, v := range getRBACFuncsForRole(role) {
			funcs[k] = v
		}
	}
	return funcs
}
